# Mongo API

The mongoapi implements a mock API for the specification defined in mongoapi.md.
It listens on port 8080 and returns mock data.

## Run

### Locally

Prerequisites:
- go 1.22

```
go run main.go
```

### In k8s

This will pull the image from `docker.io/mmadert9s/mongoapi`

Prerequisites:
- kubectl pointing to an available k8s cluster (e.g. kind)

```
kubectl apply -f ./deploy/mongoapi.yaml
kubectl port-forward svc/mongoapi 8080:8080
```

# Use

Provide the api key `api-key-1` in the `Authorization` header of the requests.

Examples:
```
curl -i -H "Authorization: api-key-1" localhost:8080/mongos
curl -i -H "Authorization: api-key-1" localhost:8080/mongos/1
# These won't actually create, update or delete anything.
curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mymongo", "regions":["us","eu"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["us","eu"]}' localhost:8080/mongos/3
curl -i -H "Authorization: api-key-1" -X DELETE localhost:8080/mongos/1
```