import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '../components/Index.vue'
import Login from '../components/Auth/Login.vue'
import Main from '../components/Main/Main.vue'
import Admin from '../components/Main/Admin.vue'
import SendCoin from '../components/Main/SendCoin.vue'

const constantRoutes = [
  {
    path: '/',
    name: 'index',
    component: AppIndex
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/main',
    name: 'main',
    component: Main,
    children: [
      {
        path: 'user',
        name: 'user',
        component: Admin,
      },
      {
        path: 'sendcoin',
        name: 'sendcoin',
        component: SendCoin,
      }
    ]
  },
]
const routerHistory = createWebHistory()


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router