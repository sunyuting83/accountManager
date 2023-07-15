<template>
  <a-layout-sider v-model:collapsed="Trigger" :trigger="null" collapsible :style="{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0, top: 0, bottom: 0 }">
    <div class="logo" />
    <a-menu
      v-model:selectedKeys="selectedKeys"
      v-model:openKeys="openKeys"
      theme="dark"
      mode="inline"
      @click="handleClick">
      <a-menu-item key="user">
        <span>
          <user-outlined />
          <span>用户中心</span>
        </span>
      </a-menu-item>
      <a-menu-item key="product">
        <desktop-outlined />
        <span>产品列表</span>
      </a-menu-item>
      <a-menu-item key="cart">
        <shopping-cart-outlined />
        <span>购物车</span>
      </a-menu-item>
      <a-menu-item key="order">
        <account-book-outlined />
        <span>订单管理</span>
      </a-menu-item>
      <a-menu-item key="logout">
        <export-outlined />
        <span>退出</span>
      </a-menu-item>
    </a-menu>
  </a-layout-sider>
</template>
<script lang="ts">
import {
  DesktopOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  AccountBookOutlined,
  ExportOutlined,
} from '@ant-design/icons-vue';
import { defineComponent, ref } from 'vue'
import { Logout } from '../../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'
import type { MenuProps } from 'ant-design-vue';
export default defineComponent({
  props: {
    Trigger: {
        type: Boolean,
        default: false
      }
  },
  components: {
    DesktopOutlined,
    UserOutlined,
    ShoppingCartOutlined,
    AccountBookOutlined,
    ExportOutlined,
  },
  setup() {
    const openKeys = ref<string[]>(['sub1']);
    const router = useRouter()
    const handleClick: MenuProps['onClick'] = e => {
      if (e.key == "logout") {
        clickLogout()
      }else{
        router.push({
          'name': String(e.key),
        })
      }
    };
    const clickLogout = async() => {
      // console.log("asdf")
      const data = await Logout()
      // console.log(data)
      if (data.status == 0) {
        router.push({
          'name': 'login',
        })
      }
    }
    return {
      selectedKeys: ref<string[]>(['user']),
      openKeys,
      clickLogout,
      handleClick
    };
  },
});
</script>
<style>
.trigger {
  font-size: 18px;
  line-height: 64px;
  padding: 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

.logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.3);
  margin: 16px;
  line-height: 32px;
  color: #f3f3f3
}
</style>
