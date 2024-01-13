<script setup lang="ts">
import "bootstrap";
import SearchbarComponent from '@/components/SearchbarComponent.vue'
import { fetchAddresses } from '@/api/api'
import ResultListComponent from '@/components/ResultListComponent.vue'
import { type Ref, ref } from 'vue'
import type { Address } from '@/api/entities/address'
import { debounce } from 'lodash'
import type { Point } from '@/api/entities/point'

const props = defineProps<{
  mapRef?: Ref
}>()

const resultList: Ref<Address[]> = ref([])

function onSearchChange(value: string) {
  fetchAddresses(value).then(val => resultList.value = val).catch(() => {})
}

function onFocusEvent(value: { address: Address, p: Point }) {
  console.log(value);
  props.mapRef?.value.focusPoint(value.p)
}

const onSearchChangeDebounced = debounce(onSearchChange, 500);

</script>

<template>
  <div class="m-4" id="sidebar">
    <div class="row">
      <div class="col">
        <SearchbarComponent type="start" @update="onSearchChangeDebounced"/>
        <SearchbarComponent type="end" @update="onSearchChangeDebounced"/>
        <ResultListComponent :list=resultList @triggerFocus="onFocusEvent"/>
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