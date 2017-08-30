# Obalkyknih.cz microservice


Webová služba obalkyknih.cz nepodporuje automatický failover (viz [API DOC][obalky] a je potřeba failover zajistit aplikačně.

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


### Parametry spuštění

    ./obalkyknih -checkInterval=90s
    ./obalkyknih -listenAddress=:8080

### Stav microservice

    curl http://localhost:8000/status
    curl http://localhost:8000/metrics

## License:

2-BSD

[obalky]: https://www.obalkyknih.cz/doc/Dokumentace_API_OKCZ_3.3.pdf
