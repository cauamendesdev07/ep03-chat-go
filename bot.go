package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4444")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("BOT Inversor"))
	if err != nil {
		log.Fatal(err)
	}

	input := bufio.NewScanner(conn)

	for input.Scan() {
		txt := input.Text()
		remetente, chars := formatarMensagem(txt)
		if remetente == "" {
			continue
		}
		inverso := inverter(chars)

		fmt.Println("Mensagem de: @" + remetente)
		fmt.Println("Recebi:", chars)
		fmt.Println("Resposta:", inverso)

		_, err = conn.Write([]byte("\\msg @" + remetente + " " + inverso + "\n"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func formatarMensagem(txt string) (remetente, chars string) {
	args := strings.Fields(txt)
	if len(args) < 5 {
		return "", txt
	}
	return args[0][1:], strings.Join(args[4:], " ")
}

func inverter(msg string) string {
	chars := []rune(msg)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}
