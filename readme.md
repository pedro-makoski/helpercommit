# Documentação
Primeiro adicione o exe, numa pasta separada na pasta do seu usuário. E depois adicione no seu homepath, a pasta.  
## Adicione repositórios 
```sh
helpercommit new <identificador-do-repositorio> <caminho-absoluto-do-repositorio>
```
## Comite e faça o push dos seus repositórios
```sh
helpercommit commit -desc=<descricao-opcional> -branch=<branch-do-projeto-opcional> -files=<arquivos-do-projeto-se-nao-tiver-vai-ser-todos> <identificador-do-repositorio> <nome-do-commit>
```
## Faça o pull dos seus projetos
Atualize seu repositório local.
```sh
helpercommit pull <identificador-do-repositorio>
```
## Importante
O identificador do repositório é você que cria ao definir no new. 
