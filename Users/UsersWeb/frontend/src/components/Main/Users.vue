<template>
  <a-layout-content :style="{background: '#fff' }" v-if="userState.status == 0">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}">
      <div class="content">
        <div class="main">
          <a-descriptions size="small" :column="2">
            <a-descriptions-item label="用户名">{{userState.users.UserName}}</a-descriptions-item>
            <a-descriptions-item label="钱包地址">{{userState.users.WalletAddress}}</a-descriptions-item>
            <a-descriptions-item label="E-mail">{{userState.users.Email}}</a-descriptions-item>
            <a-descriptions-item label="电话">{{userState.users.PhoneNumber}}</a-descriptions-item>
            <a-descriptions-item label="最后登陆时间">{{foramTime(userState.users.UpdatedAt)}}</a-descriptions-item>
            <a-descriptions-item label="最后登陆IP">{{userState.users.IPAddress}}</a-descriptions-item>
            <a-descriptions-item label="所在地">
              {{userState.users.LocalAddress}}
            </a-descriptions-item>
          </a-descriptions>
        </div>
        <div class="extra">
          <div
            :style="{
              display: 'flex',
              width: 'max-content',
              justifyContent: 'flex-end',
            }"
          >
            <a-statistic
              title="状态"
              :value="userState.users.NewStatus == 0 ? '正常' : '锁定'"
              :style="{
                marginRight: '32px',
              }"
            />
            <a-statistic title="余额" prefix="¥" :value="userState.users.Coin" />
          </div>
        </div>
      </div>
    </div>
  </a-layout-content>
</template>
<script lang="ts" setup>
import PageHeader from './PageHeader.vue'
import { onMounted, ref, h } from 'vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { GetUsers } from '../../../wailsjs/go/main/App'

interface PageHeaderData {
  title: string;
  subtitle: string;
  routers: Array<
    {
      path: string;
      breadcrumbName: string;
    }
  >
}

const pageHeader = ref<PageHeaderData>({
  title: "用户中心",
  subtitle: '用户详细信息',
  routers: [
    {
      path:'main',
      breadcrumbName: '用户中心'
    },
    {
      path:'main',
      breadcrumbName: '用户中心'
    }
  ],
});

interface Users {
  Coin: number;
  CreatedAt: number;
  Email: string;
  ID: number;
  IPAddress: string;
  LocalAddress: string;
  NewStatus: number;
  Password: string;
  PhoneNumber: string;
  UpdatedAt: number;
  UserName: string;
  WalletAddress: string;
}

interface User {
  status: number;
  message: string;
  users: Users;
}

const userState = ref<User>({
  status: 0,
  message: '',
  users: {
    Coin: 0,
    CreatedAt: 0,
    Email: '',
    ID: 0,
    IPAddress: '',
    LocalAddress: '',
    NewStatus: 0,
    Password: '',
    PhoneNumber: '',
    UpdatedAt: 0,
    UserName: '',
    WalletAddress: '',
  },
});
onMounted(async() => {
  const data = await GetUsers()
  if (data.status == 0) {
    userState.value = data as User
    // console.log(userState.value)
  }else {
    openNotification(data.message)
  }
})
const openNotification = (text: string) => {
  notification.open({
    message: '发生错误',
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}
const foramTime = (d: number) => {
  const date = new Date(d)
  const Y = date.getFullYear()
  const M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1)
  const D = date.getDate()
  const H = date.getHours() < 10 ? `0${date.getHours()}` : date.getHours()
  const minute = date.getMinutes() < 10 ? `0${date.getMinutes()}` : date.getMinutes()
  const S = date.getSeconds() < 10 ? `0${date.getSeconds()}` : date.getSeconds()
  return `${Y}-${M}-${D} ${H}:${minute}:${S}`
}
</script>
<style scoped>
.content {
  display: flex;
}
</style>
