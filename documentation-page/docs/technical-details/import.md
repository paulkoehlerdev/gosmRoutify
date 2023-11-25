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

