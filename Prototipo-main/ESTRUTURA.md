# 📁 Estrutura de Pastas Necessária

Para o servidor Go funcionar corretamente, organize os arquivos assim:

```
seu-projeto/
├── main.go                           # Servidor Go principal
├── go.mod                            # (criar com: go mod init nome-do-projeto)
│
├── src/
│   └── imports/
│       ├── lista.html                # ✅ Já existe
│       ├── calendario.html           # ✅ Já existe
│       ├── formulariotreinamento.html # ✅ Já existe
│       ├── detalhesdotreinamento.html # ✅ Já existe
│       ├── detalhesdaloja.html       # ✅ Já existe
│       └── engajamento.html          # ✅ Já existe
│
├── static/
│   └── css/
│       └── style.css                 # CSS do Tailwind (se necessário)
│
└── assets/
    └── logo_2024_(1)-nvcRYmem.png    # Logo e outras imagens
```

---

## 🚀 Como Rodar

### 1. Inicializar módulo Go (primeira vez)
```bash
go mod init painel-jpmall
```

### 2. Rodar o servidor
```bash
go run main.go
```

### 3. Acessar no navegador
```
http://localhost:8080
```

---

## 📋 Rotas Disponíveis

| Rota                        | Arquivo HTML                     | Descrição                          |
|-----------------------------|----------------------------------|------------------------------------|
| `/`                         | → `/treinamentos/lista`          | Redireciona para lista             |
| `/treinamentos/lista`       | `lista.html`                     | Lista de treinamentos              |
| `/treinamentos/calendario`  | `calendario.html`                | Calendário de treinamentos         |
| `/treinamentos/novo`        | `formulariotreinamento.html`     | Criar novo treinamento             |
| `/treinamentos/editar`      | `formulariotreinamento.html`     | Editar treinamento                 |
| `/treinamentos/detalhes`    | `detalhesdotreinamento.html`     | Detalhes do treinamento            |
| `/engajamento`              | `engajamento.html`               | Dashboard de engajamento           |
| `/lojas/detalhes`           | `detalhesdaloja.html`            | Detalhes da loja                   |

---

## ⚠️ Limitações do HTML Estático

Funcionalidades que **NÃO funcionam** (apenas visual):

- ❌ Filtros e busca
- ❌ Formulários (criar/editar)
- ❌ Paginação
- ❌ Modais de confirmação
- ❌ Gráficos interativos
- ❌ Abas dinâmicas

Funcionalidades que **funcionam**:

- ✅ Navegação entre páginas
- ✅ Layout responsivo
- ✅ Visual/Design

---

## 💡 Próximos Passos

Para ter funcionalidades completas, considere:

1. **Abordagem 1 (Recomendada):** Usar `npm run build` e servir a pasta `dist/` com React funcionando
2. **Abordagem 2 (Atual):** Adicionar backend Go que processa formulários e dados reais
3. **Abordagem 3:** Usar templates Go (html/template) com dados dinâmicos
