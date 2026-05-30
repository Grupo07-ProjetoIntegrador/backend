# Backend Go

Este diretório contém a API em Go que conversa com o Supabase e expõe as rotas usadas pelo frontend e pelas automações.

## O que fica aqui

* Código Go da API
* Arquivo [backend/.env.example](.env.example)
* Arquivo local [backend/.env](.env) com as chaves da máquina atual

## Sobre `node_modules`, `.venv` e `package.json`

* `node_modules` não pertence ao backend Go; ele é da parte frontend e deve existir só onde o `npm install` foi executado.
* `.venv` não pertence ao backend Go; ele é um ambiente virtual do Python e deve ficar dentro de `automacoes/` se você quiser isolar dependências Python.
* `package.json` do frontend também não pertence ao backend; ele deve ficar no diretório do frontend.

Se esses itens aparecerem na raiz do repositório, normalmente é porque algum comando foi executado na pasta errada. Para este projeto, o ideal é manter cada dependência dentro da sua subpasta.

## Pré-requisitos

* Go instalado
* Banco Supabase disponível
* Variáveis de ambiente configuradas

## Configuração do ambiente

1. Copie [backend/.env.example](.env.example) para [backend/.env](.env).
2. Preencha as chaves com os dados do seu ambiente.

### Variáveis usadas

* `DATABASE_URL`
* `GOOGLE_CLIENT_ID`
* `GOOGLE_CLIENT_SECRET`
* `GOOGLE_OAUTH_REDIRECT_URL`
* `FRONTEND_BASE_URL`

## Executar o backend

```powershell
go run .\cmd\api\main.go
```

O backend deve subir em:

```text
http://localhost:8080
```

## Rotas principais

* `POST /api/treinamentos/cadastrar`
* `GET /api/treinamentos`
* `POST /api/oauth/google/start`
* `GET /api/oauth/google/callback`

## Observações para GitHub

Não envie para o repositório:

* [backend/.env](.env)
* `node_modules/`
* `.venv/`

O arquivo [backend/.env.example](.env.example) deve ficar versionado para servir de modelo.
