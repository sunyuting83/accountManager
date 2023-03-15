<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            项目名称：{{projects.ProjectsName}} --- 已提取帐号管理
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
                <div class="field mr-3"></div>
                <div class="field mr-3">
                  <div class="buttons is-horizontal are-small has-addons">
                    <span v-if="total !== 0" class="is-size-7 mr-3">帐号总数 <span class="has-text-danger ml-1">{{total}}</span></span>
                    <button 
                      class="button is-info"
                      :class="buttonLoading ? 'is-loading' : '' "
                      @click="pullData">
                      导出当前日期数据到文本
                    </button>
                    <button 
                      class="button is-success"
                      :class="buttonLoading ? 'is-loading' : '' "
                      @click="pullData">
                      导出当前日期数据到Excel
                    </button>
                  </div>
                </div>
              </div>
              <div class="buttons are-small has-addons" v-if="dateList.length > 0">
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
      <PaginAtion v-if="total >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetDateList"></PaginAtion>
    </div>
    <NotIfication
      :showData="openerr">
    </NotIfication>
  </div>
</template>
<script>
import { reactive, toRefs, onMounted, defineComponent } from 'vue'
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
      projects: {},
      CurrentDate: "",
      dateList: [],
      loading: false,
      data: [],
      total: 0,
      username: '',
      buttonLoading: false,
      openerr: {
        active: false,
        message: "",
        color: ""
      },
      pageLoading: false,
      limit: Config.Limit,
    })
    const router = useRouter()
    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-帐号管理`
      
      const data = await CheckLogin()
      if (data == 0) {
        states.AccountKey = router.currentRoute._value.params.key
        const username = localStorage.getItem('user')
        states.username = username
        
        const dlist = await GetDate()
        if (states.dateList.length > 0) {
          states.CurrentDate = dlist[0]
          states.loading = false
          GetDateList()
        }else {
          states.loading = false
        }
      }else{
        setStorage(false)
        router.push("/")
      }
    })

    // const CleanData = () => {
    //   states.data = []
    //   states.temp = []
    //   states.checkTemp = []
    //   states.total = 0
    //   states.page = []
    //   states.projects = {}
    //   states.CurrentDate= ""
    //   states.dateList= []
    //   states.pageLoading = true
    //   states.loading = false
    //   states.buttonLoading = false
    // }

    const GetDate = async() => {
      const token = localStorage.getItem("token")
      const data = {}
      const url = `${Config.RootUrl}${states.AccountKey}/GetAllDateForDrawed`
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
      const url = `${Config.RootUrl}${states.AccountKey}/AccountDrawedDateList`
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

    return {
      ...toRefs(states),
      backRouter,
      getDateData,
      GetDateList,
      ShowMessage
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