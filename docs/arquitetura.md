# arquitetura do validator-hub

o validator-hub adota uma arquitetura baseada em microsserviços. a ideia central é possuir um serviço para cada tipo de validação, isolando a complexidade e permitindo escalabilidade independente.

## componentes principais

### gateway
ponto de entrada de todas as requisições externas. suas responsabilidades incluem:
- receber as chamadas HTTP
- aplicar rate limiting
- verificar a integridade dos serviços (health check)
- rotear a requisição para o serviço de validação adequado

### microsserviços de validação
cada microsserviço (como `service-cpf`, `service-cnpj`, etc) é focado em processar uma única regra de negócio. eles recebem os dados do gateway, aplicam o algoritmo de validação e retornam o resultado.

### infraestrutura
os serviços são orquestrados através de containers docker. a comunicação externa chega em uma porta unificada no gateway, enquanto a comunicação interna ocorre na rede do docker.

## fluxo de requisição

1. o cliente faz uma requisição HTTP para a URL pública (gateway)
2. o gateway processa o rate limiting e valida o path da rota
3. se aprovado, o gateway encaminha via proxy reverso ou chamada HTTP interna para o serviço alvo (ex: `service-cpf`)
4. o serviço valida os dados e devolve um JSON com o resultado
5. o gateway repassa a resposta HTTP para o cliente

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
