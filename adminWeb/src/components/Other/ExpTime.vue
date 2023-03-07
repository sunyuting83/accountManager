<template>
  <span>{{FormatTime(DateTime)}}</span>
</template>
<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ExpTime',
  props: {
    DateTime: {
      type: Number
    }
  },
  setup() {
    const FormatTime = (d) => {
      const t = Date.parse(new Date()) / 1000
      const s = d - t
      if (s > 0) {
        if (s < 86400) {
          return Math.ceil(s / 3600) + "小时"
        }else if(s < 3600){
          return Math.ceil(s / 60) + "分钟"
        }else {
          return Math.ceil(s / 86400) + "天"
        }
      }else {
        return 0
      }
    }
    return {
      FormatTime
    }
  },
})
</script>