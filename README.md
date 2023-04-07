# go-configserver
A Spring cloud-config inspired clone written in Go

`This is a work in progress project and is not ready for production usage`

## About

go-configserver is built for those scenarios where you have mutliple micro-services running within your namespace and you want to be able to version their configuration and push changes without redeploying your pods.
Simply store your configuration files in a git repository and configure go-configure to securely serve your configuration files to your microservices.

## Configuration

go-configserver will try to load its master configuration file from the path pointed by the `$CONFIGSERVER_HOME` environment variable, if this variable is not set or if it points to an unexisting directory, the configuration will be loaded from `/var/run/configserver`.
go-confiserver expects a file named configserver.yaml at this location.

```yaml
cacheStorageSeconds: 120
cacheEvictorIntervalSeconds: 10
listenOn: ":8090"
repositories:
  - name: configserver
    token: personnal_access_token
    url: https://github.com/fredjeck/go-configserver
    refreshinterval: 600
    clients:
      - bda6f1d0-ecbc-4da0-b513-8c5555b6c155
```
- **cacheStorageSeconds:** Duration a served file is maintained in memory once served
- **cacheEvictorIntervalSeconds:** Interval at which the cache evictor searches for outdated files
- **listenOn:** IP:PORT on which go-configserver listends for incoming connections
- **repositories/name:** The repository name is important for go-configsever, as it will be used for server repository contents. A repository named _myrepo_ will be a served from the _/git/myrepo_ url.
  go-configserver also strongly binds generated client secrets to a repository, therefore choose your names wisely
- **repositories/name:** List of Client IDs allowed to access the repository
- **repositories/token:** PAT used to checkout the repository, at the moment only PATs are supported

## Registering new clients

Adding clients to a repository is a two step process.

First you will need to generate a client secret for this client/repository combination:
```shell
configserver register -r configserver -i clientid
Repository: configserver
ClientId: clientid
ClientSecret: 1MQ4PT2L5bSgt/04vaP1bApZM/tMmuWcwQkAOpg88iTeYiaMMjXxQ1kDtXhKc+M0DiD+IVMOJQVvdBBCgV+2Kd8aBQz+tIY+GI4uGO0=
Please store the client secret carefully and do not forget to register the ClientID in the configserver.yaml file
```
Then you will need to declare the client id in the repository's client list

**It is important to note that a ClientSecret is bound to the repositories declared when generating the token and cannot be used to access another repository**

## Securing sensitive values
If your configuration files hold sensitive values (you should be using a kubernetes secret or a keyvault you know) you can encrypt those values in your files.
Those will be automatically replaced when served.

Encryption can be done from the command line or via API call :
```shell
configserver encrypt -v thisissometexttoencrypt
Encrypted token: {enc:qHlW+n0xsb54KnzcNDzWthkorBuWnNcHlNU9WBDvG+Siz2MCyyI9v/sdz7pOAdNSTu+W}
```

```shell
curl --location 'http://localhost:8090/api/encrypt' \
--header 'Content-Type: text/plain' \
--data 'Yeah'
```