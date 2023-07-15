<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}" v-if="!state.loading">
      <div class="content">
        <a-empty :image="simpleImage" v-if="dataState.status != 0" />
        <div v-else>
          <a-descriptions title="订单详情" bordered>
            <a-descriptions-item label="订单号">{{dataState.data.OrderCode}}</a-descriptions-item>
            <a-descriptions-item label="订单状态" :span="2">
              <a-badge status="processing" v-if="dataState.data.Status == 0" text="正常" />
              <a-badge status="error" v-if="dataState.data.Status == 1" text="整单退款中" />
              <a-badge status="warning" v-if="dataState.data.Status == 2" text="单号退款中" />
              <a-badge status="success" v-if="dataState.data.Status == 3" text="已整单退款" />
              <a-badge status="success" v-if="dataState.data.Status == 4" text="已单号退款" />
            </a-descriptions-item>
            <a-descriptions-item label="下单时间">{{foramTime(dataState.data.CreatedAt)}}</a-descriptions-item>
            <a-descriptions-item label="更新时间" :span="2">{{foramTime(dataState.data.UpdatedAt)}}</a-descriptions-item>
            <a-descriptions-item label="订单金额" :span="3">${{dataState.data.Coin}}</a-descriptions-item>
            <a-descriptions-item label="备注" :span="3">
              {{dataState.data.Remarks}}
            </a-descriptions-item>
            <a-descriptions-item label="操作">
              <a-popover v-model:open="state.visible" placement="topLeft" arrow-point-at-center title="确定退款吗？" trigger="click">
                <template #content>
                  <a-textarea v-model:value="remarks" placeholder="请输入备注" size="small" allow-clear />
                  <a-space :style="{'margin-top':'10px'}">
                    <a-button type="primary" size="small" small @click="orderRefund">确定</a-button>
                  </a-space>
                </template>
                <a-button type="primary">整单退款</a-button>
              </a-popover>
            </a-descriptions-item>
          </a-descriptions>
          <a-table
            sticky
            :style="{'margin-top': '1rem'}"
            :columns="columns"
            :data-source="dataState.data.Accounts"
            :loading="state.loading"
            size="small"
            :scroll="{ x: 1500 }">
            <template v-slot:bodyCell="{column,record}">
              <template v-if="column.dataIndex==='Active'">
                <a-popover v-model:open="state.visible" placement="topLeft" arrow-point-at-center title="确定退款吗？" trigger="click">
                  <template #content>
                    <a-textarea v-model:value="remarks" placeholder="请输入备注" size="small" allow-clear />
                    <a-space :style="{'margin-top':'10px'}">
                      <a-button type="primary" size="small" small>确定</a-button>
                    </a-space>
                  </template>
                  <a-button type="primary" ghost>退款</a-button>
                </a-popover>
              </template>
            </template>
          </a-table>
        </div>
      </div>
    </div>
  </a-layout-content>
</template>
<script lang="ts" setup>
import PageHeader from './PageHeader.vue'
import { onMounted, ref, h } from 'vue'
import { InfoCircleOutlined, CheckCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { notification, Empty } from 'ant-design-vue'
import { GetOrdersDetail, OrderRefund } from '../../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'
const router = useRouter()
const remarks = ref<string>('')

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const OrderID = router.currentRoute.value.params.order_id

const columns = [
  {
    title: '帐号',
    dataIndex: 'Account',
    fixed: 'left',
    width: 150
  },
  {
    title: '密码',
    dataIndex: 'Password',
    fixed: 'left',
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
    title: '单价',
    dataIndex: 'Price',
  },
  {
    title: '更新时间',
    dataIndex: 'UpdatedAt',
    customRender: function (t: any) {
      return foramTime(t.value)
    },
    width: 150
  },
  {
    title: '备注',
    dataIndex: 'Remarks',
  },
  {
    title: '操作',
    dataIndex: 'Active',
    fixed: 'right',
  }
];

interface Account {
  Account: string;
  Cold: number;
  Cover: string;
  Crazy: number;
  Diamond: number;
  GameName: string;
  Gold: string;
  ID: number;
  Multiple: number;
  Password: string;
  Precise: number;
  Price: number;
  Remarks: string;
  Status: number;
  UpdatedAt: number;
}

interface OrderDatas {
  OrderCode: string;
  Status: number;
  Coin: number;
  ID: number;
  Remarks: string;
  CreatedAt: number;
  UpdatedAt: number;
  Accounts: Array<Account>
}

interface OrderResponse {
  status: number
  data: OrderDatas
}
const dataState = ref<OrderResponse>({
  status: 1,
  data: {
    OrderCode: "",
    Status: 0,
    Coin: 0,
    ID: 0,
    Remarks: "",
    CreatedAt: 0,
    UpdatedAt: 0,
    Accounts: []
  },
})


interface PageHeaderData {
  title: string;
  subtitle: string;
  routers: Array<
    {
      path: string;
      breadcrumbName: string;
    }
  >,
  back: boolean;
}

const pageHeader = ref<PageHeaderData>({
  title: "订单管理",
  subtitle: '订单详情',
  routers: [
    {
      path:'main',
      breadcrumbName: '订单详情'
    },
    {
      path:'main',
      breadcrumbName: '订单详情'
    }
  ],
  back: true
})
onMounted(() => {
  getOrders()
})

const getOrders = async() => {
  state.value.loading = true
  const params = {
    order_id: OrderID,
  }
  const data = await GetOrdersDetail(params)
  if (data.status == 0) {
    state.value.loading = false
    dataState.value = data as OrderResponse
    state.value.loading = false
  }else {
    state.value.loading = false
    errNotification(data.message)
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

const orderRefund = async() => {
  const params = {
    id: OrderID,
    remarks: remarks.value
  }
  const data = await OrderRefund(params)
  if (data.status == 0) {
    remarks.value = ""
    sucNotification(data.message)
    getOrders()
  }else {
    remarks.value = ""
    errNotification(data.message)
  }
}

interface State {
  loading: boolean;
  visible: boolean;
}
const state = ref<State>({
  loading: false,
  visible: false
})

</script>
