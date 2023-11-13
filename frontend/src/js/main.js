

var map = L.map('map').setView([48.137154, 11.576124], 10);

L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map);

let start = L.marker();
let end = L.marker();
let line = L.geoJSON();

function onMapClick(e) {
    if (!map.hasLayer(start)) {
        start = L.marker(e.latlng).addTo(map);
        return;
    }

    if(!map.hasLayer(end)) {
        end = L.marker(e.latlng).addTo(map);
        doRouteRequest(start.getLatLng(), end.getLatLng());
        return;
    }

    map.removeLayer(start);
    map.removeLayer(end);

    if(map.hasLayer(line)) {
        map.removeLayer(line);
    }
}

async function doRouteRequest(start, end) {
    const url = `/api/route?start=[${start.lng},${start.lat}]&end=[${end.lng},${end.lat}]`
    const response = await fetch(url);
    const json = await response.json()

    line = L.geoJSON(json, {
        style: {
            "color": "#0048ff",
            "weight": 5,
            "opacity": 0.65
        }
    }).addTo(map);
}

map.on('click', onMapClick);