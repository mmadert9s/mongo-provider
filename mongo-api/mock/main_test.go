package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testCase struct {
	name            string
	requestFunc     func(w http.ResponseWriter, r *http.Request)
	requestMethod   string
	requestEndpoint string
	requestBody     string
	requestHeaders  map[string]string
	responseCode    int
	responseBody    string
}

var validHeaders = map[string]string{"Authorization": "api-key-1"}
var invalidHeaders = map[string]string{}
var wrongKeyHeaders = map[string]string{"Authorization": "someotherkey"}

func TestList(t *testing.T) {
	testCases := []testCase{
		{
			"list mongos",
			listMongos,
			http.MethodGet,
			"/mongos",
			"",
			validHeaders,
			200,
			`[{"id":1,"name":"mongo1"},{"id":2,"name":"mongo2"},{"id":3,"name":"mongo3"},{"id":4,"name":"mongo4"}]`,
		},
		{
			"list mongos with bad api key",
			listMongos,
			http.MethodGet,
			"/mongos",
			"",
			wrongKeyHeaders,
			403,
			``,
		},
		{
			"list mongos with missing api key",
			listMongos,
			http.MethodGet,
			"/mongos",
			"",
			invalidHeaders,
			401,
			"Authorization header missing",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			execHttpTest(test, t)
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []testCase{
		{
			"get mongo",
			getMongo,
			http.MethodGet,
			"/mongos/1",
			"",
			validHeaders,
			200,
			`{"id":1,"name":"mongo1","regions":["eu","us"],"status":{"regions":{"eu":"ready","us":"ready"}},"connection":{"endpoint":"mongodb://t1m1.cloudprovider.com:27017","username":"admin","password":"password"}}`,
		},
		{
			"get mongo not found",
			getMongo,
			http.MethodGet,
			"/mongos/5",
			"",
			validHeaders,
			404,
			`Not Found`,
		},
		{
			"get mongo with bad api key",
			getMongo,
			http.MethodGet,
			"/mongos/2",
			"",
			wrongKeyHeaders,
			403,
			``,
		},
		{
			"get mongo with missing api key",
			getMongo,
			http.MethodGet,
			"/mongos/2",
			"",
			invalidHeaders,
			401,
			"Authorization header missing",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			execHttpTest(test, t)
		})
	}
}

func TestCreate(t *testing.T) {
	testCases := []testCase{
		{
			"create mongo with two regions",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mymongo", "regions":["us","eu"]}`,
			validHeaders,
			200,
			``,
		},
		{
			"create mongo with one region",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mymongo", "regions":["us"]}`,
			validHeaders,
			200,
			``,
		},
		{
			"create mongo with other region",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mymongo", "regions":["eu"]}`,
			validHeaders,
			200,
			``,
		},
		{
			"create mongo with invalid region",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mymongo", "regions":["en", "us"]}`,
			validHeaders,
			400,
			`Invalid region provided`,
		},
		{
			"create mongo with too many regions",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mymongo", "regions":["eu", "us", "eu"]}`,
			validHeaders,
			400,
			`Too many regions provided`,
		},
		{
			"create mongo with a taken name",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mongo1", "regions":["eu", "us"]}`,
			validHeaders,
			400,
			`Name already exists`,
		},
		{
			"create mongo with bad api key",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mongo1", "regions":["eu", "us"]}`,
			wrongKeyHeaders,
			403,
			"",
		},
		{
			"create mongo with missing api key",
			createMongo,
			http.MethodPost,
			"/mongos",
			`{"name": "mongo1", "regions":["eu", "us"]}`,
			invalidHeaders,
			401,
			"Authorization header missing",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			execHttpTest(test, t)
		})
	}
}

func TestPatch(t *testing.T) {
	testCases := []testCase{
		{
			"update mongo with two regions",
			updateMongo,
			http.MethodPatch,
			"/mongos/1",
			`{"regions":["us","eu"]}`,
			validHeaders,
			200,
			``,
		},
		{
			"update mongo with one region",
			updateMongo,
			http.MethodPatch,
			"/mongos/2",
			`{"regions":["eu"]}`,
			validHeaders,
			200,
			``,
		},
		{
			"update mongo with another region",
			updateMongo,
			http.MethodPatch,
			"/mongos/2",
			`{"regions":["us"]}`,
			validHeaders,
			200,
			``,
		},
		{
			"update mongo not found",
			updateMongo,
			http.MethodPatch,
			"/mongos/5",
			"",
			validHeaders,
			404,
			`Not Found`,
		},
		{
			"update mongo with invalid region",
			updateMongo,
			http.MethodPatch,
			"/mongos/3",
			`{"regions":["us", "en"]}`,
			validHeaders,
			400,
			`Invalid region provided`,
		},
		{
			"update mongo with too many regions",
			updateMongo,
			http.MethodPatch,
			"/mongos/2",
			`{"regions":["us", "eu", "eu"]}`,
			validHeaders,
			400,
			`Too many regions provided`,
		},
		{
			"update mongo with bad api key",
			updateMongo,
			http.MethodPatch,
			"/mongos/2",
			`{"regions":["us","eu"]}`,
			wrongKeyHeaders,
			403,
			"",
		},
		{
			"update mongo with missing api key",
			updateMongo,
			http.MethodPatch,
			"/mongos/2",
			`{"regions":["us","eu"]}`,
			invalidHeaders,
			401,
			"Authorization header missing",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			execHttpTest(test, t)
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []testCase{
		{
			"delete mongo",
			deleteMongo,
			http.MethodDelete,
			"/mongos/1",
			"",
			validHeaders,
			200,
			"",
		},
		{
			"delete mongo not found",
			deleteMongo,
			http.MethodDelete,
			"/mongos/5",
			"",
			validHeaders,
			404,
			`Not Found`,
		},
		{
			"delete mongo with bad api key",
			deleteMongo,
			http.MethodDelete,
			"/mongos/2",
			"",
			wrongKeyHeaders,
			403,
			``,
		},
		{
			"delete mongo with missing api key",
			deleteMongo,
			http.MethodDelete,
			"/mongos/2",
			"",
			invalidHeaders,
			401,
			"Authorization header missing",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			execHttpTest(test, t)
		})
	}
}

func execHttpTest(test testCase, t *testing.T) {
	req := httptest.NewRequest(test.requestMethod, test.requestEndpoint, strings.NewReader(test.requestBody))
	for k, v := range test.requestHeaders {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()

	test.requestFunc(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != test.responseCode {
		t.Errorf("expected statusCode to be %d, but got %d", test.responseCode, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected no error while reading body, but got %v", err)
	}

	if strings.TrimSpace(string(body)) != test.responseBody {
		t.Errorf("expected body to be `%s`, but got `%s`", test.responseBody, string(body))
	}
}
