<script setup lang="ts">
import type { Address } from '@/api/entities/address';
import ResultComponent from '@/components/ResultComponent.vue'
import type { Point } from '@/api/entities/point'

defineProps<{
  list: Address[],
}>();

const emit = defineEmits<{
  (e: 'triggerFocus', value: { address: Address, p: Point }): void
}>()

</script>

<template>

<ResultComponent v-for="item in list" :address="item" @triggerFocus="(val) => {emit('triggerFocus', val)}"/>

</template>

<style scoped>

</style>