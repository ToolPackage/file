package main

import (
	"github.com/ToolPackage/fse/service"
	"github.com/ToolPackage/fse/tx"
	"log"
	"net"
)

func main() {
	service.Init()

	// s := server.New()
	// s.Start()

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9330")
	listener, _ := net.ListenTCP("tcp", addr)
	defer listener.Close()
	log.Println("fse server started")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("new client connected:", conn.LocalAddr())
		c := tx.NewChannel(conn, conn)
		go c.Process()
	}
}
