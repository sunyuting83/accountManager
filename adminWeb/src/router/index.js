import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '@/components/Index/Index'
import AdminList from '@/components/Admin/AdminList'
import UserList from '@/components/User/UserList'
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
    path: '/adminlist',
    name: 'adminlist',
    component: AdminList
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
    path: '/account/:id',
    name: 'account',
    component: AccountList
  }
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router