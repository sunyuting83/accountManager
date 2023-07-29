<template>
  <span v-if="props.DateTime !== 0">{{ dateTime }}</span>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { FormatDateTime } from '../../wailsjs/go/main/App'

const props = defineProps({
  DateTime: {
    type: Number,
    required: true,
    default: 0,
  },
})

const dateTime = ref('')

watch(
  () => props.DateTime,
  (newDateTime) => {
    if (typeof newDateTime === 'number') {
      FormatDateTime(newDateTime)
        .then((result: string) => {
          dateTime.value = result
        })
        .catch((err) => {
          dateTime.value = ''
        })
    }
  }
)
</script>
