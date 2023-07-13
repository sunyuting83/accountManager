<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}" v-if="!state.loading">
      <div class="content">
        <a-row :style="{'margin-bottom': '1rem'}" v-if="dataState.status == 0" justify="space-between" align="middle">
          <a-col :span="6">
            <span v-if="state.total !== 0">总价：{{state.total}}</span>
          </a-col>
          <a-col :span="6" :style="{'text-align': 'right'}">
            <a-space>
              <a-button type="primary" :disabled="dataState.data.length > 0 ? false : true " @click="cleanCart">清空购物车</a-button>
              <a-button type="primary" danger :disabled="dataState.data.length > 0 ? false : true " @click="postCart">直接购买</a-button>
            </a-space>
          </a-col>
        </a-row>
        <a-empty :image="simpleImage" v-if="dataState.data.length == 0" />
        <a-table
          v-else
          :columns="columns"
          :data-source="dataState.data"
          :loading="state.loading"
          size="small"
          :pagination={pageSize:20}
          :hideOnSinglePage="true"
        >
          <template v-slot:bodyCell="{column,record}">
            <template v-if="column.dataIndex==='Status'">
              <a-tag color="green" v-if="record.SellStatus == 1">待出售</a-tag>
              <a-tag color="default" v-else>已出售</a-tag>
            </template>
            <template v-if="column.dataIndex==='Delete'">
              <a-button type="primary" size="small" ghost @click="() => {delOneCart(record)}">删除</a-button>
            </template>
          </template>
        </a-table>
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
import { GetCart, CleanCart, DeleteCart, PostOrders } from '../../../wailsjs/go/main/App'
type Key = string | number;

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const columns = [
  {
    title: '帐号',
    dataIndex: 'Account',
  },
  {
    title: '游戏名称',
    dataIndex: 'GameName',
  },
  {
    title: '金币',
    dataIndex: 'Gold',
  },
  {
    title: '炮台',
    dataIndex: 'Multiple',
  },
  {
    title: '炮台',
    dataIndex: 'Multiple',
  },
  {
    title: '狂暴',
    dataIndex: 'Crazy',
  },
  {
    title: '冰冻',
    dataIndex: 'Cold',
  },
  {
    title: '钻石',
    dataIndex: 'Diamond',
  },
  {
    title: '瞄准',
    dataIndex: 'Precise',
  },
  {
    title: '价格',
    dataIndex: 'Price',
  },
  {
    title: '状态',
    dataIndex: 'Status',
  },
  {
    title: '删除',
    dataIndex: 'Delete',
  }
];

interface ProductDatas {
  Account: string;
  Cold: number;
  Cover: string;
  Crazy: number;
  Diamond: number;
  GameID: number;
  GameName: string;
  Gold: string;
  ID: number;
  Multiple: number;
  Precise: number;
  Price: number;
  Remarks: string;
  UpdatedAt: number;
  key: number;
  SellStatus: number;
}

interface ProductResponse {
  status: number;
  data: Array<ProductDatas>
}
const dataState = ref<ProductResponse>({
  status: 1,
  data: [],
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
  title: "购物车",
  subtitle: '购物车详细信息',
  routers: [
    {
      path:'main',
      breadcrumbName: '购物车'
    },
    {
      path:'main',
      breadcrumbName: '购物车'
    }
  ],
})
onMounted(() => {
  getCarts()
})


const makeTotal = (data: ProductDatas[]) => {
  let total = 0
  data.map((e) => {
    total += e.Price
  })
  return total
}

const makeIDs = () => {
  const data = dataState.value.data
  let list: number[] = []
  data.map((e) => {
    list = [...list, e.ID]
  })
  return list
}

const getCarts = async() => {
  state.value.loading = true
  const data = await GetCart()
  if (data.status == 0) {
    dataState.value = data as ProductResponse
    state.value.total = makeTotal(data.data)
    state.value.loading = false
  }else {
    state.value.loading = false
    errNotification(data.message)
  }
}

const cleanCart = async() => {
  state.value.loading = true
  const data = await CleanCart()
  dataState.value = data as ProductResponse
  state.value.loading = false
  sucNotification("清空成功")
}


const errNotification = (text: string) => {
  notification.open({
    message: "发生错误",
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}

const sucNotification = (text: string) => {
  notification.open({
    message: "成功",
    description:
    text,
    icon: () => h(CheckCircleOutlined, { style: 'color: #389e0d' }),
  });
}

const delOneCart = async(e: ProductDatas) => {
  state.value.loading = true
  const data = await DeleteCart(e.ID)
  if (data.status == 0) {
    dataState.value.data = dataState.value.data.filter((es) => { return es.ID !== e.ID})
    state.value.loading = false
  }else {
    state.value.loading = false
    errNotification(data.message)
  }
}

const postCart = async() => {
  state.value.loading = true
  const ids = makeIDs()
  const data = await PostOrders(ids)
  if (data.status == 0) {
    state.value.loading = false
    sucNotification(`购买成功，请转至订单详情查看 总价：${data.total} 余额： ${data.credit} 失败： ${data.FailedData.length}条`)
    getCarts()
  }else{
    state.value.loading = false
    errNotification(data.message)
  }
}

interface State {
  selectedRowKeys: Key[];
  loading: boolean;
  total: number;
}
const state = ref<State>({
  selectedRowKeys: [],
  loading: false,
  total: 0,
})

</script>
