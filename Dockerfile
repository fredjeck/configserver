FROM golang:1.22.1-alpine3.19 as build
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

RUN mkdir /var/run/configserver && chown configserver:configserver /var/run/configserver

# Create a default configuration file in case nothing is provided
RUN echo -e "environment:\n\
  kind: production\n\
  home: /var/run/configserver\n\
\n\
server:\n\
  passPhrase: I am not secure please change me\n\
  listenOn: ":4200"\n\
  secretExpiryDays: 365\n\
  validateSecretLifeSpan: false\n\
\n\
repositories:\n\
  checkoutLocation: /tmp/configserver" >> /var/run/configserver/configserver.yml


# Non root users are the best
USER configserver
WORKDIR /configserver
ENV CONFIGSERVER_HOME="/var/run/configserver"

EXPOSE 4200

VOLUME ["/var/run/configserver"]

CMD ["configserver"]