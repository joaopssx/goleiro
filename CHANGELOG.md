# changelog

todas as mudanças notáveis deste projeto serão documentadas neste arquivo.

## [0.4.0] - 2026-05-23

### adicionado
- health check profundo com cálculo de uptime e versão em cada serviço
- verificação paralela consolidada no gateway com timeouts (status 200, 207 ou 503)
- endpoint estendido `GET /health/detalhado` expondo dados métricos
- auto-restart no docker-compose (`unless-stopped`) e comandos nativos de healthcheck

## [0.3.0] - 2026-05-23

### adicionado
- implementação completa do `gateway` com autenticação via `X-API-Key`
- rate limiting nativo por ip usando sliding window (60 req/min)
- roteamento dinâmico via proxy reverso preservando headers e adicionando `X-Forwarded-For`
- log estruturado json para todas as requisições
- endpoint consolidador de `/health` monitorando todos os microsserviços

## [0.2.0] - 2026-05-23

### adicionado
- implementação completa do `service-cpf`
- endpoints de validação de cpf (get, post, path param)
- regras estritas da receita federal (cálculo de dígitos, sequências repetidas)
- documentação do serviço em `docs/service-cpf.md`

## [0.1.0] - 2026-05-23

### adicionado
- estrutura base do monorepo validator-hub
- gateway e serviços independentes (cpf, cnpj, cep, email)
- documentação inicial de arquitetura
- configurações de infraestrutura com docker-compose
