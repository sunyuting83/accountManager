<template>
  <a-layout has-sider style="min-height: 100vh">
    <Menu :Trigger="collapsed" />
    <a-layout :style="{ marginLeft: collapsed?'80px':'200px' }">
      <Header :Trigger="collapsed" :ChangeIt="ClickCollapsed" />
      <router-view />
      <Footer />
    </a-layout>
  </a-layout>
</template>
<script lang="ts">
import { defineComponent, reactive, toRefs, onMounted } from 'vue'
import Menu from './Menu.vue'
import Header from './Header.vue'
import Footer from './Footer.vue'
import { useRouter } from 'vue-router'
import { CheckLogin } from '../../../wailsjs/go/main/App'
export default defineComponent({
  components: {
    Menu,
    Header,
    Footer,
  },
  setup() {
    const router = useRouter()
    onMounted(async() => {
      const data = await CheckLogin()
      if (data.status === 1) {
        router.push({
          'name': 'login',
        })
      }else{
        router.push({
          'name': 'user',
        })
      }
    })
    let states = reactive({
      collapsed: false
    })
    const ClickCollapsed = () => {
      states.collapsed = !states.collapsed
    }
    return {
      ...toRefs(states),
      ClickCollapsed
    };
  },
});
</script>
<style>
.site-layout .site-layout-background {
  background: #fff;
}
.ant-statistic-content {
  font-size: 20px;
  line-height: 28px;
}
@media (max-width: 576px) {
  .content {
    display: block;
  }

  .main {
    width: 100%;
    margin-bottom: 12px;
  }

  .extra {
    width: 100%;
    margin-left: 0;
    text-align: left;
  }
}
</style>
