<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container.is-fullhd">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            项目名称：{{projects.ProjectsName}} --- 提取记录管理
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
              <div class="buttons are-small has-addons mt-1" v-if="dateList.length > 0">
                <button
                  class="button is-info is-light"
                  v-for="(item,index) in dateList"
                  :key="index"
                  :disabled="CurrentDate == item.FiledName ? true : false"
                  @click="()=>{getDateData(item.FiledName)}">
                  {{item.FiledName}}
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
                      <td>日志时间</td>
                      <td>题号人</td>
                      <td>题号时间</td>
                      <td>查看详细</td>
                    </tr>
                  </thead>
                  <tbody class="is-size-7">
                    <tr v-for="(item, index) in data" :key="item.ID" class="hasimg">
                      <td>{{index + 1}}</td>
                      <td>{{item.LogName}}</td>
                      <td>{{item.DrawUser}}</td>
                      <td>
                        <FormaTime :DateTime="item.UpdatedAt" />
                      </td>
                      <td><router-link :to="`/drawData/${item.ID}`" class="button is-success is-small">查看详细</router-link></td>
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
  </div>
</template>
<script>
import { reactive, toRefs, onMounted, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'
import EmptyEd from '@/components/Other/Empty'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'DrawedLog',
  components: { ManageHeader, LoadIng, EmptyEd, PaginAtion, FormaTime },
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
      pageLoading: false,
      limit: Config.Limit,
      IMGUri: Config.IMGUri,
    })
    const router = useRouter()
    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-游戏管理`
      const data = await CheckLogin()
      if (data == 0) {
        states.AccountKey = router.currentRoute._value.params.key
        const username = localStorage.getItem('user')
        states.username = username
        GetData()
      }else{
        setStorage(false)
        router.push("/")
      }
    })

    const backRouter = () => {
      router.push("/project")
    }
    const GetData = async(page = 1) => {
      const token = localStorage.getItem("token")
      const url = `${Config.RootUrl}${states.AccountKey}/DrawList`
      const d = await Fetch(url, {page:page, limit: states.limit}, 'GET', token)
      states.loading = true
      states.pageLoading = false
      if (d.status == 0) {
        if (d.data.length != 0) states.projects = d.data[0].Projects
        states.data = d.data
        states.total = d.total
      states.pageLoading = true
        states.loading = false
      }else{
        states.data = []
        states.total = 0
        states.page = []
        states.loading = false
      }
    }
    return {
      ...toRefs(states),
      GetData,
      backRouter
    }
  },
})
</script>