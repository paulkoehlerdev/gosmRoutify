<script setup lang="ts">
import { onMounted, type Ref, ref } from 'vue'
import "leaflet/dist/leaflet.css";
import L, { Point, Map, FeatureGroup } from 'leaflet'

const map: Ref<Map | undefined> = ref(undefined);
const focusPointGroup: Ref<FeatureGroup | undefined> = ref(undefined);

onMounted(() => {
  map.value = L.map('map', {
    zoomControl: false,
  }).setView([48.137154, 11.576124], 10);

  L.control.zoom({
    position: 'bottomright'
  }).addTo(map.value);

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map.value);

  focusPointGroup.value = L.layerGroup().addTo(map.value);
})


function focusPoint(p: Point) {
  console.log(p);
  focusPointGroup.value?.clearLayers();
  L.marker(p).addTo(focusPointGroup.value);
  map.value?.flyTo(p, 10);
}

</script>

<template>
<div id="map"></div>
</template>

<style scoped>

#map {
  flex-grow: 1;
  z-index: 1;
}

</style>