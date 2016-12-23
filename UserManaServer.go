// UserManaServer project UserManaServer.go
package main

import (
	"UserManaServer/process"
	_ "UserManaServer/routers"
	"github.com/astaxie/beego"
	// "io/ioutil"
	"log"
	"net"
)

func main() {
	//设置监听端口9999
	go func() {
		listen, err := net.Listen("tcp", ":9999")
		if err != nil {
			log.Fatalln(err)
			return
		}
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
				continue
			}
			go process.HandleConn(conn)
		}
	}()
	beego.Run()
}
