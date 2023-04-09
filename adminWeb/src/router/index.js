import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '@/components/Index/Index'
import AdminList from '@/components/Admin/AdminList'
import GamesList from '@/components/Games/GamesList'
import UserList from '@/components/User/UserList'
import ProjectList from '@/components/Project/ProjectList'
import UserProjectList from '@/components/Project/UserProjectList'
import AccountList from '@/components/Account/AccountList'
import AccountDraw from '@/components/AccountDraw/AccountList'
import AccountDrawed from '@/components/AccountDrawed/AccountList'
import AccountFiled from '@/components/AccountFiled/FiledList'
import DrawLogs from '@/components/DrawedLog/DrawedLog'
import DraweData from '@/components/DrawedLog/DraweData'
import AllDraw from '@/components/AllDraw/AllDraw'

const routerHistory = createWebHistory()

const constantRoutes = [
  {
    path: '/',
    name: 'index',
    component: AppIndex
  },
  {
    path: '/adminlist',
    name: 'adminlist',
    component: AdminList
  },
  {
    path: '/gameslist',
    name: 'gameslist',
    component: GamesList
  },
  {
    path: '/userlist',
    name: 'userlist',
    component: UserList
  },
  {
    path: '/project',
    name: 'project',
    component: ProjectList
  },
  {
    path: '/userProject/:userid',
    name: 'userProject',
    component: UserProjectList
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
  },
  {
    path: '/drawLog/:key',
    name: 'drawLog',
    component: DrawLogs
  },
  {
    path: '/drawData/:id',
    name: 'drawData',
    component: DraweData
  },
  {
    path: '/AllDraw',
    name: 'AllDraw',
    component: AllDraw
  }
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router