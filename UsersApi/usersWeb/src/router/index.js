import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '@/components/Index/Index'
import ProjectList from '@/components/Project/ProjectList'
import AccountList from '@/components/Account/AccountList'

const routerHistory = createWebHistory()

const constantRoutes = [
  {
    path: '/',
    name: 'index',
    component: AppIndex
  },
  {
    path: '/project',
    name: 'project',
    component: ProjectList
  },
  {
    path: '/account/:key',
    name: 'account',
    component: AccountList
  }
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router