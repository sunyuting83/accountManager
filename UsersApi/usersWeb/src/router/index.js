import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '@/components/Index/Index'
import ProjectList from '@/components/Project/ProjectList'
import AccountList from '@/components/Account/AccountList'
import AccountDraw from '@/components/AccountDraw/AccountList'
import AccountDrawed from '@/components/AccountDrawed/AccountList'
import AccountFiled from '@/components/AccountFiled/FiledList'

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
  },
  {
    path: '/accountDraw/:key/:type',
    name: 'accountDraw',
    component: AccountDraw
  },
  {
    path: '/accountDrawed/:key',
    name: 'accountDrawed',
    component: AccountDrawed
  },
  {
    path: '/accountFiled/:key',
    name: 'accountFiled',
    component: AccountFiled
  }
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router