<script setup lang="ts">
defineProps<{
  route?: { geojson: any, time: number, distance: number }[] | undefined,
  loading: boolean,
  error?: string,
}>()

function formatDuration(seconds: number): string {
  const time = {
    d: Math.floor(seconds / 86400),
    h: Math.floor(seconds / 3600) % 24,
    m: Math.floor(seconds / 60) % 60,
    s: Math.floor(seconds) % 60
  }
  return Object.entries(time)
    .filter(val => val[1] !== 0)
    .map(([key, val]) => `${val}${key} `)
    .join(' ')
}

function formatDistance(distance: number): string {
  if (distance > 1500) {
    return `${Math.round(distance / 1000)}km`
  }

  return `${Math.round(distance)}m`
}

</script>

<template>
  <div v-if="route !== undefined || loading || error !== undefined"
       class="card mt-1">
    <h5 class="card-header">Routeninformation</h5>

    <div v-if="loading" class="card-body">
      <div class="d-flex justify-content-center align-items-center">
        <div class="spinner-grow" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
    </div>

    <ul v-if="route !== undefined" class="list-group list-group-flush">
      <li v-for="part in route"
          class="list-group-item">
        <i class="bi bi-alarm"></i> {{ formatDuration(part.time) }}<br>
        <i class="bi bi-signpost"></i> {{ formatDistance(part.distance) }}
      </li>
    </ul>

    <div v-if="error !== undefined" class="card-body">
      {{ error }}
    </div>
  </div>
</template>

<style scoped>

</style>