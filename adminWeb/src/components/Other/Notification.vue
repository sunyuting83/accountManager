<template>
  <div :class="'notification is-light error '+showData.color" :style="{'top':hoverTop+'px'}" v-if="this.showData.active">
    <button class="delete" @click="closErr"></button>
    <p>{{this.showData.message}}</p>
  </div>
</template>

<script>
import { reactive, toRefs, watch, defineComponent } from 'vue'
export default defineComponent ({
  name: 'NotIfication',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      },
      message: {
        type: String,
        default: ""
      },
      color:{
        type: String,
        default: "is-info"
      }
    }
  },
  emits: ["update:active","update:message","update:color"],
  setup(props, context){
    let _data = reactive({
      hoverTop: 0
    })
    watch(() => props.showData.active,(newValue) => {
      context.emit("update:active", newValue)
      openError()
    })
    watch(() => props.showData.message,(newValue) => {
      context.emit("update:message", newValue)
    })
    watch(() => props.showData.color,(newValue) => {
      context.emit("update:color", newValue)
    })
    watch(() => _data.hoverTop,(newValue) => {
      _data.hoverTop = newValue
    })
    const openError = () => {
      const scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
      // context.emit("update:hoverTop", scrollTop)
      _data.hoverTop = scrollTop
      const _this = props
      setTimeout(function(){
        _this.showData.active = false
      },1500)
    }
    const closErr = () => {
      const _this = props
      _this.showData.message = ""
      _this.showData.active = false
    }
    return {
      ...toRefs(_data),
      closErr
    }
  }
})
</script>

<style scoped>
.error {
  position:fixed;
  right: 0px;
  width: 240px;
  z-index: 9999999999;
}
</style>