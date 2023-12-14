# go-configserver vnext

```
This repository is work in progress
```

Inspired by spring-cloud-config, go-configserver is aimed as cloud/kubernetes payloads and allows to centrally manage your application configuration using gitops.
Centrally manage your configuration within a git repositorz - the configuration will be automagically updated on your running pods.

go-configserver supports :

- mutliple git repositories for one central instance
- sensitive content encryption
- client based repository access
- and much more to come

## Philosophy

go-configserver aims at simplicity and security (will it tries to). Goal is to limit the configuration hassle and make use of conventions.

## Configuring

### General configuration

go-configserver is configured via a central `configserver.yml` file.
This file is located at startup either using the `-c` command line argument or is located in the path pointed by the `CONFIGSERVER_HOME`environment variable.
If nothing is provided, go-configserver will attempt to locate this file in the `/var/run/configserver` directory

```yaml
# Location in which encryption keys can be found
certsLocation:  /var/run/configserver/certs
server:
  # Network interface and port on which the server listens for incoming HTTP requests
  listenOn: ":8080"
git:
  # Path where git repositories configs are stored   
  repositoriesConfigurationLocation: ./samples/home/repositories
```

### Git repositories configuration

go-configserver supports multiple repositories - each repo needs to be configured in a separate yaml file :

```yaml
# Name of the repository - needs to be unique
name: fav
# Repository checkout URL
url: https://github.com/fredjeck/fav
# Repository checkout interval in seconds
refreshIntervalSeconds: 10
# Repository checkout location
checkoutLocation: ./samples/home/git/fav
# List of allowed client ids
clients:
  - myclientid
```

## Using go-configserver

### Registering a new client

```http request
POST http://localhost:8080/api/register HTTP/1.1
content-type: application/json

{
    "client_id": "sample_client"
}
```

Will generate a client secret for the provided client id

```json
{
  "client_id": "sample_client",
  "client_secret": "SECRET"
}
```

### Accessing Repository files

Before accessing files, the client needs to obtain a bearer token.
The targetted repositories names need to be specified using the scope parameter
Client authentication needs to be passed using the basic authentication header
```http request
POST http://localhost:8080/oauth2/authorize HTTP/1.1
content-type: application/x-www-form-urlencoded
Authorization: Basic B64(client_id:client_secret)

grant_type=client_credentials&scope=repo1 repo2 repo3
```

```json
{
  "access_token": "BEARER",
  "token_type": "bearer",
  "expires_in": 86400,
  "scope": ""
}
```

repository data can then be accessed via the `/git/REPO_NAME/PATH` for instance */git/go-configserver/README.md*

```http request
GET http://localhost:8080/git/go-configserver/README.md HTTP/1.1
Authorization: Bearer BEARER_TOKEN
```

## Encrypting sensitive content

Sensitive content like passwords can be encrypted using the encrypt (for a single value) or tokenize (for a whole pre-tokenizen file) endpoints

