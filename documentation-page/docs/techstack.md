---
title: Technologien
slug: techstack
sidebar_position: 10
---

# Technologien

Für dieses Projekt wurden die folgende Software verwendet. Diese Liste ist möglicherweise nicht vollständig, deckt
jedoch den größten Teil ab.

## Backend (Routen-engine)

Im Backend wurde fast ausschließlich Golang und dessen Standartbibliotheken verwendet. Lediglich ein SQLite-Driver
musste als Bibliothek importiert werden. Die Daten werden in SQLite gespeichert.

- [`Go` https://go.dev](https://go.dev/)
- [`go-sqlite3` https://github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [`SQLite` https://www.sqlite.org](https://www.sqlite.org/index.html)

Das Demo-Backend ist auf einem VPS gehostet.

## Frontend

Das Frontend wurde in Vue.JS 3 geschrieben. Zusätzlich wurden einige Bibliotheken verwendet:

- [`Vue.js` https://vuejs.org](https://vuejs.org/)
- [`Leaflet` https://leafletjs.com](https://leafletjs.com/)
- [`Bootstrap` https://getbootstrap.com](https://getbootstrap.com/)
- [`Bootstrap Icons` https://icons.getbootstrap.com](https://icons.getbootstrap.com/)
- [`lodash` https://lodash.com](https://lodash.com/)
- [`mitt` https://github.com/developit/mitt](https://github.com/developit/mitt)

Das Demo-Frontend ist bei [Vercel](https://vercel.com) gehostet.

## Dokumentation

Diese Dokumentation wurde mit Docusaurus erstellt. Einem Open-Source-Projekt, dass von Meta entwickelt wurde.

- [`Docusaurus` https://docusaurus.io](https://docusaurus.io/)

Die Dokumentation ist auf [GitHub Pages](https://pages.github.com/) gehostet.