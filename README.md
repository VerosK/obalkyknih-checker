# Obalkyknih.cz microservice


Webová služba obalkyknih.cz nepodporuje automatický failover (viz (API DOC)[obalky] a je potřeba failover zajistit aplikačně.

Tahle mikroslužba periodicky testuje cache1 a cache2 a vrací jako
backend tu služby, která zrovna běží.

## License:

2-BSD

[obalky]: https://www.obalkyknih.cz/doc/Dokumentace_API_OKCZ_3.3.pdf
