<template>
  <nav class="navbar is-white">
    <div class="container">
      <div class="navbar-brand">
        <span class="navbar-item brand-text">
          <img :src="logo">Account Manage
        </span>
      </div>
      <div id="navMenu" class="navbar-menu">
        <div class="navbar-start">
          <router-link class="navbar-item" :class="path === 'project'?'is-active':''" to="/project">
            项目管理
          </router-link>
          <a class="navbar-item" v-if="path === 'account'" :class="path === 'account'?'is-active':''">
            帐号管理
          </a>
          <a class="navbar-item" v-if="path === 'accountDraw'" :class="path === 'accountDraw'?'is-active':''">
            未提取帐号管理
          </a>
          <a class="navbar-item" v-if="path === 'accountDrawed'" :class="path === 'accountDrawed'?'is-active':''">
            已提取帐号管理
          </a>
        </div>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="field is-grouped">
            <div class="navbar-item has-dropdown is-hoverable">
              <span class="navbar-link">
                <span class="icon">
                  <i class="fa fa-user-circle-o"></i>
                </span>
                <span>
                用户中心
                </span>
              </span>
          
              <div class="navbar-dropdown">
                <div class="dropdown-item">
                  <p>当前登陆用户 <strong>{{openModal.username}}</strong> <br /><span class="is-size-7">努力刷号賺錢</span></p>
                </div>
                <hr class="dropdown-divider">
                <a class="navbar-item" @click="showModel">
                  修改密码
                </a>
                <hr class="navbar-divider">
                <a class="navbar-item" @click="LogOut">
                  退出登陆
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ChangePassword
      :showData="openModal"
      :ShowMessage="ShowMessage">
    </ChangePassword>
    <NotIfication
      :showData="openerr">
    </NotIfication>
  </nav>
</template>
<script>
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
import ChangePassword from "./ChangePassword"
import NotIfication from "@/components/Other/Notification"
import { reactive, toRefs, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
export default defineComponent({
  name: 'CardList',
  components: { ChangePassword, NotIfication },
  setup() {
    const router = useRouter()
    let states = reactive({
      logo: Config.images[2],
      path: router.currentRoute.value.name,
      openModal:{
        active: false,
        username: ""
      },
      openerr: {
        active: false,
        message: ""
      },
    })
    states.openModal.username = localStorage.getItem('user')
    // console.log(states.openModal.username)
    const LogOut =() =>{
      setStorage(false)
      router.push("/")
    }
    const showModel =() =>{
      states.openModal.active = true
    }
    const ShowMessage =(e) =>{
      states.openerr = e
    }
    return {
      ...toRefs(states),
      LogOut,
      router,
      showModel,
      ShowMessage
    }
  }
})
</script>

<style scoped>
nav.navbar {
  border-top: 4px solid #276cda;
}
</style>