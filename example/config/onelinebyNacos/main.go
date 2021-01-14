package main

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	_ "github.com/edana-dev/xjupiter/pkg/datasource/nacos"
)

// add following config to your nacos server
/*
[people]
    name = "jupiter"
[jupiter.logger.default]
    debug = true
    enableConsole = true
[jupiter.server.http]
    port = 9090
[jupiter.server.governor]
    enable = false
    host = "0.0.0.0"
    port = 9246
*/

//  go run main.go --config="nacos://ip:port?dataId=XXXXX&group=DEFAULT_GROUP&tenant=XXXXX&scheme=http/https&level=debug|info|warn"
func main() {
	app := NewEngine()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.printConfig,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func (s *Engine) printConfig() error {
	xlog.DefaultLogger = xlog.StdConfig("default").Build()
	peopleName := conf.GetString("people.name")
	xlog.Info("people info", xlog.String("name", peopleName), xlog.String("type", "onelineByNacos"))
	return nil
}
