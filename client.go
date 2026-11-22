package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	username := "user_" + GerarID()
	if len(os.Args) >= 2 {
		username = os.Args[1]
	}
	conn, err := net.Dial("tcp", "localhost:4444")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msgContrato := "USER " + username
	_, err = conn.Write([]byte(msgContrato))
	if err != nil {
		log.Fatal(err)
	}

	fim := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("fim")
		fim <- struct{}{}
	}()
	go func() { mustCopy(conn, os.Stdin) }()
	<-fim
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func GerarID() string {
	nanosegundos := time.Now().UnixNano()
	return fmt.Sprintf("%d", nanosegundos)
}
