<template>
  <nav class="navbar is-spaced" :class="'is-'+skin">
    <div class="navbar-brand">
      <span class="navbar-item brand-text">
        <span class="icon-text">
          <span class="icon">
            <i class="fa fa-ils"></i>
          </span>
          <span>Account Manage</span>
        </span>
      </span>
      <span role="button" class="navbar-burger" :class="isActive ? 'is-active' : ''" aria-label="menu" aria-expanded="false" data-target="navbarBasicExample" @click="toogleMenu">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </span>
    </div>
    <div id="navMenu" class="navbar-menu" :class="isActive ? 'is-active' : ''">
      <div class="navbar-start">
        <router-link class="navbar-item" :class="path === 'userlist'?'is-active':''" to="/userlist">
          用户管理
        </router-link>
        <router-link class="navbar-item" :class="path === 'project' || path === 'userProject' ? 'is-active' : ''" to="/project">
          项目管理
        </router-link>
        <router-link class="navbar-item" v-if="userid === '1'" :class="path === 'gameslist'?'is-active':''" to="/gameslist">
          游戏管理
        </router-link>
        <router-link class="navbar-item" v-if="userid === '1'" :class="path === 'AllDraw'?'is-active':''" to="/AllDraw">
          题号管理
        </router-link>
        <a class="navbar-item" v-if="path === 'account'" :class="path === 'account' ? 'is-active' : ''">
        帐号管理
        </a>
        <a class="navbar-item" v-if="path === 'accountDraw'" :class="path === 'accountDraw'?'is-active':''">
          未提取帐号管理
        </a>
        <a class="navbar-item" v-if="path === 'accountDrawed'" :class="path === 'accountDrawed'?'is-active':''">
          已提取帐号管理
        </a>
        <a class="navbar-item" v-if="path === 'drawLog' || path ===  'drawData'" :class="path === 'drawLog' || path === 'drawData' ? 'is-active' : ''">
          提取记录
        </a>
      </div>
      <div class="navbar-end">
        <span v-for="item in skinlist" :key="item" class="item" :class="skin === item ? item + ' has' : item" @click="()=>{setSkin(item)}"></span>
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
                <router-link class="navbar-item" to="/adminlist" v-if="openModal.username === 'admin'">
                  管理员管理
                </router-link>
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
      :ShowMessage="ShowMessage" />
    <NotIfication
      :showData="openerr" />
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
      isActive: false,
      userid: "0",
      skin: "",
      skinlist: [
        'primary',
        'link',
        'info',
        'success',
        'warning',
        'danger',
        'dark',
        'black',
        'light'
      ]
    })
    states.openModal.username = localStorage.getItem('user')
    states.userid = localStorage.getItem('userid')
    states.skin = localStorage.getItem('header_style')
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
    const toogleMenu = () => {
      states.isActive = !states.isActive
    }
    const setSkin = (es) => {
      localStorage.setItem('header_style', es)
      states.skin = es
    }
    return {
      ...toRefs(states),
      LogOut,
      router,
      showModel,
      ShowMessage,
      toogleMenu,
      setSkin
    }
  }
})
</script>
<style scoped>
.navbar-menu .item {
  display: inline-block;
  width: 15px;
  height: 15px;
  margin-left: 1rem;
}
.navbar-menu .item.has {
  border: 1px solid #f5f5f5
}
.navbar-menu .item.primary {
  background: #00d1b2;
}
.navbar-menu .item.link {
  background: #485fc7;
}
.navbar-menu .item.info {
  background: #3e8ed0;
}
.navbar-menu .item.success {
  background: #48c78e;
}
.navbar-menu .item.warning {
  background: #ffe08a;
}
.navbar-menu .item.danger {
  background: #f14668;
}
.navbar-menu .item.dark {
  background: #363636;
}
.navbar-menu .item.black {
  background: #0a0a0a;
}
.navbar-menu .item.light {
  background: #f5f5f5;
}
</style>