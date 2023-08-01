<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}">
      <div class="content">
        <a-input-group compact>
          <a-input v-model:value="state.wallet" addon-before="钱包地址" style="width: calc(100% - 200px)" />
          <a-button type="primary" :disabled="state.wallet.length == 16 ? false : true " @click="getUserWithKey">搜索</a-button>
        </a-input-group>
        <a-skeleton active v-if="state.loading" />
        <a-descriptions bordered :style="{'margin-top': '1rem'}" v-if="state.userid !== 0">
          <a-descriptions-item :label="state.username">
            <a-space>
              <a-input-number v-model:value="state.coin" :min="1" />
              <a-button type="primary" :disabled="state.coin > 0 ? false : true " @click="sendCoinToUser">发币</a-button>
            </a-space>
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </div>
  </a-layout-content>
</template>
<script lang="ts" setup>
import PageHeader from './PageHeader.vue'
import { ref, h } from 'vue'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { Empty } from 'ant-design-vue';
import { useRouter } from 'vue-router'
import { GetUserWithKey, SendCoinToUser } from '../../../wailsjs/go/main/App'
const router = useRouter()


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
  title: "发币给用户",
  subtitle: '发币给用户详情',
  routers: [
    {
      path:'main',
      breadcrumbName: '发币给用户'
    },
    {
      path:'main',
      breadcrumbName: '发币给用户'
    }
  ],
})

const getUserWithKey = async() => {
  state.value.loading = true
  const params = {
    key: state.value.wallet
  }
  const data = await GetUserWithKey(params)
  if (data.status == 0) {
    state.value.username = data.data.UserName
    state.value.userid = data.data.ID
    state.value.loading = false
  }else{
    state.value.loading = false
    errNotification(data.message)
  }
}

const sendCoinToUser = async() => {
  const params = {
    userid: String(state.value.userid),
    coin_count: String(state.value.coin)
  }
  const data = await SendCoinToUser(params)
  if (data.status == 0) {
    sucNotification(data.message)
  }else{
    errNotification(data.message)
    if (data.message == "403") {
      router.push({
        'name': 'login',
      })
    }
  }
}

const sucNotification = (text: string) => {
  notification.open({
    message: "成功",
    description:
    text,
    icon: () => h(CheckCircleOutlined, { style: 'color: #389e0d' }),
  });
}

const errNotification = (text: string) => {
  notification.open({
    message: "发生错误",
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}

interface State {
  loading: boolean;
  wallet: string;
  coin: number;
  username: string;
  userid: number;
}
const state = ref<State>({
  loading: false,
  wallet: "",
  coin: 100,
  username: "",
  userid: 0
})

</script>
