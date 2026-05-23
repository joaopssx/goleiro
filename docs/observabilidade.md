# guia de observabilidade do validator-hub

este documento explica o padrão de logs estruturados e métricas de telemetria coletadas pelo projeto.

## logs estruturados em json

todo o ecossistema (os 4 serviços de validação e o gateway) loga exclusivamente no formato json de linha única. isso garante parsing imediato em ferramentas como datadog, elk ou grafana. não existem `printf` soltos.

### estrutura base
cada log possui no mínimo:
```json
{
  "timestamp": "2026-05-23T14:00:00Z",
  "nivel": "info",
  "servico": "gateway",
  "mensagem": "iniciando gateway na porta 8080"
}
```

### logs de requisição http
gerados automaticamente pelo middleware para cada request finalizado:
```json
{
  "timestamp": "2026-05-23T14:00:00Z",
  "nivel": "info",
  "servico": "service-cpf",
  "mensagem": "requisição processada",
  "metodo": "GET",
  "path": "/validate",
  "ip": "192.168.1.10",
  "status_code": 200,
  "duracao_ms": 14
}
```

### logs de mudança de estado (circuit breaker / shutdown)
```json
{
  "timestamp": "2026-05-23T14:05:00Z",
  "nivel": "info",
  "servico": "gateway",
  "mensagem": "circuit breaker alterado (cpf)",
  "estado_anterior": "fechado",
  "estado_novo": "aberto"
}
```

### logs de erro interno
```json
{
  "timestamp": "2026-05-23T14:10:00Z",
  "nivel": "erro",
  "servico": "service-cep",
  "mensagem": "erro ao iniciar servidor",
  "arquivo": "main.go",
  "linha": 59,
  "erro": "address already in use"
}
```

### como filtrar por nível
a severidade mínima exibida é controlada pela variável de ambiente `LOG_LEVEL`.
- `LOG_LEVEL=info` (padrão): exibe info, warn e erro.
- `LOG_LEVEL=warn`: exibe warn e erro.
- `LOG_LEVEL=erro`: exibe apenas erro.

---

## métricas em memória

cada serviço coleta métricas atomicamente na memória (sem onerar banco de dados). a contagem reinicia ao zerar o container (comportamento planejado para instâncias efêmeras).

### endpoints disponíveis

- **`GET /metricas` (nos serviços - ex: na porta 8081):** retorna json plano com contadores de requisições, válidas, inválidas, erros e tempo médio de resposta do serviço local.
- **`GET /metricas` (no gateway - porta 8080):** o gateway executa um scatter-gather, varrendo as métricas de todos os microsserviços via rede interna, somando tudo no objeto raiz, mas mantendo a granularidade no campo `servicos`.

exemplo de resposta do gateway:
```json
{
  "total_requisicoes": 1500,
  "total_validas": 800,
  "total_invalidas": 650,
  "total_erros": 50,
  "servicos": {
    "cpf": {
      "total_requisicoes": 1000,
      "total_validas": 800,
      "total_invalidas": 200,
      "total_erros": 0,
      "tempo_medio_ms": 12
    },
    "gateway": {
      "total_requisicoes": 500,
      "total_validas": 0,
      "total_invalidas": 0,
      "total_erros": 50,
      "tempo_medio_ms": 18
    }
  }
}
```

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
