package main

import (
	log "github.com/Luncert/slog"
	"github.com/ToolPackage/fse/server"
)

func main() {
	log.InitLogger("cmd/conf/log.yml")

	s := server.New()
	s.Start()
}
