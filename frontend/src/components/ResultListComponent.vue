<script setup lang="ts">
import type { Address } from '@/api/entities/address';
import ResultComponent from '@/components/ResultComponent.vue'
import { inject, ref, type Ref } from 'vue'
import type { Emitter, EventType } from 'mitt'

const list: Ref<Address[]> = ref([]);
const currentIndex: Ref<number | undefined> = ref(undefined);

const eventBus = inject('eventBus') as Emitter<Record<EventType, any>>;

eventBus.on('searchQueryResults', ({ results, index }) => {
  list.value = results;
  currentIndex.value = index;
});

eventBus.on('selectAddress', () => {
  list.value = [];
});

</script>

<template>

<ResultComponent v-for="item in list" :address="item" :index="currentIndex"/>

</template>

<style scoped>

</style>