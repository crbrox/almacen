package ctx

import (
	"fmt"
	"log"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type Ctx struct {
	TransID     string
	InfoLogger  Logger
	DebugLogger Logger
	Debug       bool
}

func (c *Ctx) Infof(format string, args ...interface{}) {
	f, a := c.addTransID(format, args)
	c.InfoLogger.Printf(f, a...)
}

func (c *Ctx) Debugf(format string, args ...interface{}) {
	if c.Debug {
		f, a := c.addTransID(format, args)
		if l, ok := c.DebugLogger.(*log.Logger); ok {
			s := fmt.Sprintf(f, a...)
			l.Output(2, s)
		} else {
			c.DebugLogger.Printf(f, a...)
		}
	}
}

func (c *Ctx) addTransID(format string, args []interface{}) (newFormat string, newArgs []interface{}) {
	newFormat = "[%s] " + format
	newArgs = make([]interface{}, 1, len(args)+1)
	newArgs[0] = c.TransID
	newArgs = append(newArgs, args...)
	return newFormat, newArgs
}
