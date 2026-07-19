# Go Expert — Client-Server API (Cotação de Câmbio)

Desafio do curso Go Expert (FullCycle) que implementa dois sistemas em Go que trocam informações sobre cotação de câmbio USD-BRL respeitando limites estritos de timeout.

## Estrutura do Projeto

```
.
├── client.go      # Cliente HTTP
├── server.go      # Servidor HTTP + banco SQLite
├── go.mod
├── go.sum
└── README.md
```

## Requisitos

- Go 1.26+ (ou versão compatível)

## Como Rodar

### 1. Iniciar o Servidor

Em um terminal, execute:

```bash
go run server.go
```

O servidor iniciará na porta **8080** e expõe o endpoint `GET /cotacao`.

**Funcionamento do servidor:**
- Ao receber uma requisição em `/cotacao`, consome a API externa `https://economia.awesomeapi.com.br/json/last/USD-BRL` (timeout: **200ms**).
- Persiste o valor do `bid` em um banco SQLite (`cotacoes.db`) (timeout: **10ms**).
- Retorna um JSON com o campo `bid`.

### 2. Executar o Cliente

Em outro terminal, execute:

```bash
go run client.go
```

**Funcionamento do cliente:**
- Faz uma requisição HTTP para `http://localhost:8080/cotacao` (timeout: **300ms**).
- Extrai o campo `bid` do JSON de resposta.
- Salva a cotação no arquivo **cotacao.txt** no formato:
  ```
  Dólar: 5.8473
  ```

## Timeouts

| Componente | Operação            | Timeout |
|------------|---------------------|---------|
| Servidor   | Chamada API externa | 200ms   |
| Servidor   | Persistência SQLite | 10ms    |
| Cliente    | Requisição HTTP     | 300ms   |

Se qualquer timeout for excedido, o erro é registrado no console.

## Dependências

- **modernc.org/sqlite** — driver SQLite puro em Go (sem CGO).
