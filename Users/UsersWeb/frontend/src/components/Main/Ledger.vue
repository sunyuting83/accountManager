<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}" v-if="!state.loading">
      <div class="content">
        <a-empty :image="simpleImage" v-if="dataState.ledger.length == 0" />
        <a-table
          v-else
          :columns="columns"
          :data-source="dataState.ledger"
          :loading="state.loading"
          size="small"
          :rowKey="(record: LedgerDatas) => record.ID"
          :pagination={pageSize:20}
          :hideOnSinglePage="true"
        >
          <template v-slot:bodyCell="{column,record}">
            <template v-if="column.dataIndex==='Status'">
              <a-tag color="green" v-if="record.Status == 0">充值</a-tag>
              <a-tag color="red" v-if="record.Status == 1">转账给他人</a-tag>
              <a-tag color="purple" v-if="record.Status == 2">收到转账</a-tag>
              <a-tag color="orange" v-if="record.Status == 3">消费</a-tag>
              <a-tag color="orange" v-if="record.Status == 4">退款</a-tag>
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
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { Empty } from 'ant-design-vue';
import { GetLedger } from '../../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'
const router = useRouter()

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const columns = [
  {
    title: '用途',
    dataIndex: 'StatusName',
  },
  {
    title: '金额',
    dataIndex: 'Coin',
    customRender: function (t: any) {
      let status = ""
      switch (t.record.Status) {
      case 0:
        status = "+ "
        break
      case 1:
        status = "- "
        break
      case 2:
        status = "+ "
        break
      case 3:
        status = "- "
        break
      default:
        status = ""
        break
      }
      return `${status}￥${t.value}`
    }
  },
  {
    title: '状态',
    dataIndex: 'Status',
  },
  {
    title: '转账人',
    dataIndex: 'FormCoinUsers',
  },
  {
    title: '订单号',
    dataIndex: 'OrderCode',
  },
  {
    title: '创建时间',
    dataIndex: 'CreatedAt',
    customRender: function (t: any) {
      return foramTime(t.value)
    }
  }
];

interface LedgerDatas {
  ID: number;
  Coin: number;
  Status: number;
  FormCoinUsers: string;
  OrderCode: string;
  CreatedAt: number;
}

interface LedgerResponse {
  status: number
  ledger: Array<LedgerDatas>
  total: number
}
const dataState = ref<LedgerResponse>({
  status: 1,
  ledger: [],
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
  title: "账单",
  subtitle: '账单详情',
  routers: [
    {
      path:'main',
      breadcrumbName: '账单'
    },
    {
      path:'main',
      breadcrumbName: '账单'
    }
  ],
})
onMounted(() => {
  getLedger()
})


const getLedger = async(page:string = "1", pageSize:string = "50",) => {
  state.value.loading = true
  const params = {
    page: page,
    limit: pageSize
  }
  const data = await GetLedger(params)
  if (data.status == 0) {
    state.value.loading = false
    dataState.value = data as LedgerResponse
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

interface State {
  loading: boolean;
}
const state = ref<State>({
  loading: false,
})

</script>
