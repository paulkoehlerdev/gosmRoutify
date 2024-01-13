<script setup lang="ts">
import "bootstrap";
import { computed, type Ref, ref } from 'vue'

export type Type = 'start' | 'point' | 'end';

const props = defineProps<{
  type?: Type
}>()

const emit = defineEmits<{
  (e: 'update', value: string | undefined): void
}>();

const icon = computed(() => {
  switch (props.type) {
    case 'start': return 'bi-crosshair'
    case 'end': return 'bi-geo-alt'

    case 'point':
    default: return 'bi-three-dots-vertical'
  }
});

const placeholder = computed(() => {
  switch (props.type) {
    case 'start': return 'Start'
    case 'end': return 'Ziel'

    case 'point':
    default: return 'Zwischenziel'
  }
});

const value: Ref<string | undefined> = ref(undefined);

function emitChange(event: Event) {
  value.value = event.target?.value;
  emit('update', value.value)
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
           :placeholder=placeholder
           :value=value
           @focus="emitChange"
           @input="emitChange"
    />
  </div>
</template>

<style scoped>

</style>