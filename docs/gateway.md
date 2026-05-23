# gateway

gateway de entrada do validator-hub responsável pelo roteamento.

## endpoints públicos

- `GET /health`: status do gateway
- `POST /v1/cpf/validate`: roteia para o `service-cpf`
- `POST /v1/cnpj/validate`: roteia para o `service-cnpj`
- `POST /v1/cep/validate`: roteia para o `service-cep`
- `POST /v1/email/validate`: roteia para o `service-email`

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
