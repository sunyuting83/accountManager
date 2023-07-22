<template>
  <a-drawer
    v-model:visible="open"
    class="custom-class"
    root-class-name="root-class-name"
    :root-style="{ color: 'blue' }"
    style="{color: red}"
    title="筛选"
    placement="right"
  >
    <a-form :model="formState" :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item label="选择游戏">
        <a-select
        ref="select"
          v-model:value="formState.gameid"
          style="width: 100%"
          @select="handleChange"
        >
          <a-select-option :value="item.ID" v-for="(item) in gameState.data" :key="item.ID">{{item.GameName}}</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="金币">
        <a-input-number v-model:value="formState.gold" addon-before=">=" addon-after="亿"></a-input-number>
      </a-form-item>
      <a-form-item label="炮台">
        <a-input-number v-model:value="formState.multiple" addon-before=">=" addon-after="倍"></a-input-number>
      </a-form-item>
      <a-form-item label="钻石">
        <a-input-number v-model:value="formState.diamond" addon-before=">=" addon-after="个"></a-input-number>
      </a-form-item>
      <a-form-item label="狂暴">
        <a-input-number v-model:value="formState.crazy" addon-before=">=" addon-after="个"></a-input-number>
      </a-form-item>
      <a-form-item label="冰冻">
        <a-input-number v-model:value="formState.cold" addon-before=">=" addon-after="个"></a-input-number>
      </a-form-item>
      <a-form-item label="精准">
        <a-input-number v-model:value="formState.precise" addon-before=">=" addon-after="个"></a-input-number>
      </a-form-item>
      <a-form-item :wrapper-col="{ offset: 6, span: 18 }">
        <a-button type="primary" :disabled="formState.gold > 0 ? false : true " html-type="submit" @click="getSearchData">搜索</a-button>
      </a-form-item>
    </a-form>
  </a-drawer>
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
              <a-button type="primary" :disabled="open ? true : false " @click="showDrawer">筛选</a-button>
            </a-space>
          </a-col>
          <a-col :span="12" :style="{'text-align': 'right'}">
            <a-space v-if="state.selectedRowKeys.length > 0">
              总价：{{state.total}} 
              <a-button type="primary" :disabled="state.selectedRowKeys.length > 0 ? false : true " @click="pushCart">加入购物车</a-button>
              <a-button type="primary" danger :disabled="state.selectedRowKeys.length > 0 ? false : true " @click="postCart">直接购买</a-button>
            </a-space>
          </a-col>
          <a-col :span="24" v-if="state.search" :style="{'margin-top': '1rem'}">
            <a-tag color="processing" closable @close="cleanSearchStatus" v-if="state.search">
              筛选条件： {{getGameName()}}
              {{formState.gold > 0 ? ` & 金币 >= ${formState.gold}亿` : ""}}
              {{formState.multiple > 0 ? ` & 炮台 >= ${formState.diamond}倍` : ""}}
              {{formState.diamond > 0 ? ` & 钻石 >= ${formState.diamond}个` : ""}}
              {{formState.cold > 0 ? ` & 冰冻 >= ${formState.cold}个` : ""}}
              {{formState.crazy > 0 ? ` & 狂暴 >= ${formState.crazy}个` : ""}}
              {{formState.precise > 0 ? ` & 精准 >= ${formState.precise}个` : ""}}
            </a-tag>
          </a-col>
        </a-row>
        <a-empty :image="simpleImage" v-if="dataState.data.length == 0" />
        <a-table
          v-else
          sticky
          :row-selection="{ selectedRowKeys: state.selectedRowKeys, onChange: onSelectChange }"
          :style="{'margin-top': '1rem'}"
          :columns="columns"
          :rowKey="(record: ProductDatas) => record.ID"
          :data-source="dataState.data"
          :loading="state.loading"
          size="small"
          :scroll="{ x: 1500 }"
          :pagination="false">
          <template v-slot:bodyCell="{column,record}">
            <template v-if="column.dataIndex==='Cover'">
              <a-popover placement="bottomRight" arrow-point-at-center v-if="record.Cover.length != 0">
                <template #content>
                  <img :src="record.Cover" />
                </template>
                <img :src="record.Cover" :style="{'width': '55px'}" />
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
import { onMounted, computed, ref, h, reactive, toRaw } from 'vue'
import type { UnwrapRef } from 'vue'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { Empty } from 'ant-design-vue';
import { GetProducts, GetGamesList, AddCart, PostOrders, SearchProducts } from '../../../wailsjs/go/main/App'

import { useRouter } from 'vue-router'
const router = useRouter()

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

interface FormState {
  gameid: number;
  gold: number;
  multiple: number;
  diamond: number;
  crazy: number;
  cold: number;
  precise: number;
}
const formState: UnwrapRef<FormState> = reactive({
  gameid: 0,
  gold: 0,
  multiple: 0,
  diamond: 0,
  crazy: 0,
  cold: 0,
  precise: 0,
})
const labelCol = { span: 6 }
const wrapperCol = { span: 18 }

const columns = [
  {
    title: '帐号',
    dataIndex: 'Account',
    fixed: 'left',
    width: 150
  },
  {
    title: '单价',
    dataIndex: 'Price',
    customRender: function (t: any) {
      return `￥${t.value}`
    },
    fixed: 'left',
  },
  {
    title: '图片',
    dataIndex: 'Cover',
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
    title: '备注',
    dataIndex: 'Remarks',
  }
]

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
  state.value.search = false
  getProducts("1","20",String(id))
}

const getProducts = async(page:string = "1", pageSize:string = "20", gameid:string = String(gameID.value)) => {
  const search = state.value.search
  const games = await GetGamesList()
  
  if (games.status == 0) {
    // userState.value = data as User
    // console.log(data)
    gameState.value = games as GamesResponse
    if (games.data.length > 0) {
      formState.gameid = games.data[0].ID
    }
  }
  state.value.loading = true
  let params = {}
  if (search) {
    console.log(gameID.value)
    params = {
      page: page,
      limit: pageSize,
      gameid: String(gameID.value),
      gold: String(formState.gold * 100000000),
      multiple: String(formState.multiple),
      diamond: String(formState.diamond),
      crazy: String(formState.crazy),
      cold: String(formState.cold),
      precise: String(formState.precise),
    }
    const data = await SearchProducts(params)
    if (data.status == 0) {
      // userState.value = data as User
      // console.log(data)
      dataState.value = data as ProductResponse
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
  }else {
    params = {
      page: page,
      limit: pageSize,
      gameid: gameid
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
      if (data.message == "403") {
        router.push({
          'name': 'login',
        })
      }
    }
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
  selectedRowKeys: Number[];
  loading: boolean;
  total: number;
  search: boolean;
}
const state = ref<State>({
  selectedRowKeys: [],
  loading: false,
  total: 0,
  search: false,
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

const onSelectChange = (selectedRowKeys: Number[]) => {
  // console.log('selectedRowKeys changed: ', selectedRowKeys);
  state.value.selectedRowKeys = selectedRowKeys
  makeTotal()
}

const open = ref<boolean>(false)

const showDrawer = () => {
  open.value = true
}

const getSearchData = () => {
  state.value.search = true
  getProducts()
  setTimeout(()=>{open.value = false}, 500)
}

const handleChange = (value: number) => {
  gameID.value = value
}

const getGameName = () => {
  let gameName = ""
  gameState.value.data.forEach( (key: GamesDatas) => {if (key.ID == gameID.value) gameName = key.GameName})
  return gameName
}

const cleanSearchStatus = () => {
  state.value.search = false
  open.value = false
  gameID.value = 0
  formState.gameid = 0
  formState.gold = 0
  formState.multiple = 0
  formState.diamond = 0
  formState.crazy = 0
  formState.cold = 0
  formState.precise = 0
  getProducts()
}
</script>
