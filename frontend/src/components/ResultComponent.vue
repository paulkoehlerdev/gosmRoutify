<script setup lang="ts">
import type { Address } from '@/api/entities/address'
import { fetchLocateAddress } from '@/api/api'
import { inject } from 'vue'
import type { Emitter, EventType } from 'mitt'
import type { LatLng } from 'leaflet'

const props = defineProps<{
  address?: Address,
  index?: number,
}>()

const eventBus = inject('eventBus') as Emitter<Record<EventType, any>>;

function centerPoint() {
  if (props.address === undefined) {
    return;
  }

  fetchLocateAddress(import.meta.env.VITE_API_URL, props.address).then((p: LatLng) => {
    eventBus.emit('focusPoint', p)
  }).catch(console.log);
}

function selectPoint() {
  eventBus.emit('selectAddress', {
    address: props.address,
    index: props.index,
  });
}

</script>

<template>

  <div class="card mt-1">
    <div class="card-body">
      <div class="row">
        <div class="col d-flex flex-column justify-content-center align-items-start">
          <b>{{ address?.Name }}</b><br v-if='address?.Name !== ""'>
          {{ address?.Street }} {{ address?.Housenumber }}<br
          v-if='address?.Street !== "" && address?.Housenumber !== ""'>
          {{ address?.Postcode }} {{ address?.City }}
        </div>

        <div class="col-auto d-flex align-items-center">
          <button class="btn btn-primary" @mouseenter="centerPoint" @click="selectPoint">
            <i class="bi bi-geo-alt-fill"></i>
          </button>
        </div>
      </div>
    </div>
  </div>

</template>

<style scoped>

</style>