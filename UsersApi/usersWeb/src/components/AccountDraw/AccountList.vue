<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            项目名称：{{projects.ProjectsName}} --- 提号管理
          </p>
          <button class="card-header-icon">
            <button class="button is-link is-small" @click="backRouter">
              <span class="icon is-small">
                <i class="fa fa-arrow-circle-left"></i>
              </span>
              <span>返回</span>
            </button>
          </button>
        </header>
        <div class="card-content">
          <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
            <div class="com__box" v-if="loading" :style="loading? 'margin-top:5rem':''">
              <LoadIng></LoadIng>
            </div>
            <div v-else>
              <div class="columns flex-wrap is-justify-content-space-between mt-1">
                <div class="field ml-3">
                  <div class="buttons is-horizontal are-small has-addons">
                    <button class="button is-warning" :class="buttonLoading ? 'is-loading' : '' " :disabled="AccountType == 'gold'?true:false" @click="pushRouter">
                      按金币排列
                    </button>
                    <button class="button is-info" :class="buttonLoading ? 'is-loading' : '' " :disabled="AccountType == 'date'?true:false" @click="pushRouter">
                      按日期排列
                    </button>
                  </div>
                </div>
                <div class="field mr-3">
                  <div class="buttons is-horizontal are-small has-addons">
                    <span v-if="total !== 0" class="is-size-7 mr-3">帐号总数 <span class="has-text-danger ml-1">{{total}}</span></span>
                    <button 
                      class="button is-success"
                      :class="buttonLoading ? 'is-loading' : '' "
                      @click="pullData" :disabled="checkTemp.length <= 0">
                      提取选中
                    </button>
                    <button class="button is-info" :class="buttonLoading ? 'is-loading' : '' " @click="pushRouterToDrawed">
                      已提取帐号管理
                    </button>
                  </div>
                </div>
              </div>
              <div class="field mt-5">
                <div class="columns flex-wrap is-flex-wrap-wrap">
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              金币大于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="90" v-model="Filter.mingold">
                          </p>
                          <p class="control">
                            <span class="button is-small is-static">
                              亿
                            </span>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              金币小于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="500" v-model="Filter.maxgold">
                          </p>
                          <p class="control">
                            <span class="button is-small is-static">
                              亿
                            </span>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              炮台大于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="2500000" v-model="Filter.multiple">
                          </p>
                          <p class="control">
                            <span class="button is-small is-static">
                              万
                            </span>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              钻石大于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="300" v-model="Filter.diamond">
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              狂暴大于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="300" v-model="Filter.crazy">
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              冰冻大于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="300"  v-model="Filter.cold">
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-expanded">
                        <div class="field has-addons">
                          <p class="control">
                            <a class="button  is-small is-static">
                              瞄准大于
                            </a>
                          </p>
                          <p class="control is-expanded">
                            <input class="input is-small w165" min="0" type="number" placeholder="300" v-model="Filter.precise">
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="field ml-3">
                    <button
                      class="button is-primary is-small"
                      :class="buttonLoading ? 'is-loading' : '' "
                      :disabled="Filter.mingold > 0 ?false:true"
                      @click="pullSelectData">
                      按条件提取
                    </button>
                  </div>
                </div>
              </div>
              <div class="buttons are-small has-addons" v-if="AccountType == 'date'">
                <button
                  class="button is-info is-light"
                  v-for="(item,index) in dateList"
                  :key="index"
                  :disabled="CurrentDate == item ? true : false"
                  @click="()=>{getDateData(item)}">
                  {{item}}
                </button>
              </div>
              <div v-if="data.length <= 0">
                <EmptyEd></EmptyEd>
              </div>
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left" v-else>
                <thead class="is-size-7">
                  <tr>
                    <td>
                      <label class="checkbox"><input type="checkbox" @click="checkall" />全选</label>
                    </td>
                    <td>序号</td>
                    <td>帐号</td>
                    <td v-if="data[0].PhoneNumber.length > 0">手机号</td>
                    <td>今日金币</td>
                    <td>昨日金币</td>
                    <td>炮台</td>
                    <td>钻石</td>
                    <td>狂暴</td>
                    <td>冰冻</td>
                    <td>瞄准</td>
                    <td v-if="data[0].Price.length > 0">价格</td>
                    <td>过期时间</td>
                    <td>创建时间</td>
                    <td>更新时间</td>
                  </tr>
                </thead>
                <tbody class="is-size-7">
                  <tr v-for="(item, index) in data" :key="item.ID">
                    <td>
                      <label class="checkbox">
                        <input type="checkbox" v-model="item.check" @click="(e)=>checkBox(e,item.ID)">
                      </label>
                    </td>
                    <td>{{index + 1}}</td>
                    <td>{{item.UserName}}</td>
                    <td v-if="item.PhoneNumber.length > 0">{{item.PhoneNumber}}</td>
                    <td><FormaNumber :Numbers="item.TodayGold" /></td>
                    <td><FormaNumber :Numbers="item.YesterdayGold" /></td>
                    <td><FormaNumber :Numbers="item.Multiple" /></td>
                    <td>{{item.Diamond}}</td>
                    <td>{{item.Crazy}}</td>
                    <td>{{item.Cold}}</td>
                    <td>{{item.Precise}}</td>
                    <td v-if="item.Price.length > 0">{{item.Price}}</td>
                    <td><ExpTime :DateTime="item.Exptime" /></td>
                    <td><FormaTime :DateTime="item.CreatedAt" /></td>
                    <td><FormaTime :DateTime="item.UpdatedAt" /></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
      <PaginAtion v-if="AccountType == 'gold' && total >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetGoldData"></PaginAtion>
      <PaginAtion v-if="AccountType == 'date' && total >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetDateList"></PaginAtion>
    </div>
    <NotIfication
      :showData="openerr">
    </NotIfication>
  </div>
</template>
<script>
import { reactive, toRefs, onMounted, defineComponent, inject } from 'vue'
import { useRouter } from 'vue-router'
import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'
import EmptyEd from '@/components/Other/Empty'
import NotIfication from "@/components/Other/Notification"
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'
import FormaNumber from '@/components/Other/FormaNumber'
import ExpTime from '@/components/Other/ExpTime'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AccountList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, PaginAtion, FormaTime, FormaNumber, ExpTime },
  setup() {
    let states = reactive({
      AccountKey: "",
      AccountType: "",
      projects: {},
      CurrentDate: "",
      dateList: [],
      temp: [],
      checkTemp: [],
      loading: false,
      data: [],
      total: 0,
      username: "",
      buttonLoading: false,
      openerr: {
        active: false,
        message: "",
        color: ""
      },
      pageLoading: false,
      limit: Config.Limit,
      Filter: {
        mingold: 90,
        maxgold: 500,
        multiple: 250,
        diamond: 0,
        crazy: 0,
        cold: 0,
        precise: 0,
      }
    })
    const Reload = inject('reload')
    const router = useRouter()
    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-帐号管理`
      
      const data = await CheckLogin()
      if (data == 0) {
        states.AccountKey = router.currentRoute._value.params.key
        states.AccountType = router.currentRoute._value.params.type
        const username = localStorage.getItem('user')
        states.username = username
        if (states.AccountType == 'date') {
          const dlist = await GetDate()
          if (states.dateList.length > 0) {
            states.CurrentDate = dlist[0]
            states.loading = false
            GetDateList()
          }else {
            states.loading = false
          }
        }else {
          GetGoldData()
        }
      }else{
        setStorage(false)
        router.push("/")
      }
    })

    const CleanData = () => {
      states.data = []
      states.temp = []
      states.checkTemp = []
      states.total = 0
      states.page = []
      states.projects = {}
      states.CurrentDate= ""
      states.dateList= []
      states.pageLoading = true
      states.loading = false
      states.buttonLoading = false
    }

    const GetDate = async() => {
      const token = localStorage.getItem("token")
      const data = {}
      const url = `${Config.RootUrl}${states.AccountKey}/GetAllDateForDraw`
      states.loading = true
      const d = await Fetch(url, data, 'GET', token)
      if (d.status == 0) {
        states.dateList = d.dateList
        return d.dateList
      }else{
        states.data = []
        states.total = 0
        states.page = []
        states.projects = {}
        states.loading = false
      }
    }
    const GetDateList = async(page = 1) => {
      const token = localStorage.getItem("token")
      const data = {
        page:page, 
        limit: states.limit,
        date: states.CurrentDate
      }
      const url = `${Config.RootUrl}${states.AccountKey}/AccountDrawDateList`
      states.loading = true
      const d = await Fetch(url, data, 'GET', token)
      if (d.status == 0) {
        states.data = d.data
        states.temp = d.data
        states.total = d.total
        states.projects = d.projects
        states.pageLoading = true
        states.loading = false
      }else{
        states.data = []
        states.temp = []
        states.total = 0
        states.page = []
        states.projects = {}
        states.loading = false
      }
    }

    const getDateData = (date) => {
      states.CurrentDate = date
      GetDateList()
    }

    const GetGoldData = async(page = 1) => {
      const token = localStorage.getItem("token")
      const data = {
        page:page, 
        limit: states.limit,
      }
      const url = `${Config.RootUrl}${states.AccountKey}/AccountDrawList`
      states.loading = true
      const d = await Fetch(url, data, 'GET', token)
      if (d.status == 0) {
        states.data = d.data
        states.temp = d.data
        states.total = d.total
        states.projects = d.projects
        states.pageLoading = true
        states.loading = false
      }else{
        states.data = []
        states.temp = []
        states.total = 0
        states.page = []
        states.projects = {}
        states.loading = false
      }
    }
    /**
     * 
     * @param {*} e message用到的值
     * @param {*} status 0默认不传参 1添加-加入列表 2锁定-替换列表值 3删除-filter值
     */
    const ShowMessage = (e) => {
      states.openerr.active = e.active
      states.openerr.message = e.message
      states.openerr.color = e.color
    }
    

    const backRouter = () => {
      router.push("/project")
    }

    const pushRouter = () => {
      const accounType = states.AccountType
      const AccountKey = states.AccountKey
      if (accounType === 'gold') {
        states.AccountType = 'date'
      }else {
        states.AccountType = 'gold'
      }
      CleanData()
      router.push(`/accountDraw/${AccountKey}/${states.AccountType}`)
      Reload()
    }

    const checkBox = (e, account) => {
      if (e.target.checked) {
        states.checkTemp = [...states.checkTemp, account]
      }else{
        states.checkTemp = states.checkTemp.filter((el) => {
          if (el !== account) return el
        })
      }
    }
    const checkall = (e) => {
      if (e.target.checked) {
        states.checkTemp = []
        states.data.forEach((el) => {
          states.checkTemp =  [...states.checkTemp,el.ID]
        })
        states.data = states.data.map((el) => {
          el.check =  true
          return el
        })
      }else{
        states.checkTemp = []
        states.data = states.data.map((el) => {
          el.check =  false
          return el
        })
      }
    }

    const pullData = async() => {
      const list = states.checkTemp
      const e = {
        active: true,
        message: "获取失败",
        color: "is-danger",
        newtime: 0,
      }
      if (list.length > 0 ) {
        const token = localStorage.getItem("token")
        const data = {
          list: states.checkTemp
        }
        const url = `${Config.RootUrl}${states.AccountKey}/PullDrawList`
        states.buttonLoading = true
        states.pageLoading = false
        states.loading = true
        const d = await Fetch(url, data, 'PUT', token)
        if (d.status == 0) {
          states.loading = false
          states.buttonLoading = false
          e.color = 'is-success'
          e.message = d.message
          ShowMessage(e)
          const AccountType = router.currentRoute._value.params.type
          if (AccountType == 'date') {
            GetDateList()
          }else {
            GetGoldData()
          }
        }else{
          states.checkTemp = []
          states.loading = false
          states.buttonLoading = false
          ShowMessage(e)
        }
      }else{
        states.loading = false
        ShowMessage(e)
      }
    }

    const pullSelectData = async() => {
      const Filter = states.Filter
      const e = {
        active: true,
        message: "获取失败",
        color: "is-danger",
        newtime: 0,
      }
      const token = localStorage.getItem("token")
      const data = {
        mingold: Filter.mingold * 100000000,
        maxgold: Filter.maxgold * 100000000,
        multiple: Filter.multiple * 10000,
        diamond: Filter.diamond,
        crazy: Filter.crazy,
        cold: Filter.cold,
        precise: Filter.precise,
      }
      const url = `${Config.RootUrl}${states.AccountKey}/PullDrawSelect`
      states.buttonLoading = true
      states.pageLoading = false
      states.loading = true
      const d = await Fetch(url, data, 'PUT', token)
      if (d.status == 0) {
        states.loading = false
        states.buttonLoading = false
        e.color = 'is-success'
        e.message = d.message
        ShowMessage(e)
        const AccountType = router.currentRoute._value.params.type
        if (AccountType == 'date') {
          GetDateList()
        }else {
          GetGoldData()
        }
      }else{
        states.checkTemp = []
        states.loading = false
        states.buttonLoading = false
        ShowMessage(e)
      }
    }

    const pushRouterToDrawed = () => {
      const AccountKey = router.currentRoute._value.params.key
      router.push(`/accountDrawed/${AccountKey}`)
    }

    return {
      ...toRefs(states),
      GetGoldData,
      backRouter,
      pushRouter,
      getDateData,
      GetDateList,
      checkBox,
      checkall,
      pullData,
      pullSelectData,
      pushRouterToDrawed
    }
  },
})
</script>
<style scoped>
.f-1 {
  margin-left: -1px;
}
.w165 {
  min-width: 100px;
  max-width: 140px;
}
</style>