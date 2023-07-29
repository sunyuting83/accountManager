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
          size="small"
          :pagination={pageSize:20}
          :hideOnSinglePage="true"
        >
          <template v-slot:bodyCell="{column,record}">
            <template v-if="column.dataIndex==='NewStatus'">
              <a-tag color="green" v-if="record.NewStatus == 0">正常</a-tag>
              <a-tag color="red" v-if="record.NewStatus == 1">整单退款中</a-tag>
              <a-tag color="purple" v-if="record.NewStatus == 2">单号退款中</a-tag>
              <a-tag color="orange" v-if="record.NewStatus == 3">已整单退款</a-tag>
              <a-tag color="orange" v-if="record.NewStatus == 4">已单号退款</a-tag>
            </template>
            <template v-if="column.dataIndex==='Active'">
              <a-button type="primary" size="small" ghost @click="() => {pushOrderDetail(record.ID)}">订单详情</a-button>
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
import { GetOrdersList } from '../../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'
const router = useRouter()
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const columns = [
  {
    title: '订单号',
    dataIndex: 'OrderCode',
  },
  {
    title: '金额',
    dataIndex: 'Coin',
    customRender: function (t: any) {
      return `￥${t.value}`
    }
  },
  {
    title: '订单状态',
    dataIndex: 'NewStatus',
  },
  {
    title: '创建时间',
    dataIndex: 'CreatedAt',
    customRender: function (t: any) {
      return foramTime(t.value)
    }
  },
  {
    title: '更新时间',
    dataIndex: 'UpdatedAt',
    customRender: function (t: any) {
      return foramTime(t.value)
    }
  },
  {
    title: '备注',
    dataIndex: 'Remarks',
  },
  {
    title: '操作',
    dataIndex: 'Active',
  }
];

interface OrderDatas {
  OrderCode: string;
  NewStatus: number;
  Coin: number;
  CoinUsersID: number;
  ID: number;
  Remarks: string;
  UpdatedAt: number;
  CreatedAt: number;
}

interface OrderResponse {
  status: number
  data: Array<OrderDatas>
  total: number
}
const dataState = ref<OrderResponse>({
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
  title: "订单管理",
  subtitle: '订单管理详情',
  routers: [
    {
      path:'main',
      breadcrumbName: '订单管理'
    },
    {
      path:'main',
      breadcrumbName: '订单管理'
    }
  ],
})
onMounted(() => {
  getOrders()
})

const pushOrderDetail = (id: number) => {
  router.push({
    'name': 'OrdersDetail',
    'params': {
      'order_id': id,
    }
  })
}

const getOrders = async(page:string = "1", pageSize:string = "20",) => {
  state.value.loading = true
  const params = {
    page: page,
    limit: pageSize
  }
  const data = await GetOrdersList(params)
  if (data.status == 0) {
    state.value.loading = false
    dataState.value = data as OrderResponse
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


interface State {
  loading: boolean;
}
const state = ref<State>({
  loading: false,
})

</script>
