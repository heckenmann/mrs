# MRS - MultiRegexSuche (Testprojekt)
## Beschreibung
- Durchsucht Text-Dokumente nach vordefinierten Regex-Gruppen
- Mehrere Reguläre Ausdrücke werden zu einem Ausdruck zusammengefasst
- Restful Service über HTTP
- Docker-Image verfügbar
- Geschrieben in Golang

## Konfiguration
Die Konfiguration geschieht über die mrs.yml.

## docker-compose.yml
```
---
version: '3.5'

services:
  mrs:
    image: heckenmann/mrs:latest
    volumes:
      - ./mrs.yml:/opt/mrs/mrs.yml
    ports:
      - "8080:8080"
```

## Testfiles
- https://norvig.com/big.txt

## Request
```
curl -XPOST --data @big.txt http://localhost:8080/analysis
```