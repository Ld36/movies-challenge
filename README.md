# Movies Challenge

## Sobre o projeto

Este projeto foi desenvolvido como solução para o desafio técnico proposto, utilizando uma arquitetura baseada em microsserviços.

A aplicação expõe uma API REST para gerenciamento de filmes através de um API Gateway. O Gateway se comunica com um microsserviço responsável pelas regras de negócio utilizando gRPC e Protocol Buffers. Os dados são persistidos em MongoDB.

Além dos requisitos obrigatórios do desafio, foram adicionadas alguns itens para melhor immpacto do desafio, como paginação, documentação Swagger, logs estruturados, Request ID, Graceful Shutdown e Docker Compose.

---

# Arquitetura

```
                    HTTP/REST
                         │
                         ▼
                +------------------+
                |   API Gateway    |
                |      (Gin)       |
                +------------------+
                         │
                       gRPC
                         │
                         ▼
                +------------------+
                |  Movie Service   |
                | Business Rules   |
                +------------------+
                         │
                         ▼
                  +--------------+
                  |   MongoDB    |
                  +--------------+
```

---

# Tecnologias utilizadas

- Go
- Gin
- gRPC
- Protocol Buffers
- MongoDB
- Docker
- Docker Compose
- Swagger
- Zap Logger
- Testify
- UUID
- Mongo Driver

---

# Estrutura do projeto

```
movies-challenge/

│
├── api-gateway/
│   ├── cmd/
│   ├── docs/
│   ├── internal/
│   │   ├── adapters/
│   │   ├── config/
│   │   ├── logger/
│   │   ├── middleware/
│   │   └── routes/
│   ├── proto/
│   └── Dockerfile
│
├── movie-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── adapters/
│   │   ├── bootstrap/
│   │   ├── config/
│   │   ├── database/
│   │   ├── domain/
│   │   ├── mocks/
│   │   ├── ports/
│   │   ├── usecases/
│   │   └── logger/
│   ├── proto/
│   └── Dockerfile
│
├── proto/
├── docker-compose.yml
├── movies.json
└── README.md
```

---

# Arquitetura utilizada

O projeto foi desenvolvido utilizando Arquitetura Hexagonal.

A separação foi feita da seguinte forma:

### Domain

Contém as entidades de negócio.

Responsável apenas pelas regras de domínio.

Não possui dependência de infraestrutura.

---

### Ports

Define as interfaces utilizadas pelo domínio.

Exemplo:

- MovieRepository

---

### Use Cases

Implementa os casos de uso da aplicação.

Exemplos:

- List Movies
- Get Movie
- Create Movie
- Delete Movie

---

### Adapters

Implementações das portas.

Exemplos:

- MongoDB Repository
- gRPC Server
- HTTP Handlers

---

### Infrastructure

Responsável pela configuração da aplicação.

Exemplos:

- Logger
- Mongo
- Config
- Bootstrap

---

# Funcionalidades

- Listagem de filmes
- Busca por ID
- Cadastro de filme
- Exclusão de filme
- Paginação
- Seed automático do banco
- Índices no MongoDB
- Logs estruturados
- Request ID
- Graceful Shutdown

---

# Como executar

## Requisitos

- Docker Desktop
- Docker Compose

---

## Executar aplicação

Na raiz do projeto:

```bash
docker compose up --build
```

---

## Swagger

Após subir a aplicação:

```
http://localhost:8080/docs/index.html
```

---

## Health Check

```
GET /health
```

---

# Endpoints

## Listar filmes

```
GET /movies?page=1&limit=20
```

---

## Buscar filme

```
GET /movies/{id}
```

---

## Criar filme

```
POST /movies
```

Body

```json
{
    "id":999999,
    "title":"Meu Filme",
    "year":"2026"
}
```

---

## Remover filme

```
DELETE /movies/{id}
```

---

# Comunicação entre serviços

A comunicação entre os microsserviços é realizada utilizando gRPC.

O contrato é definido através do arquivo:

```
proto/movie.proto
```

A geração dos arquivos é realizada através do Protocol Buffers.

---

## Event Driven

As operações de criação e remoção de filmes são processadas de forma assíncrona utilizando RabbitMQ.

Fluxo:

POST /movies ou DELETE /movies/{id}
↓
API Gateway
↓
Movie Service
↓
Publicação de evento no RabbitMQ
↓
Worker consumidor
↓
MongoDB

# Banco de dados

MongoDB

Coleção:

```
movies
```

Durante a inicialização da aplicação:

- cria índices
- verifica existência dos dados
- realiza seed automaticamente caso necessário

---

# Testes

Executar testes do Movie Service

```bash
cd movie-service

go test ./...
```

Executar testes do API Gateway

```bash
cd api-gateway

go test ./...
```

Foram implementados:

- testes unitários utilizando mocks
- testes unitários do domínio sem mocks

---

# Logs

Foi implementado logging estruturado utilizando Zap Logger.

Cada requisição recebe um Request ID.

Exemplo:

```
request_id
method
path
status
latency
client_ip
```

---

# Docker

O projeto é composto por três containers:

- API Gateway
- Movie Service
- MongoDB

Todos inicializados através do Docker Compose.

Também foram configurados:

- Health Check
- Rede dedicada
- Volumes persistentes

---

# Decisões arquiteturais

## API Gateway

Foi utilizado para desacoplar a API REST do microsserviço responsável pelas regras de negócio.

---

## gRPC

Escolhido pela alta performance, contrato fortemente tipado e comunicação eficiente entre microsserviços.

---

## MongoDB

Adequado para o formato do dataset disponibilizado pelo desafio e pela facilidade de evolução do modelo.

---

## Arquitetura Hexagonal

Permite baixo acoplamento entre domínio e infraestrutura, facilitando testes, manutenção e substituição de implementações.

---

## Repository Pattern

Toda persistência foi abstraída através de interfaces, permitindo desacoplamento entre regra de negócio e banco de dados.

---