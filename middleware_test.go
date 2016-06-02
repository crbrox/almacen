package almacen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestRespondErrorGeneric(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := errors.New("plain error")
	respondErr(recorder, err)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("status code: wanted %d, got %d", recorder.Code, http.StatusInternalServerError)
	}
	msg := recorder.Body.String()
	wanted := fmt.Sprintf("%q\n", err)
	if msg != wanted {
		t.Errorf("body: wanted %q, got %q", wanted, msg)
	}
}

func TestRespondErrorPackage(t *testing.T) {
	for _, err := range []*Error{
		{statusCode: http.StatusInternalServerError, message: "internal server error"},
		{statusCode: http.StatusNotFound, message: "not found"},
		{statusCode: http.StatusBadRequest, message: "bad request, maybe"},
		{statusCode: http.StatusConflict, message: "conflict"},
		{statusCode: http.StatusBadGateway, message: "bad gateway"},
		ErrIdNotString,
		ErrNotFound,
		ErrExisting,
		ErrTooMany,
		ErrTraversingObject,
		ErrObjectExpected,
	} {
		recorder := httptest.NewRecorder()

		respondErr(recorder, err)
		if recorder.Code != err.statusCode {
			t.Errorf("status code: wanted %d, got %d", recorder.Code, err.statusCode)
		}
		msg := recorder.Body.String()
		wanted := fmt.Sprintf("%q\n", err.message)
		if msg != wanted {
			t.Errorf("body: wanted %q, got %q", wanted, msg)
		}
	}
}

func TestMiddlewareInvalidJSON(t *testing.T) {

	SetStore(NewMemStore())
	dummyFunc := func(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
		return nil, nil
	}
	recorder := httptest.NewRecorder()

	h := H(dummyFunc)

	buffer := bytes.NewBufferString("not a JSON parseable string")
	request, err := http.NewRequest("GET", "ruta", buffer)
	if err != nil {
		t.Fatal(err)
	}

	h(recorder, request, httprouter.Params{})

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status code: wanted %d, got %d", http.StatusBadRequest, recorder.Code)
	}

}

func TestMiddlewareMongoStore(t *testing.T) {
	mes := &MongoEntityStore{}
	err := mes.Start(&Config{MongoURL: "localhost"})
	if err != nil {
		t.Fatal(err)
	}
	defer mes.Stop()

	SetStore(mes)
	var c = &context{}
	dummyFunc := func(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
		// grab the context for later checks
		c = ctx
		return nil, nil
	}
	recorder := httptest.NewRecorder()

	h := H(dummyFunc)

	request, err := http.NewRequest("GET", "ruta", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}

	h(recorder, request, httprouter.Params{})

	if recorder.Code != http.StatusOK {
		t.Errorf("status code: wanted %d, got %d", http.StatusOK, recorder.Code)
	}

	if c.session == nil {
		t.Errorf("session is nil")
	}
}

func TestMiddlewareEcho(t *testing.T) {
	var err error
	SetStore(NewMemStore())
	object := map[string]interface{}{
		"integer": 2.0,
		"string":  "STRING",
		"array":   []interface{}{1.0, "dos", 3.0},
		"object":  map[string]interface{}{"a": "a", "b": 2.1},
	}
	echoFunc := func(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
		return ctx.input, nil
	}
	recorder := httptest.NewRecorder()

	h := H(echoFunc)

	oJSON, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewBuffer(oJSON)

	request, err := http.NewRequest("GET", "ruta", body)
	if err != nil {
		t.Fatal(err)
	}
	params := httprouter.Params{{"p1", "v1"}, {"p2", "v2"}}
	h(recorder, request, params)

	var retrievedObj map[string]interface{}
	err = json.NewDecoder(recorder.Body).Decode(&retrievedObj)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(object, retrievedObj) {
		t.Errorf("retrieved object: wanted %#v, got %#v", object, retrievedObj)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("content-type: wanted %s, got %s", contentType, "application/json")
	}

}


type testingNotJSONableObject struct{}

func (testingNotJSONableObject) MarshalJSON() ([]byte, error) {
	return nil, errors.New("JSON not supported")
}

func TestMiddlewareInvalidJSONInResponse(t *testing.T) {
	var err error
	SetStore(NewMemStore())

   notJSONFunc := func(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
		return testingNotJSONableObject{}, nil
	}
	recorder := httptest.NewRecorder()

	h := H(notJSONFunc)

	request, err := http.NewRequest("GET", "ruta",  &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	params := httprouter.Params{}
	h(recorder, request, params)

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("content-type: wanted %s, got %s", contentType, "application/json")
	}

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("status code: wanted %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

}

func TestMiddlewareErrorProcessing(t *testing.T) {
	const errText =  "a error from f"
	var err error
	SetStore(NewMemStore())

	errFunc := func(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
		return nil, &Error{statusCode: http.StatusBadRequest, message: errText}
	}
	recorder := httptest.NewRecorder()

	h := H(errFunc)

	request, err := http.NewRequest("GET", "ruta",  &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	params := httprouter.Params{}
	h(recorder, request, params)

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("content-type: wanted %s, got %s", contentType, "application/json")
	}

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status code: wanted %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	got := recorder.Body.String()
	wanted := fmt.Sprintf("%q\n", errText)
	if got != wanted {
		t.Errorf("error message: wanted %+q, got %+q", wanted, got)
	}

}
