# gateway do validator-hub

o gateway atua como ponto de entrada para todos os serviços de validação. ele é responsável por rotear chamadas, autenticar requisições e proteger os serviços subjacentes contra sobrecarga.

## fluxo de uma requisição

1. o cliente faz a chamada para o gateway
2. o gateway gera logs estruturados no formato json
3. o `ratelimit` verifica se o ip ultrapassou a janela de 60 requisições/minuto
4. o `auth` valida a presença do header `X-API-Key`
5. o `router` reescreve a url, adiciona o `X-Forwarded-For` via httputil e encaminha para o serviço
6. a resposta volta pelo mesmo caminho até o cliente

## variáveis de ambiente

- `API_KEY`: define a chave necessária no header `X-API-Key` (obrigatória na subida do container).

## endpoints

### roteamento
- `/cpf/*`: roteia chamadas para o `service-cpf`
- `/cnpj/*`: roteia chamadas para o `service-cnpj`
- `/cep/*`: roteia chamadas para o `service-cep`
- `/email/*`: roteia chamadas para o `service-email`

**exemplo de requisição:**
```shell
curl -H "X-API-Key: sua-chave-aqui" http://localhost:8080/cpf/validate?cpf=123.456.789-09
```

### health check consolidado
- `GET /health` (não exige autenticação)
- `GET /health/detalhado` (com versão e uptime)

o gateway valida em paralelo com um timeout de 2 segundos.
- retorna 200 ok se todos estiverem bem
- retorna 207 multi-status se algum estiver degradado
- retorna 503 service unavailable se algum estiver inacessível

**exemplo de resposta:**
```json
{
  "status": "ok",
  "servicos": {
    "cpf": {
      "status": "ok",
      "latencia_ms": 12
    },
    "cnpj": {
      "status": "ok",
      "latencia_ms": 8
    }
  },
  "verificado_em": "2026-05-23T14:00:00Z"
}
```

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
