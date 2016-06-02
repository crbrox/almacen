package almacen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var MaxConcurrentConn = 100
var IncomingReqSem = make(chan struct{}, MaxConcurrentConn)

func H(f func(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var ctx = NewContext()
		ctx.params = params

		switch req.Header.Get("x-trace") {
		case "TRUE", "true", "True", "ON", "on", "On":
			ctx.Debug = true
			ctx.Debugf("x-trace is on")
		}

		AddTransId(req, ctx)

		defer req.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		ctx.Debugf("incoming request: %s", &customReq{req})

		select {
		case IncomingReqSem <- struct{}{}:
		default:
			respondErr(w, &Error{message: "too many concurrent requests ", statusCode: http.StatusServiceUnavailable})
			return
		}
		defer func() { <-IncomingReqSem }()

		// decoding input object
		var object interface{}
		err := json.NewDecoder(req.Body).Decode(&object)
		if err != nil && err != io.EOF {
			ctx.Debugf("error encoding input: %v", err)
			respondErr(w, &Error{message: "json parsing: " + err.Error(), statusCode: http.StatusBadRequest})
			return
		}
		ctx.input = object
		ctx.Debugf("incoming object: %v (%T)", object, object)
		if store, ok := store.(*MongoEntityStore); ok {
			session := store.session.Copy()
			defer session.Close()
			ctx.session = session
		}

		// execute logic
		obj, err := f(ctx, w, req)

		if err != nil {
			ctx.Debugf("error processing request: %v", err)
			respondErr(w, err)
			return
		}

		// encoding response object
		ctx.Debugf("returning object: %#v (%T)", obj, obj)
		if obj != nil {
			err = json.NewEncoder(w).Encode(obj)
			if err != nil {
				ctx.Infof("error encoding response: %v", err)
				respondErr(w, &Error{message: "json encoding: " + err.Error(), statusCode: http.StatusInternalServerError})
			}
			return
		}
	}
}

func respondErr(w http.ResponseWriter, err error) {
	if localErr, ok := err.(*Error); ok {
		w.WriteHeader(localErr.statusCode)
		fmt.Fprintf(w, "%q\n", localErr.message)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%q\n", err.Error())

}

func AddTransId(req *http.Request, c *context) {
	const transIDHeaderField = "x-transid"

	var transID = req.Header.Get(transIDHeaderField)
	c.Debugf("received trans id: %q", transID)
	if transID == "" {
		// Temporary, to be replaced ...
		f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
		b := make([]byte, 16)
		f.Read(b)
		f.Close()
		transID = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		c.Debugf("generated trans id: %q", transID)
	}
	c.TransID = transID
}

// To avoid a conversion to string if not necessary
type customReq struct{ *http.Request }

func (r customReq) String() string {
	return fmt.Sprintf(
		"method=%q URL=%q header=%v remote=%v",
		r.Method, r.URL, r.Header, r.RemoteAddr)
}
