---
title: API
slug: api
sidebar_position: 3
---

# API

`gosmRoutify` bietet zwei API-Endpunkte an. Einen für die Routenberechnung und einen für die Suche nach Orten.

## Routen-API

Die Routen-API ist unter `GET /api/route` erreichbar. \
Sie erwartet einen Parameter `r`, der die Koordinaten der Route als Liste von Koordinatenpaaren in der Form `lon,lat`als
Base64Uri-Encodetes JSON enthält. \
Der Parameter kann in Javascript mit folgendem Snippet erstellt werden:

```js
function encodeRoute(route) {
    return encodeURIComponent(btoa(JSON.stringify(route)));
}

encodeRoute([
    [11.555806872727274, 48.15499445454545],
    [11.568533958333333, 48.14278539166667],
]);
```

Die Antwort enthält für jeden Wegpunkt die Distanz und die Zeit, die benötigt wird, um von diesem Wegpunkt zum nächsten
zu gelangen. Außerdem enthält sie die GeoJSON-Geometrie der Route.

### Beispiel

```bash
curl -X GET "https://api.gosmroutify.xyz/api/route?r=W1sxMS41Njg1MzM5NTgzMzMzMzMsNDguMTQyNzg1MzkxNjY2NjddLFsxMS41NTU4MDY4NzI3MjcyNzQsNDguMTU0OTk0NDU0NTQ1NDVdXQ==" -H "accept: application/json"
```

```json
[
  {
    "distance": 2241.2408995677297,
    "time": 197,
    "geojson": {
      "type": "FeatureCollection",
      "features": [
        {
          "type": "Feature",
          "geometry": {
            "type": "LineString",
            "coordinates": [
              ...
            ]
          },
          "properties": null
        }
      ]
    }
  },
  ...
]
```

## Search-API

Die Search-API ist unter `GET /api/search` erreichbar. \
Sie erwartet einen Parameter `q`, der den Suchbegriff enthält.

Die Antwort enthält eine Liste von Orten, die den Suchbegriff enthalten. Die Orte sind dabei in der Form eines `Address`
-Objekts definiert. Das Address-Objekt ist wie folgt definiert:

```go
type Address struct {
    OsmID      int64  `json:"OsmID"`
    Housenumber string `json:"Housenumber"`
    Street     string `json:"Street"`
    City       string `json:"City"`
    Postcode   string `json:"Postcode"`
    Country    string `json:"Country"`
    Suburb     string `json:"Suburb"`
    State      string `json:"State"`
    Province   string `json:"Province"`
    Floor      string `json:"Floor"`
    Name       string `json:"Name"`
}
```

### Beispiel

```bash
curl -X GET "https://api.gosmroutify.xyz/api/search?q=Roter%20Würfel" -H  "accept: application/json"
```

```json
[
  {
    "OsmID": 97390347,
    "Housenumber": "64",
    "Street": "Lothstraße",
    "City": "München",
    "Postcode": "80335",
    "Country": "DE",
    "Suburb": "",
    "State": "",
    "Province": "",
    "Floor": "",
    "Name": "Roter Würfel"
  },
  ...
]
```

## Locate-API

Die Locate-API ist unter `GET /api/locate` erreichbar. \
Sie erwartet einen Parameter `id`, der eine gültige Osm-ID einer Addresse enthält.

Die Antwort enthält die Koordinaten der Addresse als Koordinatenpaar in der Form `lon,lat`.

### Beispiel

```bash
curl -X GET "https://api.gosmroutify.xyz/api/locate?id=97390347" -H  "accept: application/json"
```

```json
[
  11.555806872727274,
  48.15499445454545
]
```