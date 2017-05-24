Negentropia - início rápido pt_BR
---------------------------------

- Passos 1-8: instalação (executados uma única vez)
- Passos 9-10: ciclo de desenvolvimento

1) Os servidores do negentropia só rodam em Windows.
   Rodar em Linux deve ser fácil, mas nunca tentei.
   Então, por enquanto, Windows.

2) Instalar Git no Windows - https://git-scm.com/downloads
   A instação do Git deve disponibilizar o utilitário "git-bash".

3) Instalar Redis no Windows - https://redis.io/download
   Instalado, o Redis deve disponibilizar algumas ferramentas na pasta "C:\redis" (este caminho é importante):
   ‪C:\redis\redis-cli.exe
   ‪C:\redis\redis-server.exe

4) Abra o git-bash e crie o diretório: /c/tmp/devel
   $ mkdir /c/tmp/devel
   (o prefixo /c no path acima é o drive C: do Windows)
   
5) Baixe o script win-gitbash-clone.sh para o diretório: /c/tmp/devel
   https://raw.githubusercontent.com/udhos/negentropia/master/win-gitbash-clone.sh

6) Dentro do git-bash, execute o script: win-gitbash-clone.sh
   Assim:
   $ /c/tmp/devel/win-gitbash-clone.sh

7) Edite o arquivo redis-zone-add.txt para atribuir uma zona inicial ao seu endereço de email.
   Esse endereço de email será usado para fazer login.

‪C:\tmp\devel\negentropia\redis-zone-add.txt

Por exemplo:

    hset everton.marques@gmail.com location z:simple_zone
    hset everton.marques@gmail.com password-sha1-hex 40bd001563085fc35165329ea1ff5c5ecbdbbeef

40bd001563085fc35165329ea1ff5c5ecbdbbeef é o SHA1 para a senha '123'.

8) Em um DOS prompt (cmd), execute:

    copy \tmp\devel\negentropia\config-common-sample.txt \tmp\devel\config-common.txt
    copy \tmp\devel\negentropia\config-webserv-sample.txt \tmp\devel\config-webserv.txt
    copy \tmp\devel\negentropia\config-world-sample.txt \tmp\devel\config-world.txt

9) Em um DOS prompt (cmd), execute:
	
    c:\tmp\devel\win-build
    c:\tmp\devel\win-run
    c:\tmp\devel\win-zone-add

10) Abra a URL do negentropia com um navegador:

http://localhost:8080/ne/

Fazer login usando o par email/senha cadastrado no passo 7.

---
FIM
---
