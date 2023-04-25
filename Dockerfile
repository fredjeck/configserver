#LABEL authors="fred"

FROM golang:1.20-bullseye as gobuild
WORKDIR /build
COPY . .
# Installs Go dependencies
RUN go mod download
# Build configserver
RUN go build ./cmd/configserver

FROM golang:1.20-bullseye as runtime
WORKDIR /var/run/configuserver
COPY --from=gobuild /build/configserver .

# Tells Docker which network port your container listens on
# EXPOSE 8080
# Specifies the executable command that runs when the container starts
CMD ["/var/run/configserver"]