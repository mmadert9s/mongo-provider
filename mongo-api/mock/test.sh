#!/bin/sh

# GET /mongos
curl -i -H "Authorization: api-key-1" localhost:8080/mongos
curl -i -H "Authorization: api-key-2" localhost:8080/mongos
curl -i -H "" localhost:8080/mongos


# GET /mongos/{id}
curl -i -H "Authorization: api-key-1" localhost:8080/mongos/1
curl -i -H "Authorization: api-key-1" localhost:8080/mongos/7
curl -i -H "Authorization: api-key-2" localhost:8080/mongos/1

# POST /mongos
curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mymongo", "regions":["us","eu"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mymongo", "regions":["us"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mymongo", "regions":["eu"]}' localhost:8080/mongos

curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mymongo", "regions":["us","eu","eu"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mymongo", "regions":["us","en"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X POST -d '{"name": "mongo1", "regions":["us","eu"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-2" -X POST -d '{"name": "mymongo", "regions":["us","eu"]}' localhost:8080/mongos

# PATCH /mongos/{id}
curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["us","eu"]}' localhost:8080/mongos/3
curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["us"]}' localhost:8080/mongos/3
curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["eu"]}' localhost:8080/mongos/3

curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["us","eu","eu"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["us","en"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-1" -X PATCH -d '{"regions":["us","eu"]}' localhost:8080/mongos
curl -i -H "Authorization: api-key-2" -X PATCH -d '{"regions":["us","eu"]}' localhost:8080/mongos/3

# DELETE /mongos/{id}
curl -i -H "Authorization: api-key-1" -X DELETE localhost:8080/mongos/1
curl -i -H "Authorization: api-key-1" -X DELETE localhost:8080/mongos/7
curl -i -H "Authorization: api-key-2" -X DELETE localhost:8080/mongos/1


