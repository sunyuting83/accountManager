<template>
  <a-page-header
    style="border: 1px solid rgb(235, 237, 240);background: #fefefe;"
    :title="props.data.title"
    :breadcrumb="{ routes }"
    :sub-title="props.data.subtitle"
    @back="() => router.back()"
    v-if="props.data.back"
  />
  <a-page-header
    style="border: 1px solid rgb(235, 237, 240);background: #fefefe;"
    :title="props.data.title"
    :breadcrumb="{ routes }"
    :sub-title="props.data.subtitle"
    v-else
  />
</template>
<script lang="ts" setup>
import { useRouter } from 'vue-router'
const router = useRouter()

const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
})
const routes = props.data.routers
</script>
