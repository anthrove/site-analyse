# E621 Analysis

## Local Development

To start running dev requirements locally you can fire up a local S3 storage and Prometheus instance using
the deployments compose file:

```sh
docker compose -f deployments/docker-compse.yml up minio prometheus prom-pushgateway
```

To compile and run the go project itself and trigger a run of the analysis system use:
```sh
go run cmd/analyse/main.go
```

If you also want to (re)load a *.env* file for configs you can use this oneliner instead:
```sh
set -a && source .env && set +a && go run cmd/analyse/main.go
```