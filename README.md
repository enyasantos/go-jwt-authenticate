# authentication-jwt

Projeto de autenticação JWT em Go, utilizando MongoDB como banco de dados e o framework Gin para API REST.

## Funcionalidades

- Cadastro de usuários com validação de dados
- Login com geração de access token e refresh token (JWT)
- Refresh de tokens
- Logout
- Middleware de autenticação para rotas protegidas
- Armazenamento seguro de senhas (bcrypt)
- Armazenamento e controle de refresh tokens no MongoDB

## Tecnologias

- Go (Golang)
- Gin (framework web)
- MongoDB (persistência)
- JWT (autenticação)
- Docker (opcional, para banco de dados)

## Estrutura

```
cmd/api/main.go              # Ponto de entrada da aplicação
internal/
  auth/                      # Lógica de autenticação e geração de tokens
  database/                  # Conexão com o MongoDB
  middlewares/               # Middlewares do Gin
  models/                    # Modelos de dados
  repositories/              # Repositórios de acesso ao banco
  server/                    # Handlers e configuração das rotas
```

## Como rodar

1. Configure as variáveis de ambiente no arquivo `.env` (exemplo já incluso).
2. Suba o MongoDB com Docker (opcional):
   ```sh
   docker-compose up -d
   ```
3. Instale as dependências:
   ```sh
   go mod tidy
   ```
4. Configure as variaveis de ambiente
    ```sh
    PORT=8080
    DATABASE_URI=mongodb://root:root@localhost:27017/authentication-jwt?authSource=admin
    DATABASE_NAME=authentication-jwt
    JWT_SECRET=mysecretkey
    JWT_SECRET_REFRESH=mysecretkeyrefresh
    ```
5 Rode a aplicação:
   ```sh
   go run ./cmd/api/main.go
   ```

## Rotas principais

- `POST /api/auth/register` — Cadastro de usuário
- `POST /api/auth/logon` — Login
- `POST /api/auth/refresh` — Refresh do token
- `POST /api/auth/logout` — Logout
- `GET /api/user` — Dados do usuário autenticado (rota protegida)

---

> Projeto para estudo de autenticação JWT com Go e MongoDB.
