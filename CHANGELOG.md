# changelog

todas as mudanças notáveis deste projeto serão documentadas neste arquivo.

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
