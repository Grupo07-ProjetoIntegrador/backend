# Projeto Integrador - Backend e Prototipo

Este projeto possui uma API em Go integrada ao Supabase e um frontend React/Vite em `Prototipo-main`.

Atualmente, a integracao principal permite cadastrar treinamentos no backend e visualizar a lista de treinamentos no prototipo.

## Pre-requisitos

- Go instalado
- Node.js e npm instalados
- Acesso ao banco Supabase do projeto
- Git instalado

## Configuracao do ambiente

Crie um arquivo `.env` na raiz do projeto com a URL de conexao do Supabase:

```env
DATABASE_URL=coloque_a_url_do_supabase_aqui
```

Use o arquivo `.env.example` como modelo.

## Instalar dependencias

Na raiz do projeto, instale as dependencias do Go:

```powershell
go mod download
```

Depois instale as dependencias do frontend:

```powershell
cd Prototipo-main
npm.cmd install
```

## Rodar o backend

Em um terminal, na raiz do projeto:

```powershell
go run .\cmd\api
```

O backend deve ficar disponivel em:

```text
http://localhost:8080
```

## Rodar o frontend

Em outro terminal:

```powershell
cd Prototipo-main
npm.cmd run dev -- --host 127.0.0.1
```

Acesse no navegador:

```text
http://127.0.0.1:5173/
```

## Rotas usadas pelo frontend

O prototipo consome as seguintes rotas do backend:

```text
POST http://localhost:8080/api/treinamentos/cadastrar
GET  http://localhost:8080/api/treinamentos
```

## Observacoes para subir no GitHub

Nao envie estes arquivos/pastas:

```text
.env
Prototipo-main/node_modules/
Prototipo-main/dist/
```

O `package-lock.json` deve ser enviado, pois ajuda outras pessoas a instalarem as mesmas versoes das dependencias.

## Status atual

- Cadastro de treinamentos conectado ao backend
- Listagem de treinamentos conectada ao backend
- Visualizacao dos treinamentos no prototipo
- Exclusao no prototipo ainda funciona apenas localmente na tela, pois o backend ainda nao possui rota de exclusao
- Edicao no prototipo ainda nao salva no banco, pois o backend ainda nao possui rota de edicao
