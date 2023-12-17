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

RUN mkdir /configserver/repositories &&  chown configserver:configserver /configserver/repositories \
    && mkdir /configserver/certs &&  chown configserver:configserver /configserver/certs

# Create a default configuration file in case nothing is provided
RUN echo -e "certsLocation:  /configserver/certs \n\
server: \n\
  listenOn: ":8080" \n\
git: \n\
  repositoriesConfigurationLocation: /configserver/repositories" >> /configserver/configserver.yml


# Non root users are the best
USER configserver
WORKDIR /configserver
ENV CONFIGSERVER_HOME="/configserver"

EXPOSE 8080

CMD ["configserver"]

VOLUME ["/configserver/certs", "/configserver/repositories"]