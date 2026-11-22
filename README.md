# EP 3 — Chat em Go

## Integrantes

- Cauã Mendes
- Stephany França

## Arquivos

- `server.go` — servidor TCP e gerenciamento das conexões e mensagens.
- `client.go` — cliente de terminal para comunicação entre usuários.
- `bot.go` — bot que inverte mensagens privadas recebidas.

## Comunicação

A comunicação entre os programas utiliza TCP, com host `localhost` e porta `4444`.

Cada programa deve ser executado em um terminal separado. O servidor deve ser iniciado primeiro.

## Como executar

**1. Servidor:**

```bash
go run server.go
```

**2. Cliente Alice:**

```bash
go run client.go Alice
```

**3. Cliente Bob:**

```bash
go run client.go Bob
```

Alice e Bob são exemplos de apelidos e podem ser substituídos por outros nomes.

**4. Bot:**

```bash
go run bot.go
```

Se o cliente for executado sem apelido, ele gera um automaticamente:

```bash
go run client.go
```

## Comandos do cliente

```text
\msg Olá pessoal          mensagem pública
\msg @Bob Oi Bob          mensagem privada
\changenick GatoPreto     troca de apelido
\exit                    sai do chat
```

## Bot

O bot se conecta ao servidor como um usuário do tipo `BOT`, utilizando o apelido fixo `Inversor`.

Ele não recebe mensagens públicas, apenas mensagens privadas.

Sua função é inverter o texto recebido e enviar a resposta ao remetente.

Exemplo:

```text
\msg @Inversor abcdef
```

Saída no terminal do bot:

```text
Mensagem de: @Alice
Recebi: abcdef
Resposta: fedcba
```