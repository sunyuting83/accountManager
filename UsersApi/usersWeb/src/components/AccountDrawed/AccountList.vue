<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container.is-fullhd">
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
              <div class="field mt-5" v-if="data.length > 0">
                <div class="columns  flex-wrap is-justify-content-space-between">
                  <div class="field is-horizontal ml-3 mr-2">
                    <label class="checkbox mr-4">
                      <input type="checkbox" v-model="form.multiple">
                      炮台
                    </label>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <label class="checkbox mr-4">
                      <input type="checkbox"  v-model="form.diamond">
                      钻石
                    </label>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <label class="checkbox mr-4">
                      <input type="checkbox"  v-model="form.crazy">
                      狂暴
                    </label>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <label class="checkbox mr-4">
                      <input type="checkbox"  v-model="form.cold">
                      冰冻
                    </label>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <label class="checkbox mr-4">
                      <input type="checkbox" v-model="form.precise">
                      瞄准
                    </label>
                  </div>
                  <div class="field is-horizontal ml-3 mr-2">
                    <label class="checkbox mr-4">
                      <input type="checkbox" v-model="form.remarks">
                      其他
                    </label>
                  </div>
                </div>
              </div>
              <div class="columns flex-wrap is-justify-content-space-between mt-1">
                <div class="field mr-3">
                  <button class="button is-dark is-small" @click="showAccountFiled">查看归档</button>
                </div>
                <div class="field mr-3">
                  <div class="buttons is-horizontal are-small has-addons">
                    <span v-if="total !== 0" class="is-size-7 mr-3">帐号总数 <span class="has-text-danger ml-1">{{total}}</span></span>
                    <DownloadFile 
                      :uri="`${RootUrl}${AccountKey}/ExportDrawed`"
                      styles="is-info"
                      :Data="form"
                      :buttonLoading="buttonLoading"
                      title="导出当前日期数据到文本"
                      ext=".txt"
                      v-if="data.length > 0" />
                    <DownloadFile 
                      :uri="`${RootUrl}${AccountKey}/ExportDrawed`"
                      styles="is-success"
                      :Data="form"
                      :buttonLoading="buttonLoading"
                      title="导出当前日期数据到Excel"
                      ext=".xlsx"
                      v-if="data.length > 0" />
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
              <div class="table-container" v-else>
                <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                  <thead class="is-size-7">
                    <tr>
                      <td>序号</td>
                      <td>帐号</td>
                      <!-- <td v-if="data[0].PhoneNumber.length > 0">手机号</td> -->
                      <td>今日金币</td>
                      <td>昨日金币</td>
                      <td>炮台</td>
                      <td>钻石</td>
                      <td>狂暴</td>
                      <td>冰冻</td>
                      <td>瞄准</td>
                      <td>价格</td>
                      <td>过期时间</td>
                      <td>更新时间</td>
                    </tr>
                  </thead>
                  <tbody class="is-size-7">
                    <tr v-for="(item, index) in data" :key="item.ID" class="hasimg">
                      <td>{{index + 1}}</td>
                      <td>{{item.UserName}}</td>
                      <!-- <td v-if="item.PhoneNumber.length > 0">{{item.PhoneNumber}}</td> -->
                      <td><FormaNumber :Numbers="item.TodayGold" /></td>
                      <td><FormaNumber :Numbers="item.YesterdayGold" /></td>
                      <td><FormaNumber :Numbers="item.Multiple" /></td>
                      <td>{{item.Diamond}}</td>
                      <td>{{item.Crazy}}</td>
                      <td>{{item.Cold}}</td>
                      <td>{{item.Precise}}</td>
                      <td>{{item.Price}}</td>
                      <td><ExpTime :DateTime="item.Exptime" /></td>
                      <td class="potd">
                        <FormaTime :DateTime="item.UpdatedAt" />
                        <div v-if="item.Cover.length > 0" class="poimg">
                          <img :data-src="IMGUri+item.Cover" :src="IMGUri+item.Cover" />
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
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
import DownloadFile from '@/components/Other/DownloadFile.vue'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AccountList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, PaginAtion, FormaTime, FormaNumber, ExpTime, DownloadFile },
  setup() {
    let states = reactive({
      RootUrl: Config.RootUrl,
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
      form: {
        date: "",
        multiple: true,
        diamond: false,
        crazy: false,
        cold: false,
        precise: false,
        remarks: false,
        excel: false,
      },
      IMGUri: Config.IMGUri,
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
          states.form.date = dlist[0]
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
      states.form.date = date
      GetDateList()
    }

    const backRouter = () => {
      router.push("/project")
    }

    const showAccountFiled = () => {
      router.push(`/accountFiled/${states.AccountKey}`)
    }

    return {
      ...toRefs(states),
      backRouter,
      getDateData,
      GetDateList,
      showAccountFiled
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

.hasimg .potd .poimg {
  position: absolute;
  right: 0;
  min-width: 570px;
  min-height: 76px;
  display: none;
  z-index: 10000;
}
.hasimg:hover .potd .poimg {
  display: block;
}
</style>