package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	tenant1 string = "api-key-1"
)

type Region string

const (
	us Region = "us"
	eu Region = "eu"
)

type RegionStatus string

const (
	statusUnready   RegionStatus = "unready"
	statusMigrating RegionStatus = "migrating"
	statusReady     RegionStatus = "ready"
)

type MongoIdentity struct {
	tenant string
	Id     int    `json:"id"`
	Name   string `json:"name"`
}

type MongoStatus struct {
	Regions map[Region]RegionStatus `json:"regions"`
}

type MongoConnection struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Mongo struct {
	MongoIdentity `json:",inline"`
	Regions       []Region        `json:"regions"`
	Status        MongoStatus     `json:"status"`
	Connection    MongoConnection `json:"connection"`
}

type MongoUpdate struct {
	Regions []Region `json:"regions"`
}

// -------------------------------- mongo 1 ---------------------------------------
// An instance with nodes running in two regions
var mongo1id = MongoIdentity{
	tenant: tenant1,
	Id:     1,
	Name:   "mongo1",
}

var mongo1status = MongoStatus{
	Regions: map[Region]RegionStatus{
		us: statusReady,
		eu: statusReady,
	},
}

var mongo1connection = MongoConnection{
	Endpoint: "mongodb://t1m1.cloudprovider.com:27017",
	Username: "admin",
	Password: "password",
}

var mongo1 = Mongo{
	MongoIdentity: mongo1id,
	Regions:       []Region{"eu", "us"},
	Status:        mongo1status,
	Connection:    mongo1connection,
}

// -------------------------------- mongo 2 ---------------------------------------
// An instance that is currently being migrated from us to eu
var mongo2id = MongoIdentity{
	tenant: tenant1,
	Id:     2,
	Name:   "mongo2",
}

var mongo2status = MongoStatus{
	Regions: map[Region]RegionStatus{
		eu: statusMigrating,
	},
}

var mongo2connection = MongoConnection{
	Endpoint: "mongodb://t1m2.cloudprovider.com:27017",
	Username: "admin",
	Password: "password",
}

var mongo2 = Mongo{
	MongoIdentity: mongo2id,
	Regions:       []Region{"eu"},
	Status:        mongo2status,
	Connection:    mongo2connection,
}

// -------------------------------- mongo 3 ---------------------------------------
// An unready instance
var mongo3id = MongoIdentity{
	tenant: tenant1,
	Id:     3,
	Name:   "mongo3",
}

var mongo3status = MongoStatus{
	Regions: map[Region]RegionStatus{
		us: statusUnready,
	},
}

var mongo3connection = MongoConnection{
	Endpoint: "mongodb://t1m3.cloudprovider.com:27017",
	Username: "admin",
	Password: "password",
}

var mongo3 = Mongo{
	MongoIdentity: mongo3id,
	Regions:       []Region{"us"},
	Status:        mongo3status,
	Connection:    mongo3connection,
}

// -------------------------------- mongo 4 ---------------------------------------
// An instance that is currently being migrated from eu to us,eu
var mongo4id = MongoIdentity{
	tenant: tenant1,
	Id:     4,
	Name:   "mongo4",
}

var mongo4status = MongoStatus{
	Regions: map[Region]RegionStatus{
		us: statusMigrating,
		eu: statusReady,
	},
}

var mongo4connection = MongoConnection{
	Endpoint: "mongodb://t1m4.cloudprovider.com:27017",
	Username: "admin",
	Password: "password",
}

var mongo4 = Mongo{
	MongoIdentity: mongo4id,
	Regions:       []Region{"eu", "us"},
	Status:        mongo4status,
	Connection:    mongo4connection,
}

func main() {
	http.HandleFunc("GET /mongos", func(w http.ResponseWriter, r *http.Request) {
		listMongos(w, r)
	})

	http.HandleFunc("GET /mongos/{id}", func(w http.ResponseWriter, r *http.Request) {
		getMongo(w, r)
	})

	http.HandleFunc("POST /mongos", func(w http.ResponseWriter, r *http.Request) {
		createMongo(w, r)
	})

	http.HandleFunc("PATCH /mongos/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateMongo(w, r)
	})

	http.HandleFunc("DELETE /mongos/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteMongo(w, r)
	})

	http.ListenAndServe(":8080", nil)
}

func listMongos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if !checkAuth(w, r) {
		return
	}

	mongos := []MongoIdentity{
		mongo1id, mongo2id, mongo3id, mongo4id,
	}

	json, err := json.Marshal(mongos)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func getMongo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if !checkAuth(w, r) {
		return
	}

	id, valid := idFromURL(w, r)
	if !valid {
		return
	}

	var mongo Mongo
	switch id {
	case 1:
		mongo = mongo1
	case 2:
		mongo = mongo2
	case 3:
		mongo = mongo3
	case 4:
		mongo = mongo4
	default:
		{
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	}

	json, err := json.Marshal(mongo)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func createMongo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if !checkAuth(w, r) {
		return
	}

	var mongo Mongo
	err := json.NewDecoder(r.Body).Decode(&mongo)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	for _, existing := range []Mongo{mongo1, mongo2, mongo3, mongo4} {
		if mongo.Name == existing.Name {
			http.Error(w, "Name already exists", http.StatusBadRequest)
			return
		}
	}

	if len(mongo.Regions) < 1 {
		http.Error(w, "No region provided", http.StatusBadRequest)
		return
	}

	if len(mongo.Regions) > 2 {
		http.Error(w, "Too many regions provided", http.StatusBadRequest)
		return
	}

	for _, region := range mongo.Regions {
		if region != eu && region != us {
			http.Error(w, "Invalid region provided", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func updateMongo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if !checkAuth(w, r) {
		return
	}

	id, valid := idFromURL(w, r)
	if !valid {
		return
	}

	found := false
	for _, existing := range []Mongo{mongo1, mongo2, mongo3, mongo4} {
		if id == existing.Id {
			found = true
		}
	}

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var mongoupdate MongoUpdate
	err := json.NewDecoder(r.Body).Decode(&mongoupdate)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	if len(mongoupdate.Regions) < 1 {
		http.Error(w, "No region provided", http.StatusBadRequest)
		return
	}

	if len(mongoupdate.Regions) > 2 {
		http.Error(w, "Too many regions provided", http.StatusBadRequest)
		return
	}

	for _, region := range mongoupdate.Regions {
		if region != eu && region != us {
			http.Error(w, "Invalid region provided", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func deleteMongo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if !checkAuth(w, r) {
		return
	}

	id, valid := idFromURL(w, r)
	if !valid {
		return
	}

	found := false
	for _, existing := range []Mongo{mongo1, mongo2, mongo3, mongo4} {
		if id == existing.Id {
			found = true
		}
	}

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// idFromURL retrieves and checks the id from the request URL. If it is invalid or an error occured,, it writes to the response and returns false.
func idFromURL(w http.ResponseWriter, r *http.Request) (int, bool) {
	paramId := strings.TrimPrefix(r.URL.Path, "/mongos/")
	if strings.Contains(paramId, "/") {
		http.Error(w, "Not Found", http.StatusBadRequest)
	}

	if paramId == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return 0, false
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Id parameter malformed", http.StatusBadRequest)
		return 0, false
	}

	return id, true
}

// checkAuth checks the Authorization header. If it is invalid or an error occured, it writes to the response and returns false.
func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	tenant := r.Header.Get("Authorization")
	if tenant == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return false
	}

	if tenant != tenant1 {
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}
