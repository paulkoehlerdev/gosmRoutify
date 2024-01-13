<script setup lang="ts">
import MapCompontent from '@/components/MapCompontent.vue'
import NavigationComponent from '@/components/NavigationComponent.vue'
import { ref } from 'vue'

const mapRef = ref()

</script>

<template>
  <main>
    <NavigationComponent :mapRef="mapRef"/>
    <MapCompontent :ref="mapRef"/>
  </main>
</template>

<style scoped>

main {
  display: flex;
  flex-direction: row;
  width: 100vw;
  height: 100vh;
}

</style>
