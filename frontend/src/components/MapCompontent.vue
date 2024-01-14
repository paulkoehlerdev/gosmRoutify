<script setup lang="ts">
import { inject, onMounted, type Ref, ref } from 'vue'
import "leaflet/dist/leaflet.css";
import markerIconUrl from "@/assets/images/marker-icon.png";
import markerIconRetinaUrl from "@/assets/images/marker-icon-2x.png";
import markerShadowUrl from "@/assets/images/marker-shadow.png";
import L, { Map, LatLng, FeatureGroup } from 'leaflet'
import type { Emitter, EventType } from 'mitt'

const map: Ref<Map | undefined> = ref(undefined);
const focusPointGroup: Ref<FeatureGroup | undefined> = ref(undefined);

const mapRef: Ref<HTMLElement | undefined> = ref(undefined);

const eventBus = inject('eventBus') as Emitter<Record<EventType, any>>;

L.Icon.Default.prototype.options.iconUrl = markerIconUrl;
L.Icon.Default.prototype.options.iconRetinaUrl = markerIconRetinaUrl;
L.Icon.Default.prototype.options.shadowUrl = markerShadowUrl;
L.Icon.Default.imagePath = "";

onMounted(() => {
  if (mapRef.value === undefined) {
    return;
  }

  map.value = L.map(mapRef.value, {
    zoomControl: false,
  }).setView([48.137154, 11.576124], 10);

  L.control.zoom({
    position: 'bottomright'
  }).addTo(map.value);

  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map.value);

  focusPointGroup.value = L.featureGroup().addTo(map.value);

})

eventBus.on('focusPoint', (p: LatLng) => {
  if (focusPointGroup.value === undefined) {
    return
  }

  focusPointGroup.value.clearLayers();
  L.marker(p).addTo(focusPointGroup.value);

  map.value?.flyTo(p, 15);
});

eventBus.on('startRoute', ({ addresses }) => {
  if (focusPointGroup.value === undefined) {
    return
  }

  focusPointGroup.value?.clearLayers();

  addresses.forEach(({ point }: { point: LatLng }) => {
    if (focusPointGroup.value === undefined) {
      return
    }

    L.marker(point).addTo(focusPointGroup.value);
  })

  map.value?.flyToBounds(focusPointGroup.value.getBounds(), {
    padding: [100, 100],
  })
});

eventBus.on('foundRoute', ({ route }) => {
  route.forEach(({ geojson }: { geojson: any }) => {
    if (focusPointGroup.value === undefined) {
      return
    }

    L.geoJSON(geojson).addTo(focusPointGroup.value)
  });

  if (focusPointGroup.value === undefined) {
    return
  }

  map.value?.flyToBounds(focusPointGroup.value.getBounds(), {
    padding: [100, 100],
  })
})

</script>

<template>
<div class="map" ref="mapRef"></div>
</template>

<style scoped>

.map {
  flex-grow: 1;
  z-index: 1;
}

</style>