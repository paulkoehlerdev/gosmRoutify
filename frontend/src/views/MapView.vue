<script setup lang="ts">
import MapCompontent from '@/components/MapCompontent.vue'
import NavigationComponent from '@/components/NavigationComponent.vue'

</script>

<template>
  <main>
    <NavigationComponent />
    <MapCompontent />
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
