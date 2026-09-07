# Backend Core (Go) — Shopping Flamboyant

API RESTful de alta performance desenvolvida em **Go (Golang)**, responsável pelo núcleo de regras de negócio, persistência de dados no **PostgreSQL (Supabase)**, cálculo de validação presencial por **Geofencing (Haversine)**, autenticação **Google OAuth 2.0** e orquestração de disparos e relatórios.

Repositório oficial: **[Grupo07-ProjetoIntegrador/backend](https://github.com/Grupo07-ProjetoIntegrador/backend)**

---

## 📑 Sumário

1. [Visão Geral & Responsabilidades](#-visão-geral--responsabilidades)
2. [Estrutura de Pastas e Arquitetura](#-estrutura-de-pastas-e-arquitetura)
3. [Script SQL de Criação do Banco (Supabase)](#-script-sql-de-criação-do-banco-supabase)
4. [Configuração do Google Cloud OAuth 2.0](#-configuração-do-google-cloud-oauth-20)
5. [Variáveis de Ambiente (`.env`)](#-variáveis-de-ambiente-env)
6. [Como Rodar Localmente](#-como-rodar-localmente)
7. [Referência Completa de Endpoints da API](#-referência-completa-de-endpoints-da-api)
8. [Regras de Negócio Especiais](#-regras-de-negócio-especiais)

---

## 🎯 Visão Geral & Responsabilidades

O Backend Core em Go atua como o cérebro central da plataforma:
- **Gestão de Lojas e Treinamentos**: CRUD completo de capacitações, locais, categorias e lojas do shopping.
- **Validação de Presença com Geofencing**: Cálculo de distância real via fórmula de Haversine para garantir que o participante estava fisicamente no auditório durante o check-in por QR Code.
- **Fluxo OAuth 2.0 do Google**: Permite que o organizador conecte sua conta Google Workspace para criação de formulários e envio de e-mails em seu nome.
- **Integração com Microsserviço Python**: Comunicação assíncrona para geração de Google Forms, e-mails pelo Gmail e emissão de PDFs.
- **Métricas e Dashboards**: Agregação de dados para os painéis analíticos e explorador de lojas.

---

## 📂 Estrutura de Pastas e Arquitetura

O projeto segue a convenção canônica de layout da comunidade Go:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Ponto de entrada, configuração de CORS e inicialização do servidor HTTP
├── internal/
│   ├── database/
│   │   └── supabase.go             # Conexão singleton com PostgreSQL via driver github.com/lib/pq
│   ├── handlers/
│   │   ├── rotas.go                # Mapeamento central de todas as rotas e métodos HTTP
│   │   ├── loja_handler.go         # Cadastro e listagem de lojas
│   │   ├── treinamento_handler.go  # Criação, edição, exclusão e geração de formulários
│   │   ├── presenca_handler.go     # Check-in por geolocalização e gestão de participantes
│   │   ├── oauth_handler.go        # Ciclo de vida OAuth 2.0 do Google (start, callback, refresh, status)
│   │   ├── explorador_handler.go   # Métricas consolidadas de lojas e engajamento
│   │   ├── historico_loja_handler.go # Histórico individual de capacitações por loja
│   │   ├── dashboard_handler.go    # Estatísticas gerais do dashboard
│   │   ├── inscricao_handler.go    # Recepção de webhooks de inscrição do Google Forms
│   │   ├── relatorio_handler.go    # Proxy para geração e download de PDFs de dossiê e chamada
│   │   └── upload_handler.go       # Importação de planilhas de presença
│   ├── models/
│   │   ├── loja.go                 # Estruturas de dados de lojas e explorador
│   │   ├── treinamento.go          # Estruturas de treinamentos e filtros
│   │   ├── presenca.go             # Estruturas de presenças e confirmações
│   │   ├── dashboard.go            # Modelos de estatísticas e gráficos
│   │   └── inscricaoFormsRequest.go# Payload do webhook do Google Forms
│   └── repositories/
│       ├── loja_repo.go            # Queries SQL para a tabela 'lojas'
│       ├── treinamento_repo.go     # Queries SQL para a tabela 'treinamentos'
│       ├── presenca_repo.go        # Queries SQL para a tabela 'presencas'
│       ├── explorador_repo.go      # Queries analíticas agregadas
│       ├── dashboard_repo.go       # Consultas consolidadas de KPIs
│       ├── historico_loja_repo.go  # Consultas de histórico por LUC
│       └── job_repo.go             # Enfileiramento de tarefas na tabela 'job_queue'
├── .env.example                    # Modelo oficial de variáveis de ambiente
├── Dockerfile                      # Build multistage para container de produção
├── go.mod                          # Módulo Go e dependências declaradas
└── go.sum                          # Checksums de integridade das dependências
```

---

## 🗄️ Script SQL de Criação do Banco (Supabase)

Para inicializar o banco de dados em um novo projeto Supabase, acesse o **SQL Editor** do Supabase Dashboard e execute:

```sql
-- Ativa extensão para UUIDs automáticos
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tabela de Lojas do Shopping
CREATE TABLE IF NOT EXISTS lojas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    luc VARCHAR(50) NOT NULL UNIQUE,
    nome VARCHAR(255) NOT NULL,
    segmento VARCHAR(100) DEFAULT 'Não Informado',
    status BOOLEAN DEFAULT true,
    email VARCHAR(255),
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Tabela de Locais com Geofencing
CREATE TABLE IF NOT EXISTS locais_treinamento (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome_local VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    raio_amplitude INTEGER NOT NULL DEFAULT 100, -- Raio de tolerância em metros
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Tabela de Treinamentos
CREATE TABLE IF NOT EXISTS treinamentos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tema VARCHAR(255) NOT NULL,
    descricao TEXT,
    categoria VARCHAR(100),
    data DATE NOT NULL,
    horario_inicio TIMESTAMP WITH TIME ZONE NOT NULL,
    horario_fim TIMESTAMP WITH TIME ZONE NOT NULL,
    local TEXT,
    modalidade VARCHAR(50) DEFAULT 'PRESENCIAL',
    conteudo TEXT,
    capacidade_maxima INTEGER DEFAULT 50,
    segmento_alvo VARCHAR(100) DEFAULT 'Geral',
    status VARCHAR(50) DEFAULT 'AGENDADO',
    objetivo TEXT,
    observacoes TEXT,
    material_apoio TEXT,
    responsavel VARCHAR(255),
    area_responsavel VARCHAR(255),
    tags TEXT,
    recorrente BOOLEAN DEFAULT false,
    local_id UUID REFERENCES locais_treinamento(id) ON DELETE SET NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Tabela de Presenças / Inscrições
CREATE TABLE IF NOT EXISTS presencas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    treinamento_id UUID NOT NULL REFERENCES treinamentos(id) ON DELETE CASCADE,
    loja_id UUID NOT NULL REFERENCES lojas(id) ON DELETE RESTRICT,
    nome_participante VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefone VARCHAR(50),
    cargo VARCHAR(100),
    status_presenca VARCHAR(50) DEFAULT 'PENDENTE', -- 'PENDENTE', 'CONFIRMADO', 'AUSENTE'
    data_confirmacao TIMESTAMP WITH TIME ZONE,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Vínculo entre Treinamentos e Formulários do Google Forms
CREATE TABLE IF NOT EXISTS formularios_treinamento (
    treinamento_id UUID PRIMARY KEY REFERENCES treinamentos(id) ON DELETE CASCADE,
    google_form_id VARCHAR(255) NOT NULL,
    url_formulario TEXT NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 6. Tokens OAuth 2.0 do Google por Usuário
CREATE TABLE IF NOT EXISTS google_oauth_tokens (
    user_id VARCHAR(255) PRIMARY KEY,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_type VARCHAR(50) DEFAULT 'Bearer',
    scope TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 7. Fila de Tarefas Assíncronas (Worker Python)
CREATE TABLE IF NOT EXISTS job_queue (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 🔑 Configuração do Google Cloud OAuth 2.0

Para permitir que o Backend gerencie o fluxo de consentimento do organizador:
1. Acesse o **[Google Cloud Console](https://console.cloud.google.com)**.
2. Em **APIs e Serviços > Credenciais**, clique em **Criar Credenciais > ID do cliente OAuth**.
3. Selecione o tipo **Aplicativo da Web (Web Application)**.
4. Defina o nome como `Backend Flamboyant`.
5. Em **URIs de redirecionamento autorizados**, adicione:
   - Desenvolvimento: `http://localhost:8080/api/oauth/google/callback`
   - Produção: `https://jpmallflamboyant.live/api/oauth/google/callback`
6. Copie o **Client ID** e o **Client Secret** e preencha no arquivo `.env`.

---

## ⚙️ Variáveis de Ambiente (`.env`)

Crie o arquivo `backend/.env` baseado no modelo [backend/.env.example](.env.example):

```env
# Conexão direta ou Transaction Pooler com o banco PostgreSQL no Supabase
DATABASE_URL=postgresql://postgres:sua_senha@db.sua_referencia.supabase.co:5432/postgres

# Credenciais do OAuth 2.0 (Google Cloud Console > Web Application)
GOOGLE_CLIENT_ID=seu-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=seu_client_secret

# Redirecionamento após consentimento no Google
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8080/api/oauth/google/callback

# URL do frontend para onde o usuário retorna após o callback
FRONTEND_BASE_URL=http://localhost:5173

# URL base do serviço de automações Python
AUTOMACOES_PUBLIC_URL=http://localhost:8000

# Credenciais opcionais para armazenamento de mídia no Supabase Storage
SUPABASE_URL=https://sua_referencia.supabase.co
SUPABASE_KEY=sua_anon_ou_service_role_key
```

---

## 🚀 Como Rodar Localmente

### Pré-requisitos
- **Go 1.22+** instalado (`go version`)
- Acesso à internet para conexão com o banco Supabase

### Passos de Execução
```powershell
# 1. Entre no diretório do backend
cd backend

# 2. Crie o arquivo .env com suas credenciais reais
cp .env.example .env

# 3. Baixe e sincronize os módulos Go
go mod tidy

# 4. Inicie o servidor
go run .\cmd\api\main.go
```

O terminal exibirá:
```text
Backend rodando e .env processado!
Conexão com o Supabase estabelecida com sucesso!
Servidor rodando na porta: http://localhost:8080
```

---

## 📚 Referência Completa de Endpoints da API

### 1. Treinamentos
| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `GET` | `/api/treinamentos` | Lista todos os treinamentos cadastrados. |
| `POST` | `/api/treinamentos/cadastrar` | Cadastra um novo treinamento e inicia geração de formulário. |
| `PUT` | `/api/treinamentos/editar` | Atualiza os dados de um treinamento existente. |
| `DELETE` | `/api/treinamentos/deletar?id={id}` | Remove um treinamento e seus registros vinculados. |
| `GET` | `/api/treinamentos/dashboard` | Retorna KPIs e totais consolidados para o dashboard. |

### 2. Google Forms & Convites
| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `POST` | `/api/treinamentos/gerar-formulario` | Solicita ao Python a criação manual do Google Forms. |
| `GET` | `/api/treinamentos/formulario?id={id}` | Retorna a URL pública do Google Forms gerado. |
| `POST` | `/api/treinamentos/apagar-formulario` | Exclui o vínculo e apaga o arquivo no Google Drive. |
| `POST` | `/api/treinamentos/regerar-formulario` | Remove o form anterior e cria um novo do zero. |
| `POST` | `/api/treinamentos/disparar-convite` | Dispara e-mails com o link do Forms segmentado por setor. |
| `POST` | `/api/treinamentos/webhook-forms` | Recebe a inscrição repassada pelo Python e insere como `PENDENTE`. |

### 3. Presença, Check-in & Geofencing
| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `POST` | `/api/presencas/confirmar` | Valida coordenadas GPS via Haversine e confirma presença do lojista. |
| `GET` | `/api/treinamentos/presencas?treinamento_id={id}` | Lista participantes de um treinamento e respectivos status. |
| `POST` | `/api/treinamentos/presencas/manual` | Adiciona um participante manualmente pelo painel admin. |
| `PUT` | `/api/treinamentos/presencas/editar` | Altera status manual (`PENDENTE`, `CONFIRMADO`, `AUSENTE`). |
| `DELETE` | `/api/treinamentos/presencas/deletar` | Exclui um registro de presença. |
| `GET` | `/api/treinamentos/geofencing?id={id}` | Retorna as coordenadas e raio de tolerância do local do treinamento. |
| `POST` | `/api/locais/cadastrar` | Cadastra um novo local (ex: Auditório Cristal, latitude, longitude, raio). |
| `GET` | `/api/locais` | Lista todos os locais cadastrados. |

### 4. Lojas & Explorador
| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `POST` | `/api/lojas/cadastrar` | Cadastra uma nova loja no shopping (LUC, Nome, Segmento). |
| `GET` | `/api/lojas` | Lista todas as lojas com opção de filtro por segmento. |
| `GET` | `/api/lojas/explorador?data_inicio=YYYY-MM-DD&data_fim=YYYY-MM-DD` | Retorna métricas de taxa de participação e total de presenças por loja. |
| `GET` | `/api/lojas/historico?luc={LUC}` | Retorna todo o histórico de capacitações frequentadas por aquela loja. |

### 5. Google OAuth 2.0
| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `POST` | `/api/oauth/google/start` | Inicia o fluxo de autorização gerando a URL de consentimento do Google. |
| `GET` | `/api/oauth/google/callback` | Callback chamado pelo Google com o `code`, salva o token no banco. |
| `GET` | `/api/oauth/google/status?user_id={id}` | Informa se o usuário possui token Google válido e não expirado. |
| `POST` | `/api/oauth/google/disconnect` | Remove as credenciais do usuário do banco de dados. |

### 6. Relatórios & Exportações
| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `GET` | `/api/relatorios/loja/dossie?luc={LUC}` | Solicita ao Python o PDF do dossiê histórico da loja. |
| `GET` | `/api/relatorios/treinamento/chamada?treinamento_id={id}` | Solicita ao Python o PDF da ata de chamada do evento. |
| `POST` | `/api/treinamentos/upload` | Importa lista de presenças em formato de planilha. |

---

## 📐 Regras de Negócio Especiais

### 1. Fórmula de Haversine (Validação Presencial)
Para evitar que lojistas confirmem presença remotamente de suas lojas ou residências, o endpoint `POST /api/presencas/confirmar` recebe a latitude e longitude capturadas pelo navegador do smartphone do participante e compara com o ponto geográfico cadastrado no local do evento:

$$\Delta\sigma = 2 \arcsin\left(\sqrt{\sin^2\left(\frac{\Delta\phi}{2}\right) + \cos(\phi_1)\cos(\phi_2)\sin^2\left(\frac{\Delta\lambda}{2}\right)}\right)$$
$$d = R \cdot \Delta\sigma \quad (\text{com } R = 6.371.000\text{ metros})$$

Se a distância calculada $d$ for menor ou igual ao `raio_amplitude` configurado (por exemplo, 100 metros), a presença é aprovada e marcada como `CONFIRMADO`. Caso contrário, a API rejeita com status `400 Bad Request` indicando a distância atual do participante.

### 2. Auto-Refresh de Tokens OAuth
O backend verifica a validade do token OAuth armazenado na tabela `google_oauth_tokens`. Se o token estiver expirado, ele utiliza o `refresh_token` permanente para obter um novo `access_token` transparente perante o Google.
