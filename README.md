# validator-hub

o validator-hub é um projeto focado na centralização e padronização da validação de dados no brasil. ele recebe requisições em um gateway único e as roteia para microsserviços especializados, como validação de CPF, CNPJ, CEP, entre outros.

## como rodar localmente

para iniciar o projeto na sua máquina, certifique-se de ter o docker e o docker-compose instalados.

1. clone o repositório
2. acesse a raiz do projeto
3. execute o comando para subir os serviços:

```shell
docker-compose -f infra/docker-compose.yml up -d
```

## estrutura de pastas

- `/gateway`: serviço central de roteamento, rate limiting e health checks.
- `/service-cpf`: microsserviço responsável pela validação de CPF.
- `/service-cnpj`: microsserviço responsável pela validação de CNPJ.
- `/service-cep`: microsserviço responsável pela validação de CEP.
- `/service-email`: microsserviço responsável pela validação de email.
- `/infra`: configurações de infraestrutura (docker, nginx).
- `/docs`: documentações técnicas detalhadas sobre a arquitetura e serviços.

## endpoints disponíveis

o gateway centraliza as chamadas. a documentação específica de contratos pode ser encontrada na pasta `docs/`.

- `GET /health`: verifica a saúde do gateway e dos serviços
- `POST /v1/cpf/validate`: envia dados para validação no `service-cpf`
- `POST /v1/cnpj/validate`: envia dados para validação no `service-cnpj`
- `POST /v1/cep/validate`: envia dados para validação no `service-cep`
- `POST /v1/email/validate`: envia dados para validação no `service-email`

## como contribuir

para contribuir com o projeto:

1. crie uma branch com a sua feature ou correção
2. mantenha o padrão de código em letras minúsculas (exceto siglas técnicas)
3. atualize o `CHANGELOG.md` na raiz do projeto
4. abra um pull request detalhando suas alterações

## autor

projeto desenvolvido por joaopssx - joão pedro de sena santana
