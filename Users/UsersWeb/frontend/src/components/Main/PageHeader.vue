<template>
  <a-page-header
    style="border: 1px solid rgb(235, 237, 240);background: #fefefe;"
    :title="props.data.title"
    :breadcrumb="{ routes }"
    :sub-title="props.data.subtitle"
  />
</template>
<script lang="ts" setup>

const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
})
const routes = props.data.routers
</script>
