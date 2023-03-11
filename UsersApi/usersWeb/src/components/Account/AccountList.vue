<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            项目名称：{{projects.ProjectsName}}
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
              <div v-if="statusList.length > 0">
                <div class="buttons are-small has-addons">
                  <span class="f-1" v-for="(item,index) in statusList" :key="item.status" >
                    <button class="button" v-if="item.status !== '108'" :class="item.status === CurrentStatus.status?'is-success':''" @click="()=>{pushToData(index)}">
                      {{item.title}}
                    </button>
                  </span>
                </div>
                <div>
                  <table>
                    <tr>
                      <td align="right">
                        <div class="buttons are-small has-addons is-justify-content-flex-end">
                          <button class="button is-small is-success is-light" v-if="CurrentStatus.import">
                            导入{{CurrentStatus.title}}帐号
                          </button>
                          <button class="button is-small is-link is-light" v-if="CurrentStatus.callback">
                            退回{{CurrentStatus.title}}帐号
                          </button>
                          <button class="button is-small is-danger" v-if="CurrentStatus.delete">
                            删除{{CurrentStatus.title}}帐号
                          </button>
                          <button class="button is-small is-warning is-light" v-if="CurrentStatus.export">
                            导出{{CurrentStatus.title}}帐号
                          </button>
                        </div>
                      </td>
                    </tr>
                  </table>
                </div>
              </div>
              <div v-if="data.length <= 0">
                <EmptyEd></EmptyEd>
              </div>
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left" v-else>
                <thead>
                  <tr>
                    <td width="10%">项目名称</td>
                    <td>所属用户</td>
                    <td>Key</td>
                    <td>状态</td>
                    <td>可乐API</td>
                    <td>创建时间</td>
                    <td>用户名</td>
                    <td>密码</td>
                    <td>帐号数</td>
                    <td width="18%">操作</td>
                  </tr>
                </thead>
                <tbody class=" is-size-7">
                  <tr v-for="(item) in data" :key="item.ID">
                    <td>{{item.ProjectsName}}</td>
                    <td>{{item.Users.UserName}}</td>
                    <td>{{item.Key}}</td>
                    <td>
                      <span class="has-text-success" v-if="item.NewStatus === 0">正常</span>
                      <span class="has-text-danger" v-else>锁定</span>
                    </td>
                    <td>
                      <span class="has-text-success" v-if="item.ColaAPI">是</span>
                      <span class="has-text-danger" v-else>否</span>
                    </td>
                    <td><FormaTime :DateTime="item.CreatedAt"></FormaTime></td>
                    <td>{{item.UserName}}</td>
                    <td>{{item.Password}}</td>
                    <td>{{item.AccNumber}}</td>
                    <td>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
      <PaginAtion v-if="data.length >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetData"></PaginAtion>
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


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AccountList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, PaginAtion, FormaTime },
  setup() {
    let states = reactive({
      AccountKey: "",
      CurrentStatus: {},
      statusList: [],
      projects: {},
      loading: false,
      data: [],
      UserData: [],
      total: 0,
      username: "",
      userLoading: false,
      openModal:{
        active: false,
        username: ""
      },
      openAddModal:{
        active: false
      },
      openerr: {
        active: false,
        message: "",
        color: ""
      },
      pageLoading: false,
      limit: Config.Limit
    })
    const router = useRouter()
    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-帐号管理`
      const data = await CheckLogin()
      states.AccountKey = router.currentRoute._value.params.key
      if (data == 0) {
        const username = localStorage.getItem('user')
        states.username = username
        GetData(true)
      }else{
        setStorage(false)
        router.push("/")
      }
    })
    const GetData = async(first = false, page = 1) => {
      const token = localStorage.getItem("token")
      const data = {
        page:page, 
        limit: states.limit,
        status: states.status,
      }
      const url = Config.MakeAccountListUri(states.AccountKey)
      const d = await Fetch(url, data, 'GET', token)
      states.loading = true
      states.pageLoading = false
      if (d.status == 0) {
        states.data = d.data
        states.total = d.total
        states.projects = d.projects
        states.statusList = JSON.parse(d.projects.StatusJSON)
        states.pageLoading = true
        states.loading = false
        if (first) states.CurrentStatus = states.statusList[0]
      }else{
        states.data = []
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
    const ShowMessage = (e, status = 0, id=0) => {
      states.openerr.active = e.active
      states.openerr.message = e.message
      states.openerr.color = e.color
      switch (status) {
        case 1:
          states.data = [...states.data, e.data]
          states.userLoading = false
          break;
        case 2:
          states.data = states.data.map((e)=>{
            if(e.ID == id) {
              if(e.NewStatus === 0) {
                e.NewStatus = 1
              }else {
                e.NewStatus = 0
              }
            }
            return e
          })
          break;
        case 3:
          states.data = states.data.filter((e)=>{
            return e.ID !== id
          })
          break;
        case 4:
          states.userLoading = false
          states.UserData = []
          break;
        default:
          break;
      }
    }
    const showModel = (username) => {
      states.openModal.active = true
      states.openModal.username = username
    }
    

    const backRouter = () => {
      router.back()
    }

    const pushToData = (index) => {
      states.CurrentStatus = states.statusList[index]
      GetData()
    }

    return {
      ...toRefs(states),
      ShowMessage,
      showModel,
      GetData,
      backRouter,
      pushToData
    }
  },
})
</script>
<style scoped>
.f-1 {
  margin-left: -1px;
}
</style>