# adnettest

Утилита для тестирования сетевой связанности между АРМ и контроллерами домена

## Использование

в командной строке за пустить утилиту adnettest в выходных данных будут перечислены адреса серверов и статусы необходимых портов, состояние должно быть open

adnettest [-domain <example.com>] [-dns <ip addr server[:port]>] [host.example.com] [192.168.0.1] [192.168.0.2] [192.168.0.3] ...

при запуске утилиты без параметров будет выполнено попытка получить суффикс домена из имени компьютера, в случае успеха выполняется резолв SRV записей контроллеров домена для полученного суффикса домена, и проверка сетевой связанности выполнится для записей SRV.

Явно задать днс сервер для поиска записей можно указав параметр "-dns <ip addr server[:port]>", порт по умолчанию 53.

Днс нужен для резолва адресов, т.к. в тесте используются IP адреса.


## Тестируемые порты и протоколы
   |Порт сервера  |    Служба|
   |--------------|---------------|
   |53/TCP/UDP  |     DNS |
   |88/TCP/UDP    |    Kerberos |
   |123/UDP      |     W32Time |
   |135/TCP      |     Сопоставитель конечных точек RPC|
   |137/UDP       |    netbios-ns |
   |138/UDP       |    netbios-dgm |
   |139/TCP       |    netbios-ssn |
   |389/TCP/UDP   |    LDAP |
   |445/TCP       |    SMB |
   |464/TCP/UDP    |   Изменение пароля в Kerberos |
   |636/TCP       |    LDAP SSL |
   |3268/TCP     |     LDAP GC |
   |3269/TCP     |     LDAP GC SSL |
