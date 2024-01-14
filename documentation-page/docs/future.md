---
title: Zukünftige Entwicklung
slug: future
sidebar_position: 11
---

# Zukünftige Entwicklung

## Einbauen der Unterscheidung zwischen KFZ-, Fußgänger- und Fahrradrouten in das Frontend

Das Routing Backend hat aktuell einen Hardcoded-Modus für KFZ. Dieser soll in Zukunft in das Frontend eingebaut werden.

## Erweiterung des Routing-Backends für Isochronen

Das Routing-Backend soll in Zukunft um Isochronen erweitert werden. Dabei soll es möglich sein, Isochronen für einen
Startpunkt und eine maximale Zeit zu berechnen. Die Isochronen sollen dann als GeoJSON-Geometrien zurückgegeben werden.

## Verbesserung der Routen

Die Routen sind noch unvollständig. Verschiedene OSM-Tags werden noch nicht beachtet. So werden z. B. Ampeln nicht
berücksichtigt, sowie andere Vorfahrtsregelungen. Abbiegevorgänge werde ausschließlich durch eine "Strafe" nach
Abbiegewinkel bestraft. Dies führt dazu, dass die Routen nicht immer optimal sind und manchmal viele Ampeln einem leicht
längeren Weg auf einer Schnellstraße bevorzugt werden.

## Verbesserung der Performance

Die Performance des Backends ist aktuell noch nicht optimal. Die Berechnung der Routen dauert noch zu lange. Dies liegt
vor allem daran, das die Gewichtungen der Kanten im Graph und der Graph Allgemein on the fly generiert werden. Dies
könnte durch eine Vorberechnung der Gewichtungen und des Graphen verbessert werden. Dabei müsste wahrscheinlich nicht
der ganze Graph vorberechnet werden, sondern vor allem bevorzugte Routen (z. B. Autobahnen).

## Verbesserung der Datenqualität

Die Suche nach Orten ist noch nicht optimal. Hier werden aktuell ausschließlich die Tags der Nodes und Ways verwendet.
Allerdings gibt es auch Komplexere Addressen, die ausschließlich in Relationen abgebildet sind. Außerdem werden die
Ortsnamen nicht ergänzt, wie in der OSM-Wiki
beschrieben<sup>[[Quelle]](https://wiki.openstreetmap.org/wiki/DE:Key:addr:*)</sup>.

## Verbesserung des Daten-Imports

Der Daten-Import ist aktuell noch sehr langsam. Hier kann eventuell mit parallelen Imports gearbeitet werden. Außerdem
könnte Datenbank-Sharding verwendet werden, um die Datenbank (für den Import) zu beschleunigen.