<script setup lang="ts">
import type { Address } from '@/api/entities/address'
import { fetchLocateAddress } from '@/api/api'
import type { Point } from '@/api/entities/point'

defineProps<{
  address?: Address,
}>()

const emit = defineEmits<{
  (e: 'triggerFocus', value: { address: Address, p: Point }): void
}>()

function centerPoint(address: Address | undefined) {
  if (address === undefined) {
    return
  }

  fetchLocateAddress(address).then((p: Point) => {
    emit('triggerFocus', { address, p })
  })
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
          <button class="btn btn-primary" @mouseenter="centerPoint(address)">
            <i class="bi bi-geo-alt-fill"></i>
          </button>
        </div>
      </div>
    </div>
  </div>

</template>

<style scoped>

</style>