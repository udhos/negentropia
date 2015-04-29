@rem win-run

pushd c:\tmp\devel

start cmd /k c:\redis\redis-server.exe

start cmd /k c:\tmp\devel\negentropia\win-run-world.cmd
start cmd /k c:\tmp\devel\negentropia\win-run-webserv.cmd
start cmd /k c:\tmp\devel\negentropia\win-run-webserv2.cmd

popd

@rem eof
