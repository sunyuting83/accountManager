<template>
  <span>{{FormatTime(DateTime)}}</span>
</template>
<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'FormaTime',
  props: {
    DateTime: {
      type: Number
    }
  },
  setup() {
    const FormatTime = (timestamp) => {
      const date = new Date(timestamp);
      const Y = date.getFullYear()
      const M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1)
      const D = date.getDate()
      const H = date.getHours() < 10 ? `0${date.getHours()}` : date.getHours()
      const minute = date.getMinutes() < 10 ? `0${date.getMinutes()}` : date.getMinutes()
      // const S = date.getSeconds() < 10 ? `0${date.getSeconds()}` : date.getSeconds()
      return `${Y}-${M}-${D} ${H}:${minute}`
    }
    return {
      FormatTime
    }
  },
})
</script>