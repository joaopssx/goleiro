# changelog

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
