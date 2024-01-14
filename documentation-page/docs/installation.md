---
title: Installation
slug: installation
sidebar_position: 2
---

# Installation

## Installation des Backends

Aktuell kann das Backen nur lokal aus dem Quellcode gebaut werden. In Zukunft wird es eventuell Docker-Images geben.

## Voraussetzungen

Das Projekt benötigt folgende Software:

- [`Go 1.21.6` https://go.dev](https://go.dev/)
- [`gcc` https://gcc.gnu.org/](https://gcc.gnu.org/) (Wird für CGO benötigt)
- [`make` https://www.gnu.org/software/make/](https://www.gnu.org/software/make/)
- [`git` https://git-scm.com](https://git-scm.com/)

## Installation

1. Klonen Sie das Repository mit `git clone`
```bash
git clone https://github.com/paulkoehlerdev/gosmRoutify.git
```
oder (wenn sie bereits einen ssh-key für GitHub hinterlegt haben)
```bash
git clone git@github.com:paulkoehlerdev/gosmRoutify.git
```

2. Wechseln Sie in das Verzeichnis des Projektes
```bash
cd gosmRoutify
```

3. Bauen Sie das Projekt mit `make build`
```bash
make build
```

4. Downloaden Sie OSM-Rohdaten (`.osm.pbf`) von [Geofabrik](https://www.geofabrik.de/). Die Daten werden im Beispiel in `./resources/data/germany-latest.osm.pbf` gespeichert.
```bash
wget https://download.geofabrik.de/europe/germany-latest.osm.pbf -O ./resources/data/germany-latest.osm.pbf
```

5. Importieren Sie die Daten in die Datenbank mit der loader Binary, die `make build` erstellt hat.

:::info
Dieser Schritt kann etwas dauern. Der import von Deutschland dauert ca. 30 Minuten.
Für die Entwicklung empfiehlt es sich daher einen kleineren Datensatz zu verwenden. (z. B. Oberbayern, wobei der Import nurnoch ca. 2 Minuten dauert)
:::

```bash
./bin/loader -import ./resources/data/germany-latest.osm.pbf -database ./resources/germany.db
```

6. Kopieren Sie die Beispiel-Konfiguration in die Konfigurationsdatei. Hier müssen Sie die Datenbank-URL anpassen, wenn Sie einen anderen Datensatz verwenden.
```bash
cp ./resources/config.example.json ./resources/config.json
```

7. Starten Sie den Server mit der server Binary, die `make build` erstellt hat.
```bash
./bin/router -config ./resources/config.json
```

8. Der Server ist nun unter `http://localhost:3000` erreichbar. Sie können nun die API verwenden. Der Port und der Bind-Host können in der Konfigurationsdatei angepasst werden.

## Installation des Frontends

Das Frontend ist in TypeScript geschrieben und verwendet `Vue.JS 3`.

## Voraussetzungen

Das Projekt benötigt folgende Software:

- [`Node.JS 21.2.0` https://nodejs.org](https://nodejs.org/)
- [`npm 10.2.3` https://www.npmjs.com](https://www.npmjs.com/)
- [`git` https://git-scm.com](https://git-scm.com/)

## Installation

1. Klonen Sie das Repository mit `git clone`
:::info
Wenn sie bereits das Backend geklont haben, können Sie diesen Schritt überspringen.
:::
```bash
git clone https://github.com/paulkoehlerdev/gosmRoutify.git
```
oder (wenn sie bereits einen ssh-key für GitHub hinterlegt haben)
```bash
git clone git@github.com:paulkoehlerdev/gosmRoutify.git
```

2. Wechseln Sie in das Verzeichnis des Projektes
```bash
cd gosmRoutify/frontend
```

3. Installieren Sie die Abhängigkeiten mit `npm install`
```bash
npm install
```

4. Starten Sie den Entwicklungsserver mit `npm run dev`
```bash
npm run dev
```

5. Das Frontend ist nun unter `http://localhost:5173` erreichbar. Wenn sie den Port oder host des Backends geändert haben, müssen Sie die Datei `.env.development` anpassen. Die Konfiguration `VUE_API_URL` muss auf die URL des Backends zeigen.