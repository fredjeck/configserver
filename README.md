# configserver

```
This repository is work in progress

Current status : TOTAL MESS
```

Inspired by spring-cloud-config, configserver is aimed as cloud/kubernetes payloads and allows to centrally manage your application configuration using gitops.
Centrally manage your configuration within a git repositorz - the configuration will be automagically updated on your running pods.

configserver supports :

- mutliple git repositories for one central instance
- sensitive content encryption
- client based repository access
- and much more to come

## Philosophy

Configserver aims at simplicity and security (well it tries to). Goal is to limit the configuration hassle and make use of conventions.

## Configuring

### General configuration

Configserver is configured via a central `configserver.yml` file.
This file is located at startup either using the `-c` command line argument or is located by the path pointed by the `CONFIGSERVER_HOME`environment variable.
If nothing is provided, configserver will attempt to locate this file in the `/var/run/configserver` directory

```yaml
# Location in which encryption keys can be found
certsLocation:  /var/run/configserver/certs
server:
  # Network interface and port on which the server listens for incoming HTTP requests
  listenOn: ":8080"
  authorization:
    - basic
    - bearer
git:
  # Path where git repositories configs are stored   
  repositoriesConfigurationLocation: ./samples/home/repositories
```

### Git repositories configuration

configserver supports multiple repositories - each repo needs to be configured in a separate yaml file :

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

The `clients` array lists all the client ids allowed to access the repository, please see *registering a new client* below

## Using configserver

Configserver is mostly used via its API

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

### Auhtorization

#### Configuring Authorization

Configserver supports the following authorization schemes: Basic, Bearer or None.
Enabled authorization schemes are defined in the main configuration file.

#### Bearer Tokens

Before accessing files, the client needs to obtain a bearer token.
The targeted repository names need to be specified using the scope parameter.
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

repository data can then be accessed via the `/git/REPO_NAME/PATH` for instance */git/configserver/README.md*

```http request
GET http://localhost:8080/git/configserver/README.md HTTP/1.1
Authorization: Bearer BEARER_TOKEN
```

## Encrypting sensitive content

Sensitive content like passwords can be encrypted using the encrypt (for a single value) or tokenize (for a whole pre-tokenizen file) endpoints.
go-config server replaces sensitive values by encrypted tokens using the following formalism `{enc:ENCRYPTED_VALUE}`

You can encrypt a single value using the encrypt endpoint or using the tokenize command:
```http request
POST http://localhost:8080/api/encrypt HTTP/1.1
content-type: application/json

This text value will be encrypted
This one as well
```

Which will return

```json
{
  "token": "{enc:ZkcF7Xk+bnU6axHs/UdtmXKQxVS71+7a13ctfYrRhpbXeKW2ZnkzFujwzx4IJcAGppgdd9hybsrEXA8YUbB1+CqAFjcQj8Yfzi+HuxV1}"
}
```

or you could pre-tokenize your configuration file
```http request
POST http://localhost:8080/api/tokenize HTTP/1.1
content-type: text/plain

contentToTokenize:
 -p1: '{enc:value1}'
 -p2: '{enc:value2}'
```

Which would return the tokenized configuration ready to be copied and pasted
```yaml
contentToTokenize:
 -p1: '{enc:ZkcF7Xk+bnU6axHs/UdtmXKQxVS71+7a13ctfYrRhpbXeKW2ZnkzFujwzx4IJcAGppgdd9hybsrEXA8YUbB1+CqAFjcQj8Yfzi+HuxV1}'
 -p2: '{enc:ZkcF7Xk+bnU6axHs/UdtmXKQxVS71+7a13ctfYrRhpbXeKW2ZnkzFujwzx4IJcAGppgdd9hybsrEXA8YUbB1+CqAFjcQj8Yfzi+HuxV1}'
```

### Using the tokenize command

Like with the endpoint, prepare your files to be tokenized.
Then simply call the following command (be aware you will require the same private key used on your instance to be available locally)

```shell
configserver tokenize -k ./certs -f configuration.yaml -o tokenized.yaml
```

## Using the provided container image

```shell
# Build the image
$ podman build . -t configserver:latest
$ podmarun -it --rm -p 8080:8080 -v ./samples/home/certs/:/configserver/certs -v ./samples/home/repositories/:/configserver/repositories localhost/configserver:latest
```