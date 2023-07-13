<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}">
      <div class="content">
        <a-row :style="{'margin-bottom': '1rem'}">
          <a-col :span="12" v-if="gameState.status == 0">
            <a-space>
              <span>
                <a-button type="link" :disabled="gameID == 0 ? true : false " @click="()=>{changeGame(0)}">所有游戏</a-button>
              </span>
              <span v-for="(item) in gameState.data" :key="item.ID">
                <a-button type="link" :disabled="item.ID == gameID ? true : false " @click="()=>{changeGame(item.ID)}">{{item.GameName}}</a-button>
              </span>
            </a-space>
          </a-col>
          <a-col :span="12" :style="{'text-align': 'right'}">
            <a-space v-if="state.selectedRowKeys.length > 0">
              总价：{{state.total}} 
              <a-button type="primary" :disabled="state.selectedRowKeys.length > 0 ? false : true " @click="pushCart">加入购物车</a-button>
              <a-button type="primary" danger :disabled="state.selectedRowKeys.length > 0 ? false : true " @click="postCart">直接购买</a-button>
            </a-space>
          </a-col>
        </a-row>
        <a-empty :image="simpleImage" v-if="dataState.data.length == 0" />
        <a-table
          v-else
          :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
          :columns="columns"
          :data-source="dataState.data"
          :loading="state.loading"
          size="small"
          :pagination="false"
        />
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
import { onMounted, computed, ref, h } from 'vue'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { Empty } from 'ant-design-vue';
import { GetProducts, GetGamesList, AddCart, PostOrders } from '../../../wailsjs/go/main/App'
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
}

interface ProductResponse {
  status: number;
  total: number;
  data: Array<ProductDatas>
}
interface GamesDatas {
  GameName: string;
  ID: number;
}

interface GamesResponse {
  status: number;
  total: number;
  data: Array<GamesDatas>
}
const dataState = ref<ProductResponse>({
  status: 1,
  total: 0,
  data: [],
})

const gameState = ref<GamesResponse>({
  status: 0,
  total: 0,
  data: [
    {
      GameName: "",
      ID: 0,
    }
  ],
});

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

const current = ref(1)
const gameID = ref(0)

const pageHeader = ref<PageHeaderData>({
  title: "产品列表",
  subtitle: '产品列表详细信息',
  routers: [
    {
      path:'main',
      breadcrumbName: '产品列表'
    },
    {
      path:'main',
      breadcrumbName: '产品列表'
    }
  ],
})
onMounted(() => {
  getProducts()
})

const changeGame = (id:number) => {
  gameID.value = id
  getProducts("1","20",String(id))
}

const getProducts = async(page:string = "1", pageSize:string = "20", gameid:string = String(gameID.value)) => {
  const games = await GetGamesList()
  if (games.status == 0) {
    // userState.value = data as User
    // console.log(data)
    gameState.value = games as GamesResponse
  }
  state.value.loading = true
  const params = {
    page: page,
    limit: pageSize,
    gameid: gameid,
  }
  const data = await GetProducts(params)
  if (data.status == 0) {
    // userState.value = data as User
    // console.log(data)
    dataState.value = data as ProductResponse
    state.value.loading = false
  }else {
    state.value.loading = false
    errNotification(data.message)
  }
}

const changePage = (page: number, pageSize: number) => {
  getProducts(String(page), String(pageSize), String(gameID.value))
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
  selectedRowKeys: Key[];
  loading: boolean;
  total: number;
}
const state = ref<State>({
  selectedRowKeys: [],
  loading: false,
  total: 0,
})

const postCart = async() => {
  state.value.loading = true
  let IDList: number[] = [];
  const d = state.value.selectedRowKeys
  d.forEach((k) => IDList = [...IDList, Number(k)])
  const data = await PostOrders(IDList)
  if (data.status == 0) {
    state.value.loading = false
    sucNotification(`购买成功，请转至订单详情查看 总价：${data.total} 余额： ${data.credit} 失败： ${data.FailedData.length}条`)
    getProducts()
  }else{
    state.value.loading = false
    errNotification(data.message)
  }
}


const pushCart = async() => {
  let IDList: number[] = [];
  const d = state.value.selectedRowKeys
  d.forEach((k) => IDList = [...IDList, Number(k)])
  // console.log(IDList)
  const data = await AddCart(IDList)
  // console.log(data)
  sucNotification(data.message)
}

const hasSelected = computed(() => state.value.selectedRowKeys.length > 0);

const toDecimal2 = (x: number) => { 
  var f = parseFloat(String(x)); 
    if (isNaN(f)) { 
      return 0 
    } 
    var f = Math.round(x*100)/100; 
    var s = f.toString(); 
    var rs = s.indexOf('.'); 
    if (rs < 0) { 
        rs = s.length; 
        s += '.'; 
    } 
    while (s.length <= rs + 2) { 
      s += '0'; 
    } 
    return Number(s)
} 

const makeTotal = () => {
  let total = 0
  dataState.value.data.map((e) => {
    state.value.selectedRowKeys.map(el => {
      if (el == e.ID) total += e.Price
    })
  })
  state.value.total = toDecimal2(total)
}

const onSelectChange = (selectedRowKeys: Key[]) => {
  // console.log('selectedRowKeys changed: ', selectedRowKeys);
  state.value.selectedRowKeys = selectedRowKeys
  makeTotal()
}

</script>
