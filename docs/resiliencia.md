# resiliência e proteção no validator-hub

o gateway atua como escudo principal para todos os microsserviços. ele implementa dois grandes padrões de resiliência sem depender de nenhuma biblioteca externa ou cache distribuído.

## rate limiting (token bucket)

evita abusos limitando a quantidade de requisições que um ip pode fazer em uma janela de tempo.

- **algoritmo:** token bucket por ip.
- **limite padrão:** 60 rpm (requisições por minuto). configurável via variável de ambiente `RATE_LIMIT_RPM`.
- **isenções:** ips locais (`127.0.0.1`, `::1`) e redes privadas (`10.x.x.x`, `192.168.x.x`) não sofrem limitação de tráfego.
- **limpeza (garbage collection):** o mapa de tokens é varrido a cada 5 minutos para limpar ips que não acessam mais a api, evitando vazamento de memória.

### comportamento ao exceder o limite

o gateway interrompe a requisição e devolve status http `429 too many requests`. ele também injeta o cabeçalho `Retry-After` informando em quantos segundos o cliente poderá tentar de novo.

exemplo de resposta json:
```json
{
  "erro": "limite de requisições excedido",
  "limite": 60,
  "janela": "1 minuto",
  "tente_novamente_em": "2026-05-23T14:01:00Z"
}
```

## circuit breaker

previne falhas em cascata isolando serviços defeituosos rapidamente.

- **rastreamento por downstream:** cada serviço roteado pelo gateway possui o seu próprio circuito (ex: o circuito do `cpf` é independente do circuito do `cep`).
- **abertura do circuito:** se ocorrerem **5 falhas de rede ou erros 5xx consecutivos num intervalo de 30 segundos**, o circuito muda para o estado `aberto`.
- **timeout do bloqueio:** o estado `aberto` dura exatos **20 segundos**. nesse meio-tempo, o proxy falha as requisições de imediato sem nem tentar bater na rede (*fast-fail*).
- **estado semi-aberto:** após os 20 segundos de resfriamento, ele transiciona para `semi-aberto` permitindo passar **1 requisição de teste**.
  - se falhar: reabre o circuito por mais 20 segundos.
  - se sucesso: fecha o circuito restaurando o tráfego total e zera a contagem de falhas.

### comportamento de bloqueio

quando o circuito está aberto (ou cai na falha de rede), o gateway responde imediatamente com o status http `503 service unavailable`:

exemplo de resposta json:
```json
{
  "erro": "serviço temporariamente indisponível",
  "servico": "cpf"
}
```

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
