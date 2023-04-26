#LABEL authors="fred"

FROM golang:1.20-bullseye as gobuild
WORKDIR /build
COPY . .
# Installs Go dependencies
RUN go mod download
# Build configserver
RUN go build ./cmd/configserver

FROM node:20 as nodebuild
WORKDIR /build
COPY ./ui .
RUN npm install
RUN npm run build

FROM golang:1.20-bullseye as runtime
WORKDIR /configserver
COPY --from=gobuild /build/configserver .
COPY --from=nodebuild /build/dist ./static
RUN mkdir /var/run/configserver/
RUN mkdir /var/run/configserver/repositories

ARG CONFIGSERVER_HOME=/configserver
ENV CONFIGSERVER_HOME=$CONFIGSERVER_HOME

ARG CONFIGSERVER_REPOSITORIES=/var/run/configserver/repositories
ENV CONFIGSERVER_REPOSITORIES=$CONFIGSERVER_REPOSITORIES

ARG CONFIGSERVER_CFG=/var/run/configserver
ENV CONFIGSERVER_CFG=$CONFIGSERVER_CFG

ARG PORT=8090
EXPOSE $PORT
CMD ["/configserver/configserver"]