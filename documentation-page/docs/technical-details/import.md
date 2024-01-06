---
title: Daten-Import
---

## OSM-Daten

[Geofabrik](https://www.geofabrik.de/) stellt tagesaktuelle OpenStreetMap daten in verschiedenen Formaten zur Verfügung: 

- `.osm.pbf`: OSM Protobuf File<sup>[[mehr]](https://wiki.openstreetmap.org/wiki/PBF_Format)</sup>
- `.osm.xml`/`.osm.bz2`: Klassisches OSM-XML (Wird durch `.osm.pbf` ersetzt)<sup>[[mehr]](https://wiki.openstreetmap.org/wiki/OSM_XML)</sup>
- `.shp.zip`: ESRI-Shapefile (von [Geofabrik](https://www.geofabrik.de/) generiert und nur für kleine Ausschnitte erhältlich)

Da [`protobuf`](https://github.com/protocolbuffers/protobuf) exzellent von Golang unterstützt wird und hierfür von OpenStreetMap fertige Schemata herausgegeben werden (siehe [hier](https://github.com/openstreetmap/OSM-binary) für mehr dazu) fällt meine Entscheidung darauf, dieses Format zu verwenden. `protobuf` hat außerdem den nützlichen Vorteil, dass die Daten seriell geladen werden können, was überhaupt erst das Einhalten der unter [Zielsetzung](/docs/goals) gesetzten Ziele zu erreichen.

## Die Pipeline

Die Pipeline besteht aus 3 einfachen Schritten:

- Lesen des `protobuf`
- Transformieren der Daten in den Graphen
- Schreiben des Graphen in das Dateisystem

Um die Daten seriell lesen zu können und mit möglichst geringem zusätzlichen Speicheraufwand in den Graphen Transformieren zu können wird die ursprüngliche OSM-Datei zweimal durchlaufen.

Dabei werden im ersten Durchgang für jeden relevanten Weg die dazugehörigen OSM-IDs gespeichert. Hierbei ergeben sich bereits alle knoten und kanten des Graphen, jedoch fehlt die Georeferenz dessen noch völlig.  

Dafür werden die Dateien ein zweites mal durchlaufen, wobei für alle relevanten Nodes die Koordinaten und Tags gespeichert werden.