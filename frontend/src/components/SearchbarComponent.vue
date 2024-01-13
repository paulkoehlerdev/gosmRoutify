<script setup lang="ts">
import 'bootstrap'
import { computed, inject, type Ref, ref } from 'vue'
import { debounce } from 'lodash'
import { fetchAddresses } from '@/api/api'
import type { Emitter, EventType } from 'mitt'
import type { Address } from '@/api/entities/address'

export type Type = 'start' | 'point' | 'end';

const props = defineProps<{
  type: Type
  index: number
  address: Address | undefined
}>()

const eventBus = inject('eventBus') as Emitter<Record<EventType, any>>

const icon = computed(() => {
  switch (props.type) {
    case 'start':
      return 'bi-crosshair'
    case 'end':
      return 'bi-geo-alt'

    case 'point':
    default:
      return 'bi-three-dots-vertical'
  }
})

const placeholder = computed(() => {
  switch (props.type) {
    case 'start':
      return 'Start'
    case 'end':
      return 'Ziel'

    case 'point':
    default:
      return 'Zwischenziel'
  }
})


const value: Ref<string | undefined> = ref(undefined)

function emitChange() {
  eventBus.emit('searchQuery', {
    query: value.value,
    index: props.index
  })

  triggerSearchDebounced(value.value ?? '')
}

function triggerSearch(query: string) {

  fetchAddresses(query).then((results) => {
    eventBus.emit('searchQueryResults', {
      results,
      index: props.index
    })
  }).catch(console.log)
}

const triggerSearchDebounced = debounce(triggerSearch, 500)

function getValue(): string {
  if (props.address !== undefined) {
    if (props.address.Name !== "") {
      return props.address.Name;
    }

    return `${props.address.Street} ${props.address.Housenumber}, ${props.address.City}`;
  }

  return value.value ?? "";
}

function unselectResult() {
  eventBus.emit('selectAddress', {
    index: props.index,
    address: undefined,
  })
}

</script>

<template>
  <div class="input-group mb-2">
    <span class="input-group-text">
        <i :class="'bi ' + icon"></i>
    </span>
    <input type="text"
           autocomplete="false"
           class="form-control"
           :disabled="address !== undefined"
           :placeholder=placeholder
           v-model=value
           @focus="emitChange"
           @input="emitChange"
    />
    <button v-if="address !== undefined"
            class="btn btn-light"
            @click="unselectResult"
    >
      <i class="bi bi-x"></i>
    </button>
  </div>
</template>

<style scoped>

</style>