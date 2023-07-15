import { createRouter, createWebHistory } from 'vue-router'
import AppIndex from '../components/Index.vue'
import Login from '../components/Auth/Login.vue'
import Main from '../components/Main/Main.vue'
import Users from '../components/Main/Users.vue'
import Product from '../components/Main/Product.vue'
import Cart from '../components/Main/Cart.vue'
import Order from '../components/Main/Order.vue'
import OrdersDetail from '../components/Main/OrdersDetail.vue'


const routerHistory = createWebHistory()

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
        component: Users,
      },
      {
        path: 'product',
        name: 'product',
        component: Product,
      },
      {
        path: 'cart',
        name: 'cart',
        component: Cart,
      },
      {
        path: 'order',
        name: 'order',
        component: Order,
      },
      {
        path: 'OrdersDetail:order_id',
        name: 'OrdersDetail',
        component: OrdersDetail,
      }
    ]
  },
]


let router = createRouter({
  history: routerHistory,
  routes: constantRoutes,
})

export default router