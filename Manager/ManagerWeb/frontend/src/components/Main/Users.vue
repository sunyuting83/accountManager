<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}" v-if="!state.loading">
      <div class="content">
        <a-empty :image="simpleImage" v-if="dataState.data.length == 0" />
        <a-table
          v-else
          :columns="columns"
          :data-source="dataState.data"
          :loading="state.loading"
          :rowKey="(record: UserDatas) => record.ID"
          size="small"
          :hideOnSinglePage="true"
          :pagination="false"
        >
          <template v-slot:bodyCell="{column,record}">
            <template v-if="column.dataIndex==='NewStatus'">
              <a-tag :color="record.NewStatus === 0 ? 'green' : 'red'" @click="() => {upStatus(record.ID)}">{{ record.NewStatus === 0 ? '正常' : '锁定' }}</a-tag>
            </template>
            <template v-if="column.dataIndex==='Active'">
              <a-popover placement="topRight" arrow-point-at-center title="发币" v-model:open="record.sendCoin" trigger="click">
                <template #content>
                  <a-space>
                    <a-input-number v-model:value="state.coin" :min="1" />
                    <a-button type="primary" @click="() => {sendCoinToUser(record.ID)}">发币</a-button>
                  </a-space>
                </template>
                <a-button type="primary" size="small" ghost>发币</a-button>
              </a-popover>
            </template>
          </template>
        </a-table>
        <a-pagination
          v-model:current="current"
          :total="dataState.total"
          :pageSize="20"
          :hideOnSinglePage="true"
          :style="{'margin-top': '15px'}"
          @change="changePage"
          :showSizeChanger="false">
          <template #itemRender="{ type, originalElement }">
            <a v-if="type === 'prev'">上一页</a>
            <a v-else-if="type === 'next'">下一页</a>
            <component :is="originalElement" v-else></component>
          </template>
        </a-pagination>
      </div>
    </div>
  </a-layout-content>
</template>
<script lang="ts" setup>
import PageHeader from './PageHeader.vue'
import { onMounted, ref, h } from 'vue'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { Empty } from 'ant-design-vue';
import { SendCoinToUser, UsersList, UpStatusUser } from '../../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'
const router = useRouter()
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const columns = [
  {
    title: '用户名',
    dataIndex: 'UserName',
  },
  {
    title: '余额',
    dataIndex: 'Coin',
    customRender: function (t: any) {
      return `￥${t.value}`
    }
  },
  {
    title: '钱包地址',
    dataIndex: 'WalletAddress'
  },
  {
    title: '用户状态',
    dataIndex: 'NewStatus',
  },
  {
    title: '操作',
    dataIndex: 'Active',
  }
];

interface UserDatas {
  UserName: string;
  IPAddress: string;
  LocalAddress: string;
  Email: string;
  PhoneNumber: string;
  NewStatus: number;
  Coin: number;
  ID: number;
  UpdatedAt: number;
  WalletAddress: string;
  sendCoin: boolean;
}

interface UsersResponse {
  status: number
  data: Array<UserDatas>
  total: number
}
const dataState = ref<UsersResponse>({
  status: 1,
  data: [],
  total: 0,
})


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
  title: "用户管理",
  subtitle: '用户管理详情',
  routers: [
    {
      path:'main',
      breadcrumbName: '用户管理'
    },
    {
      path:'main',
      breadcrumbName: '用户管理'
    }
  ],
})
onMounted(() => {
  getUsersList()
})

const getUsersList = async(page:string = "1", pageSize:string = "20",) => {
  state.value.loading = true
  const params = {
    page: page,
    limit: pageSize
  }
  const data = await UsersList(params)
  if (data.status == 0) {
    state.value.loading = false
    dataState.value = data as UsersResponse
    state.value.loading = false
  }else {
    state.value.loading = false
    errNotification(data.message)
    if (data.message == "403") {
      router.push({
        'name': 'login',
      })
    }
  }
  
}

const current = ref(1)

const changePage = (page: number, pageSize: number) => {
  getUsersList(String(page), String(pageSize))
}

const errNotification = (text: string) => {
  notification.open({
    message: "发生错误",
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}


const sendCoinToUser = async(id: number) => {
  dataState.value.data.map((e: UserDatas) => {
    if (e.ID == id) {
      e.sendCoin = true
      return 
    }
  })
  const params = {
    userid: String(id),
    coin_count: String(state.value.coin)
  }
  const data = await SendCoinToUser(params)
  if (data.status == 0) {
    dataState.value.data.map((e: UserDatas) => {
      if (e.ID == id) {
        e.sendCoin = false
        e.Coin += state.value.coin
        return 
      }
    })
    state.value.coin = 100
    sucNotification(data.message)
  }else{
    state.value.coin = 100
    dataState.value.data.map((e: UserDatas) => {
      if (e.ID == id) {
        e.sendCoin = false
        return 
      }
    })
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

const upStatus = async(id: number) => {
  const params = {
    id: id
  }
  const data = await UpStatusUser(params)
  if (data.status == 0) {
    dataState.value.data.map((e: UserDatas) => {
      if (e.ID == id) {
        if (e.NewStatus == 0) {
          e.NewStatus = 1
          return e
        }else{
          e.NewStatus = 0
          return e
        }
      }
    })
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


interface State {
  loading: boolean;
  coin: number;
  sendcoin: boolean;
}
const state = ref<State>({
  loading: false,
  coin: 100,
  sendcoin: false
})

</script>
