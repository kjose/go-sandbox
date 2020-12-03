package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprintln(conn, "POST /index.html HTTP/1.1")
	fmt.Fprintln(conn, "Host: client-go-test")
	fmt.Fprintln(conn, "")

	bs, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(bs))

}
