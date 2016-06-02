package main

import (
	"net/http"
	"os"

	"github.com/crbrox/almacen"
	"github.com/julienschmidt/httprouter"

	"github.com/crbrox/almacen/ctx"
)

const (
	ExitStatusConfig = iota + 2
	ExitStatusStore
	ExitStatusServer
)

func main() {
	var err error
   _ = "breakpoint"
	cB := &ctx.Ctx{TransID: "n/a",
		DebugLogger: almacen.DebugLogger,
		InfoLogger:  almacen.InfoLogger}

	c, err := almacen.LoadConfig("./config.json")
	if err != nil {
		cB.Infof("config: %v ", err)
		os.Exit(ExitStatusConfig)
	}
	cB.Infof("config loaded %#v", c)
	mes := &almacen.MongoEntityStore{}
	err = mes.Start(c)
	if err != nil {
		cB.Infof("store start: %v", err)
		os.Exit(ExitStatusStore)
	}
	cB.Infof("store started")
	defer mes.Stop()

	almacen.SetStore(mes)

	router := httprouter.New()

	almacen.AddRoutes(router)

	cB.Infof("starting http server %v", c.Address)
	err = http.ListenAndServe(c.Address, router)
	if err != nil {
		cB.Infof("start server: %v", err)
		os.Exit(ExitStatusServer)
	}

}
