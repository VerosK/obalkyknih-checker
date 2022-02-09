# Obalkyknih.cz microservice


Webová služba obalkyknih.cz nepodporuje automatický failover (viz [API DOC][obalky]) a je potřeba failover zajistit aplikačně.

Tahle mikroslužba periodicky testuje cache1 a cache2 a vrací jako
backend tu služby, která zrovna běží.

## Zkompilovat

    # zkompilovat
    go build

## Spustit

    ./obalkyknih 

## Zjistit  živý server

Microservice vrací JSON řetězec s odkazem na naposledy živý testovací server.

    curl http://localhost:8000

    "https://cache1.obalkyknih.cz/"


### Zjistit živý server z PHP

Program vytváří soubor, který obsahuje aktuálně běžící server. Standardně ho vytváří v `/tmp/obalkyknih.php` 
Soubor obsahuje šablonu `$OBALKYKNIH_BASEURL="https://cache1.obalkyknih.cz";` a mělo by být možné ho naincludovat do PHP.

Umístění souboru a jeho obsah se dá změnit z příkazové řádky.


### Parametry spuštění

    ./obalkyknih -checkInterval=90s
    ./obalkyknih -listenAddress=:8080
    ./obalkyknih -help

### Stav microservice

    curl http://localhost:8000/status

### Stav souboru

    cat /tmp/obalkyknih.php

## License:

2-BSD

[obalky]: https://www.obalkyknih.cz/doc/Dokumentace_API_OKCZ_3.3.pdf
