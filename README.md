# Rate Limiter em Go
Este projeto implementa um Rate Limiter simples em Go, que limita o número de requisições que um cliente pode fazer a um servidor em um determinado período de tempo. Ele utiliza um sistema de bloqueio temporário para impedir abusos.

## Requisitos
Antes de executar o projeto, certifique-se de que você tem as seguintes dependências instaladas:

* Go 1.18+ (para compilar e rodar o projeto);
* Dependências do projeto que serão instaladas com o comando:
```shell script
go mod
```

## Como Executar
Siga os passos abaixo para executar o Rate Limiter localmente.

* Clone o repositório para o seu ambiente local:
```shell script
git clone https://github.com/seu-usuario/rate-limiter.git
```
```shell script
cd rate-limiter/cmd
```
* Instalar Dependências

    O Go usa um sistema de módulos para gerenciar dependências. Para garantir que todas as dependências do projeto sejam baixadas, execute o comando abaixo:

```shell script
go mod tidy
```

* Para rodar o servidor, use o comando:
```shell script
go run main.go
```
> O servidor estará disponível em http://localhost:8080. Você pode testar as funcionalidades de rate limiting enviando requisições HTTP. Aqui estão alguns exemplos de como você pode fazer isso:

* Teste com curl:
```shell script
curl --location 'http://localhost:8080/' \
--header 'API_KEY: my-token'
```

* Esse endpoint estará protegido pelo Rate Limiter. Se você fizer várias requisições rápidas para o mesmo IP ou token (usando o cabeçalho API_KEY), as requisições serão bloqueadas após o limite de requisições ser atingido.

* Teste com Ferramenta de API (Postman, Insomnia, etc):
Envie requisições GET para http://localhost:8080/ e observe o comportamento do Rate Limiter.

Passo 5: Parar o Servidor
Para parar o servidor, pressione Ctrl+C no terminal onde o servidor está rodando.

Como Funciona
O Rate Limiter utiliza um sistema de controle de requisições baseado em uma chave identificadora (como o IP do cliente ou o token API_KEY). Ele conta o número de requisições feitas dentro de um período de tempo e bloqueia o cliente por um tempo configurado após atingir o limite.

Configuração do Rate Limiter
O Rate Limiter pode ser configurado na inicialização do servidor. Aqui está um exemplo de como você pode configurar o limitador no código:

## Estrutura do Projeto
A estrutura do projeto é organizada da seguinte forma:

```
rate-limiter/
├── cmd/                            # Contém o arquivo de entrada (main.go)
│   ├── .env
│   └── main.go
├── config/                         # Contém as configurações do projeto
│   └── config.go
├── internal/                       # Contém a lógica interna do projeto
│   ├── rate_limiter/               # Implementação do Rate Limiter
│   │   ├── rate_limiter.go
│   │   ├── rate_limiter_test.go  
│   ├── middleware/                 # Middleware para interceptação das requisições
│   │   ├── middleware.go
│   │   └── middleware_test.go
│   └── persistence/                # Implementação da persistência (ex.: Redis)
│       ├── persistence.go
│       └── persistence_test.go     # Testes para a persistência
├── docker-compose.yml   
├── go.mod                          # Dependências do Go
└── README.md           
```

## Licença
Este projeto está licenciado sob a [Licença MIT](LICENSE).
