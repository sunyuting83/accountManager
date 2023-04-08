<template>
  <div id="app">
    <router-view v-if="isRouterAlive" />
  </div>
</template>

<script>
export default {
  name: 'App',
  mounted(){
    localStorage.setItem('header_style','link')
  },
  provide() {
    return{
      reload: this.reload	
    }
  },
  data(){
    return{
      isRouterAlive: true	
    }
  },
  methods: {
    reload(){
      this.isRouterAlive = false
      this.$nextTick(() =>{
        this.isRouterAlive = true		
      })
    }
  }
}
</script>
