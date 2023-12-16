FROM golang:alpine3.19 as build
# Step one build everything
WORKDIR /usr/src/configserver

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/src/configserver/configserver ./cmd

# Generate a different runtime image - smaller
FROM alpine:3.19 as runtime
COPY --from=build /usr/src/configserver/configserver /usr/bin

# Create a user and its group for configserver
ARG PUID=4200
ARG PGID=4200

RUN addgroup --gid ${PGID} configserver \
    && adduser --disabled-password --uid ${PUID} -G configserver  --gecos "" --home /configserver configserver

RUN mkdir /var/run/configserver \
    && chown configserver:configserver /var/run/configserver

# Create a default configuration file in case nothing is provided
RUN echo -e "certsLocation:  /var/run/configserver/certs \n\
server: \n\
  listenOn: ":8080" \n\
git: \n\
  repositoriesConfigurationLocation: /var/run/configserver/repositories" >> /var/run/configserver/configserver.yml


# Non root users are the best
USER configserver
WORKDIR /configserver

EXPOSE 8080

CMD ["configserver"]