# changelog

## [0.7.0] - 2026-05-23

### adicionado
- framework proprietário de logs estruturados em json (`encoding/json`) substituindo `log` nativo em todos os componentes
- campos unificados de observabilidade (`duracao_ms`, `arquivo`, `estado_anterior`) com controle por severidade (`LOG_LEVEL`)
- middleware nativo interceptador no gateway e serviços injetando request tracing nos logs
- mecanismo de métricas em memória baseado em `sync/atomic` imune a data race (tempo_medio, hits, erros)
- endpoint `GET /metricas` implementado em cada serviço unitariamente
- endpoint `GET /metricas` consolidado no gateway (aggregating endpoints filhos com concorrência)
- elaborada a documentação `docs/observabilidade.md` com as instruções completas das payloads de telemetria

## [0.6.0] - 2026-05-23

### adicionado
- rate limit customizado via algoritmo token bucket (utilizando `sync.map`) no gateway
- isenção automática de tráfego para ips internos (localhost e redes de infraestrutura) no middleware de limit
- bloqueio elegante de requests retornando `429` com os parâmetros `tente_novamente_em` e o cabeçalho `retry-after`
- integração de circuit breaker por serviço no gateway (fechado, aberto, semi-aberto) 100% stdlib
- acoplamento do estado dos disjuntores ao reverse proxy através de transição de `http.RoundTripper` (`ErrorHandler` injetando 503 customizado)
- logging estruturado para a máquina de estado ("de: fechado", "para: aberto")
- testes unitários puros para transições do breaker e vazamento de tokens do rate limiter
- documentação arquitetônica da nova camada no `docs/resiliencia.md`
