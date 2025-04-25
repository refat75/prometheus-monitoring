# Monitor Go Server Using Prometheus

## Start & Stop Application
```shell
docker compose up
docker compose down
```

## Access Endpoints

```shell
http://localhost:8080/ # Server
http://localhost:9090/ # Prometheus 
http://localhost:3000/ # Grafana
```


## Reference
[Golang Application monitoring using Prometheus](https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/)