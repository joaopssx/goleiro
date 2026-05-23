# service-cpf

o `service-cpf` é o microsserviço responsável pela validação e formatação de cpfs no validator-hub.

## regras de validação

- remoção de qualquer caractere não numérico (ex: traços e pontos)
- verificação do tamanho exato (11 dígitos)
- bloqueio de sequências onde todos os números são repetidos (ex: 111.111.111-11)
- validação matemática utilizando o cálculo dos dois dígitos verificadores estabelecidos pela receita federal

## endpoints

### `GET /health`

retorna o status de saúde do serviço.

**resposta (200 OK):**
```json
{
  "status": "ok"
}
```

### `GET /validate?cpf={cpf}`

valida um cpf via query string.

**exemplo de requisição:**
`GET /validate?cpf=123.456.789-09`

### `GET /cpf/{cpf}`

valida um cpf via path param.

**exemplo de requisição:**
`GET /cpf/12345678909`

### `POST /validate`

valida um cpf via payload json.

**exemplo de requisição:**
```json
{
  "cpf": "123.456.789-09"
}
```

### formato de resposta da validação

em caso de sucesso (200 OK):
```json
{
  "valido": true,
  "formatado": "123.456.789-09",
  "mensagem": "cpf válido"
}
```

em caso de falha na validação (200 OK):
```json
{
  "valido": false,
  "formatado": "111.111.111-11",
  "mensagem": "cpf inválido: sequência de dígitos repetidos"
}
```

em caso de requisição malformada (400 Bad Request):
```json
{
  "valido": false,
  "formatado": "",
  "mensagem": "cpf não informado"
}
```

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
