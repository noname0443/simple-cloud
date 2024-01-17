package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

const (
	MYSQL_CONN  = "%s:3306"
	SOCKET_CONN = ":8080"
)

func handleConnection(conn net.Conn, mysql_addr string) {
  fmt.Println(fmt.Sprintf(MYSQL_CONN, mysql_addr))
	mysql, err := net.Dial("tcp", fmt.Sprintf(MYSQL_CONN, mysql_addr))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	defer mysql.Close()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	go io.Copy(conn, mysql)
	io.Copy(mysql, conn)
}

func handleConnectionWithBash(conn net.Conn) {
	out, err := exec.Command("/bin/bash", "-c", "/sample.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	handleConnection(conn, string(out))
}

func main() {
	ln, err := net.Listen("tcp", SOCKET_CONN)
	for err != nil {
		ln, err = net.Listen("tcp", SOCKET_CONN)
		fmt.Printf("%s", err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		go handleConnectionWithBash(conn)
	}
}
