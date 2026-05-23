# guia de operações do validator-hub

este documento explica como operar e gerenciar os microsserviços do hub em ambiente docker local ou de produção.

## como reiniciar um serviço

se precisar forçar o restart de algum serviço específico (ex: gateway), utilize o comando do docker-compose:
```shell
docker-compose restart gateway
```

## como fazer redeploy com zero downtime

o projeto está configurado com a política de rolling update (ordem `start-first`). isso significa que o docker sobe o container novo e espera ele ficar saudável antes de derrubar o antigo, garantindo que o tráfego não caia.

para forçar um redeploy limpo de um serviço específico após alterar o código:
```shell
docker-compose up -d --no-deps --build nome-do-servico
```
*substitua `nome-do-servico` pelo nome correto, como `service-cpf`.*

## como verificar logs

para checar a saúde geral, acompanhe os logs no seu terminal:
```shell
docker-compose logs -f
```

se quiser os logs de um serviço só:
```shell
docker-compose logs -f gateway
```

## como forçar encerramento

todos os serviços possuem graceful shutdown, aguardando até 15 segundos para encerrar as conexões ativas de forma limpa.

- ao dar `docker-compose stop` ou `Ctrl+C`, o sistema enviará um SIGTERM.
- após 15 segundos (ou o tempo configurado na variável `SHUTDOWN_TIMEOUT_SECONDS`), se os serviços ainda não tiverem finalizado suas tarefas, eles farão um *exit 1* automático forçando a queda.
- o docker está configurado com um `stop_grace_period` de 20 segundos para não matar o processo via SIGKILL antes do nosso timeout interno atuar.

---
projeto desenvolvido por joaopssx - joão pedro de sena santana
