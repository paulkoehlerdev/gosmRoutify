<script setup lang="ts">
import "bootstrap";
import SearchbarComponent from '@/components/SearchbarComponent.vue'
import ResultListComponent from '@/components/ResultListComponent.vue'
import { inject, ref, type Ref } from 'vue'
import type { Address } from '@/api/entities/address'
import type { Emitter, EventType } from 'mitt'
import type { LatLng } from 'leaflet'
import { fetchLocateAddress, fetchRoute } from '@/api/api'
import RouteInformationComponent from '@/components/RouteInformationComponent.vue'

const eventBus = inject('eventBus') as Emitter<Record<EventType, any>>;

const selectedAddresses: Ref<{
  query: string,
  address?: Address
  point?: LatLng
}[]> = ref([
  { query: "", address: undefined, point: undefined },
  { query: "", address: undefined, point: undefined },
])

type searchType = "start" | "point" | "end";
function indexToSearchtype(index: number, length: number): searchType {
  if (index === 0) {
    return "start"
  }

  if (index === length - 1) {
    return "end"
  }

  return "point"
}

eventBus.on('searchQuery', ({ query, index }) => {
  selectedAddresses.value[index].query = query;
});

const route: Ref<{ geojson:any, time:number, distance:number }[] | undefined> = ref(undefined);
const loading: Ref<boolean> = ref(false);
const error: Ref<string | undefined> = ref(undefined);

eventBus.on('selectAddress', async ({ address, index }) => {
  selectedAddresses.value[index].address = address

  if (!selectedAddresses.value.map((v) => v.address === undefined).reduce((a, b) => a || b, false)) {
    const addresses = await Promise.all(selectedAddresses.value.map(async (v) => {
      const point = await fetchLocateAddress(import.meta.env.VITE_API_URL, v.address as Address)
      return { ...v, point }
    }));

    loading.value = true;
    eventBus.emit('startRoute', { addresses })
  } else {
    route.value = undefined;
    error.value = undefined;
  }
});

eventBus.on('startRoute', async ({ addresses }: {addresses: { point: LatLng }[]}) => {
  try {
    route.value = await fetchRoute(import.meta.env.VITE_API_URL, addresses.map((v) => v.point))

    eventBus.emit('foundRoute', { route: route.value });
    loading.value = false;
    error.value = undefined;

  } catch (e) {
    loading.value = false;
    route.value = undefined;
    error.value = "Beim Planen der Route ist ein Fehler aufgetreten";
  }
});

</script>

<template>
  <div class="m-4" id="sidebar">
    <div class="row">
      <div class="col">
        <SearchbarComponent v-for="({ address }, index) in selectedAddresses"
                            :type="indexToSearchtype(index, selectedAddresses.length)"
                            :index="index"
                            :address="address"/>
        <RouteInformationComponent :route="route" :loading="loading" :error="error"/>
        <ResultListComponent />
      </div>
    </div>
  </div>
</template>

<style scoped>
#sidebar {
  z-index: 2;
  position: absolute;
  width: 30%;
}
</style>