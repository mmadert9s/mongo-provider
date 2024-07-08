# Mongo API Spec

## Model
A mongo instance is a set of MongoDB nodes running across one or multiple regions. There is only one instance per specified region. The set of regions can be updated, in which case a copying or migration of data occurs.

A mongo region (i.e. node) can be in one of the following states:
- migrating: If the node is currently being migrated from another region.
- unready: This can occur if the region's node is down for any reason other than migration.
- ready: If the region's node is running and fully migrated. 

Note: the status only contains information about regions in the current set of regions. For example, when updating the set of regions from {r1} to {r2}, the status will show `{"r2": "migrating"}` even if the node in r1 is still running during the migration. When updating the set of regions from {r1} to {r1, r2}, the status will show `{"r1": "ready", "r2":"migrating"}` since the node in r1 is still up and running and part of the set of desired regions.

The endpoint returned for a mongo instance points to one region at a time. It will only be reachable if at least one region in the applied set has the status `ready`.

A mongo resource has following schema:
```json
{
  "title": "mongo",
  "type": "object",
  "description": "A cluster of MongoDB nodes running across one or multiple regions",
  "properties": {
    "id": {
      "type": "integer",
      "description": "The mongo's id"
    },
    "name": {
      "type": "string",
      "description": "The mongo's name"
    },
    "regions": {
      "type": "array",
      "description": "The mongo's desired regions",
      "items": {
        "type": "string",
        "minItems": 1,
        "uniqueItems": true
      }
    },
    "status": {
      "type": "object",
      "description": "The mongo's status. Contains information about the status of each region",
      "properties": {
        "regions": {
          "type": "object",
          "additionalProperties": {
            "type": "string",
            "enum": ["unready", "migrating", "ready"]
          }
        }
      }
    },
    "connection": {
      "type": "object",
      "description": "The connection information",
      "properties": {
        "endpoint": {
          "type": "string",
          "description": "The endpoint for connecting to the database"
        },
        "username": {
          "type": "string",
          "description": "The user for connecting to the database"
        },
        "password": {
          "type": "string",
          "description": "The password for connecting to the database"
        }
      }
    }
  }
}
```

## Endpoint overview

| Method   | Path   | Description |
|----------|--------|-------------|
| GET | /mongos  | List mongo instances |
| GET | /mongos/{id}| Get a mongo instance |
| POST| /mongos| Create a mongo instance |
| PATCH | /mongos/{id}| Update a mongo instance's regions |
| DELETE | /mongos/{id} | Delete a mongo instance |

## General parameters

**Authorization**

The header ```Authorization: <API_KEY>``` is expected for all requests. The API key uniquely identifies tenants. If the key is missing or not valid, the response will be ```401: Unauthorized```

## Endpoints

### List mongo instances

```
GET /mongos
```

Returns the list of mongo instances belonging to the tenant using the following schema:

```json
{
  "title": "mongos",
  "type": "array",
  "description": "Short listing of mongo instances",
  "items": {
    "type": "object",
    "properties": {
      "id": {
        "type": "integer",
        "description": "The mongo's id"
      },
      "name": {
        "type": "string",
        "description": "The mongo's name"
      }
    }
  }
}
```

#### Responses

`200 OK`

Example response:
```json
[
  {
    "id": 1,
    "name": "mongo-instance-1"
  },
  {
    "id": 2,
    "name": "mongo-instance-2"
  }
]
```
Other responses:

`500 Internal Server Error` If an unexpected error occured.

---

### Get a mongo instance

```
GET /mongos/{id}
```

Returns the details and status of the mongo instance with the given `id` using the mongo schema.

#### Responses

`200 OK`

Example response:
```json
{
  "id": 1,
  "name": "mongo-instance-1",
  "regions": [
    "region1",
    "region2"
  ],
  "status": {
    "region1": "ready",
    "region2": "migrating"
  },
  "connection": {
    "endpoint": "mongodb://8f91a0b3.cloudprovider.com:27017",
    "username": "8f91a0b3",
    "password": "la59KI714gokAYZBc4eO"
  }
}
```

Other responses:

`404 Not Found` If the instance does not exist or if the instance does not belong to the tenant.

`500 Internal Server Error` If an unexpected error occured.

---

### Create a mongo instance

```
POST /mongos
```

Creates a mongo instance. The request body expects following schema:

```json
{
  "title": "mongodefinition",
  "type": "object",
  "description": "Definition of a mongo instance for creation",
  "properties": {
    "name": {
      "type": "string",
      "description": "The mongo's name"
    },
    "regions": {
      "type": "array",
      "description": "The mongo's desired regions",
      "items": {
        "type": "string",
        "minItems": 1,
        "uniqueItems": true
      }
    }
  }
}

```
#### Example Body

```json
{
  "name": "my-mongo-instance",
  "regions": ["region1", "region2"]
}
```

#### Responses

`201 Created`

`400 Bad Request` If the instance name already exists, is missing, or regions are invalid or missing, or the request is otherwise malformed.

`500 Internal Server Error` If an unexpected error occured.

---

### Update a mongo instance's regions

```
PATCH /mongo/{id}
```

Updates the regions of the mongo instance with the given `id`. The request body expects following schema:

```json
{
  "title": "mongoupdate",
  "type": "object",
  "description": "Update of a mongo instance's regions",
  "properties": {
    "regions": {
      "type": "array",
      "description": "The mongo's desired regions",
      "items": {
        "type": "string",
        "minItems": 1,
        "uniqueItems": true
      }
    }
  }
}

```
#### Example Body

```json
{
  "regions": ["region1", "region2", "region3"]
}
```

#### Responses

`200 OK`

`400 Bad Request` If regions are invalid or missing, or the request is otherwise malformed.

`404 Not Found` If the instance does not exist or if the instance does not belong to the tenant.

`500 Internal Server Error` If an unexpected error occured.

---

### Delete a mongo instance

```
DELETE /mongos/{id}
```

Deletes the mongo instance with the given `id`.

#### Responses

`200 OK`

`404 Not Found` If the instance does not exist or if the instance does not belong to the tenant.

`500 Internal Server Error` If an unexpected error occured.