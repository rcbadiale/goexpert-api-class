# goexpert-api-class

API desenvolvida durante as aulas do curso Go Expert.

Se trata de uma API de exemplo para explicação/demonstração da implementação de
uma API em Go.

Obs.: Essa é uma aplicação de exemplo apenas, na qual existem diversos pontos
que precisam de refatoração e melhorias.

## Execução

1. Ir para a pasta `cmd/server` para execução do projeto
```shell
cd cmd/server
```

2. Criar um arquivo `.env` com as configurações abaixo
```shell
JWT_SECRET=<chave secreta>
JWT_EXPIRESIN=300
```
3. Executar o projeto
```shell
go run main.go
```

## Gerar documentação

Instalar o pacote `swag` com o comando abaixo.
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

Executar o comando para gerar os documentos na raíz do projeto.
```shell
swag init -g cmd/server/main.go
```

## Execução dos testes

Executar na raíz do projeto o comando para executar todos os testes em subpastas.
```shell
go test ./...
```

Para gerar a cobertura e exibir os relatórios utilizar os comandos na pasta raíz.
```shell
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Requisições

Na pasta `test` existem os arquivos `*.http` que definem as requisições
possíveis na aplicação.

Para executar essas requisições é necessário instalar a extensão
[REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) no VSCode.
