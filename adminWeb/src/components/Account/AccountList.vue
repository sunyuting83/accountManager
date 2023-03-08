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
                      <div class="buttons">
                        <button class="button is-success is-small" @click="()=>{showAccount(item.ID)}">帐号管理</button>
                        <PopoButton
                          :message="item.NewStatus === 0?'锁定':'解锁'" color="is-info" :callBack="()=>{lockIt(item.ID)}" v-if="item.UserName !== username"></PopoButton>
                        <PopoButton message="删除" color="is-danger" :callBack="()=>{deleteIt(item.ID)}" v-if="item.UserName !== username"></PopoButton>
                      </div>
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
import PopoButton from '@/components/Other/PopoButton'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AccountList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, PopoButton, PaginAtion, FormaTime },
  setup() {
    let states = reactive({
      id: 0,
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
      states.id = router.currentRoute._value.params.id
      if (data == 0) {
        const username = localStorage.getItem('user')
        states.username = username
        GetData()
      }else{
        setStorage(false)
        router.push("/")
      }
    })
    const GetData = async(page = 1) => {
      const token = localStorage.getItem("token")
      const data = {
        page:page, 
        limit: states.limit,
        projectsID: states.id
      }
      const d = await Fetch(Config.Api.accountList, data, 'GET', token)
      states.loading = true
      states.pageLoading = false
      if (d.status == 0) {
        states.data = d.data
        states.total = d.total
        states.projects = d.projects
        states.pageLoading = true
        states.loading = false
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
    const showAddModel = async() => {
      states.UserData = await getUserData()
      if (states.userLoading) states.openAddModal.active = true
    }
    const getUserData = async() => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.UsersAllList, {}, 'GET', token)
      console.log(d)
      if (d.status == 0) {
        states.UserData = d.data
        states.userLoading = true
      }else{
        states.UserData = []
        states.userLoading = false
      }
    }
    const lockIt = async(id) => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.upprojectstatus, {id: id}, 'PUT', token)
      if (d.status == 0) {
        const data = {
          active: true,
          message: d.message,
          color: 'is-success'
        }
        ShowMessage(data, 2, id)
      }else{
        const data = {
          active: true,
          message: d.message,
          color: 'is-danger'
        }
        ShowMessage(data, 0)
      }
    }
    const deleteIt = async(id) => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.delproject, {id: id}, 'DELETE', token)
      if (d.status == 0) {
        const data = {
          active: true,
          message: d.message,
          color: 'is-success'
        }
        ShowMessage(data, 3, d.id)
      }else{
        const data = {
          active: true,
          message: d.message,
          color: 'is-danger'
        }
        ShowMessage(data, 0)
      }
    }

    const backRouter = () => {
      router.back()
    }


    return {
      ...toRefs(states),
      ShowMessage,
      showModel,
      showAddModel,
      lockIt,
      deleteIt,
      GetData,
      backRouter
    }
  },
})
</script>
