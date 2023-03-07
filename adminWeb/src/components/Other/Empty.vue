<template>
  <img :src="image" alt="nothing" />
</template>
<script>
import { reactive, toRefs, defineComponent } from 'vue'
import Config from '@/helper/config'

export default defineComponent({
  name: 'EmptyEd',
  setup() {
    let states = reactive({
      image: Config.images[3]
    })
    return {
      ...toRefs(states)
    }
  },
})
</script>