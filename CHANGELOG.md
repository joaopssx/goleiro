# changelog

## [0.5.0] - 2026-05-23

### adicionado
- suporte completo a graceful shutdown nos 5 serviços utilizando `signal.NotifyContext`
- fechamento brando com sinal do os (sigterm/sigint) e timeout via `SHUTDOWN_TIMEOUT_SECONDS`
- logs estritos de parada: "encerrando serviço...", e timeout
- suporte a rolling update (zero downtime) via `update_config: start-first`
- documentação `docs/operacoes.md` detalhando gerenciamento dos containers

## [0.4.0] - 2026-05-23
