# Obalkyknih.cz microservice


Webová služba obalkyknih.cz nepodporuje automatický failover (viz [API DOC][obalky] a je potřeba failover zajistit aplikačně.

Tahle mikroslužba periodicky testuje cache1 a cache2 a vrací jako
backend tu služby, která zrovna běží.

## Jak zkompilovat

    # stáhnout knihovny
    export GOPATH=/tmp/gopath-$$
    mkdir $GOPATH
    go get github.com/prometheus/client_golang/prometheus/promhttp
    # zkompilovat
    go build

## License:

2-BSD

[obalky]: https://www.obalkyknih.cz/doc/Dokumentace_API_OKCZ_3.3.pdf
