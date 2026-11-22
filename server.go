package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type MensagemPrivada struct {
	destinatario string
	texto        string
}

type ConexaoCliente struct {
	apelido string
	tipo    string
	ch      chan<- string
}

var (
	entrada      = make(chan ConexaoCliente)
	saida        = make(chan ConexaoCliente)
	mensagensPub = make(chan string)
	mensagensPvt = make(chan MensagemPrivada)
	clientes     = make(map[string]ConexaoCliente)
)

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:4444")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Portas abertas!")
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-mensagensPub:
			for cli := range clientes {
				if clientes[cli].tipo == "USER" {
					clientes[cli].ch <- msg
				}
			}
			fmt.Println(msg)
		case msg := <-mensagensPvt:
			clientes[msg.destinatario].ch <- msg.texto
			fmt.Println(msg.texto)
		case conexao := <-entrada:
			clientes[conexao.apelido] = conexao
		case conexao := <-saida:
			delete(clientes, conexao.apelido)
			close(conexao.ch)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	ch := make(chan string)
	go clientWriter(conn, ch)

	contratoBytes := make([]byte, 1024)
	n, err := conn.Read(contratoBytes)
	if err != nil {
		fmt.Fprintln(conn, "Conexão Recusada:", err)
		return
	}
	contratoArr := strings.Fields(string(contratoBytes[:n]))

	tipo := contratoArr[0]
	if tipo != "USER" && tipo != "BOT" {
		fmt.Fprintln(conn, "Conexão Recusada:", "Cliente deve ser USER ou BOT")
		return
	}

	apelido := contratoArr[1]
	if err = validarApelido(apelido); err != nil {
		fmt.Fprintln(conn, "Conexão Recusada:", err)
		return
	}

	cc := ConexaoCliente{apelido, tipo, ch}

	ch <- "Conectado!"
	entrada <- cc
	mensagensPub <- mapearTipoCliente(cc) + " @" + cc.apelido + " acabou de entrar"

	input := bufio.NewScanner(conn)

SCAN_LOOP:
	for input.Scan() {
		txt := input.Text()
		tipoMsg, msg, dst := mapearMensagem(txt)

		switch tipoMsg {
		case "MSG_PUB":
			mensagensPub <- "@" + cc.apelido + " disse: " + msg
		case "MSG_PVT":
			mensagensPvt <- MensagemPrivada{dst, "@" + cc.apelido + " disse em privado: " + msg}
		case "TROCAR_APELIDO":
			apelidoAntigo := cc.apelido
			cc.apelido = msg
			clientes[cc.apelido] = cc
			delete(clientes, apelidoAntigo)
			mensagensPub <- mapearTipoCliente(cc) + " @" + apelidoAntigo + " agora é @" + cc.apelido
		case "SAIR":
			break SCAN_LOOP
		case "ERRO":
			mensagensPvt <- MensagemPrivada{cc.apelido, "Comando inválido: " + msg}
		case "DESCONHECIDO":
			mensagensPvt <- MensagemPrivada{cc.apelido, "Comando desconhecido"}
		default:
			continue
		}
	}
	if err := input.Err(); err != nil {
		fmt.Println("Erro no Scanner:", err)
	}

	saida <- cc
	mensagensPub <- mapearTipoCliente(cc) + " @" + cc.apelido + " saiu"
}

func validarApelido(apelido string) error {
	apelido = strings.TrimSpace(apelido)
	if len(apelido) == 0 || len(apelido) > 32 {
		return errors.New("Apelido deve ter entre 1 e 32 caracteres")
	}

	for i := range apelido {
		if apelido[i] == ' ' {
			return errors.New("Apelido não pode conter espaços")
		}
	}

	clienteExistente := clientes[apelido]
	if clienteExistente.apelido == apelido {
		return errors.New("Apelido já existe")
	}

	return nil
}

func mapearMensagem(msg string) (tipo, msgFormatada, destinatario string) {
	args := strings.Fields(msg)
	if len(args) == 0 {
		return "DESCONHECIDO", "", ""
	}

	switch args[0] {
	case "\\changenick":
		if len(args) < 2 {
			return "ERRO", "Informe o novo apelido", ""
		}
		err := validarApelido(args[1])
		if err != nil {
			return "ERRO", err.Error(), ""
		}
		return "TROCAR_APELIDO", args[1], ""
	case "\\msg":
		if len(args) < 2 {
			return "ERRO", "Informe a mensagem a ser enviada", ""
		}
		if args[1][0] == '@' {
			apelido := args[1][1:]
			destinatario := clientes[apelido]
			if destinatario.apelido != apelido {
				return "ERRO", "Usuário ou Bot @" + apelido + " não encontrado", ""
			}
			if len(args) < 3 {
				return "ERRO", "Informe a mensagem a ser enviada para @" + destinatario.apelido, ""
			}
			return "MSG_PVT", strings.Join(args[2:], " "), destinatario.apelido
		}
		return "MSG_PUB", strings.Join(args[1:], " "), ""
	case "\\exit":
		return "SAIR", "", ""
	default:
		break
	}

	return "DESCONHECIDO", "", ""
}

func mapearTipoCliente(cc ConexaoCliente) string {
	switch cc.tipo {
	case "USER":
		return "Usuário"
	case "BOT":
		return "Bot"
	}
	return ""
}
