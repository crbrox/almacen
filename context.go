package almacen

import (
	"log"
	"os"

	"github.com/crbrox/almacen/ctx"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

var (
	DebugLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	InfoLogger  = log.New(os.Stderr, "", log.LstdFlags)
)

type context struct {
	params  httprouter.Params
	session *mgo.Session
	input   interface{}
	ctx.Ctx
}

func NewContext() *context {
	return &context{
		Ctx: ctx.Ctx{
			DebugLogger: DebugLogger,
			InfoLogger:  InfoLogger}}
}
