# str-ap-internal 
[![Go](https://github.com/thierryturpin/str-ap-internal/actions/workflows/go.yaml/badge.svg)](https://github.com/thierryturpin/str-ap-internal/actions/workflows/go.yaml)
[![Pulumi-charts](https://github.com/thierryturpin/str-ap-infra-internal/actions/workflows/pulumi-charts.yml/badge.svg)](https://github.com/thierryturpin/str-ap-infra-internal/actions/workflows/pulumi-charts.yml)  

Private development repository used for https://github.com/SEMICeu/STR-AP

## Setup notes
To generate a local certificate use: https://github.com/FiloSottile/mkcert

For testing the endpoints
* provided notebook
* [collection.http](collection.http)

### Swagger
Swagger ui: https://str.local/swagger/index.html
Update of the swagger documentation, can be done by running the following command from the project root:
```
swag fmt  -g ./cmd/str/main.go
swag init -g ./cmd/str/main.go -o docs
```

### helm chart setup notes

```
cd chart/str

# Create a chart name sparse
helm create str

# Uninstall
helm uninstall str -n str

# Create or update
helm upgrade -i str . --namespace str --create-namespace --set image.password=ghp***

helm upgrade -i str . --namespace str --create-namespace \
--set image.password=ghp_*** \
--set tls.key=$(base64 -i ../../str.local-key.pem) \
--set tls.crt=$(base64 -i ../../str.local.pem) \
--set kafka.boostrapServers=pkc-***:9092 \
--set kafka.username=*** \
--set kafka.password=*** \
-f values-local.yaml
```
