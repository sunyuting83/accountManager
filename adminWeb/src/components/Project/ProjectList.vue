<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            项目列表
          </p>
          <button class="card-header-icon">
            <button class="button is-link is-small" @click="showAddModel">
              <span class="icon is-small">
                <i class="fa fa-plus"></i>
              </span>
              <span>添加项目</span>
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
              <div class="table-container" v-else>
                <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                  <thead>
                    <tr>
                      <td width="8%">项目名称</td>
                      <td>所属用户</td>
                      <td>Key</td>
                      <td>状态</td>
                      <td>可乐API</td>
                      <td>创建时间</td>
                      <td>用户名</td>
                      <td>密码</td>
                      <td>帐号数</td>
                      <td width="25%">操作</td>
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
                          <button class="button is-success is-small" @click="()=>{showAccount(item.Key)}">帐号管理</button>
                          <button class="button is-info is-small" @click="()=>{showAccountDraw(item.Key)}">提号管理</button>
                          <button class="button is-primary is-small" @click="()=>{showAccountDrawed(item.Key)}">已提取帐号</button>
                          <button class="button is-warning is-small" @click="()=>{showModifyModal(item.ID)}">修改项目</button>
                          <PopoButton
                            :message="item.NewStatus === 0?'锁定':'解锁'" color="is-dark" :callBack="()=>{lockIt(item.ID)}" v-if="item.UserName !== username"></PopoButton>
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
      </div>
      <PaginAtion v-if="data.length >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetData"></PaginAtion>
    </div>
    <AddProject
      v-if="userLoading"
      :showData="openAddModal"
      :ShowMessage="ShowMessage"
      :UserData="UserData">
    </AddProject>
    <ModifyProject
      v-if="modifyStatus"
      :showData="openModifyModal"
      :ShowMessage="ShowMessage">
    </ModifyProject>
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
import AddProject from "@/components/Project/AddProject"
import ModifyProject from "@/components/Project/ModifyProject"
import PopoButton from '@/components/Other/PopoButton'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'ProjectList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, AddProject, PopoButton, PaginAtion, FormaTime, ModifyProject },
  setup() {
    let states = reactive({
      loading: false,
      data: [],
      UserData: [],
      total: 0,
      username: "",
      userLoading: false,
      modifyStatus: false,
      openModal:{
        active: false,
        username: ""
      },
      openAddModal:{
        active: false
      },
      openModifyModal:{
        active: false,
        Project: {}
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
      document.title = `${Config.GlobalTitle}-项目管理`
      const data = await CheckLogin()
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
      const d = await Fetch(Config.Api.projectList, {page:page, limit: states.limit}, 'GET', token)
      states.loading = true
      states.pageLoading = false
      if (d.status == 0) {
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
        case 5:
          states.modifyStatus = false
          states.openModifyModal.Project = {}
          break;
        case 6:
          states.modifyStatus = false
          states.data = states.data.map((el)=>{
            console.log(el.ID, e.data.ID)
            if (el.ID == e.data.ID) {
              return el = e.data
            }
            return el
          })
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
    const showModifyModal = async(id) => {
      states.openModifyModal.Project = states.data.filter((e) => {
        return e.ID == id
      })[0]
      states.openModifyModal.active = true
      states.modifyStatus = true
      // console.log(states.openModifyModal.Project)
    }
    const getUserData = async() => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.UsersAllList, {}, 'GET', token)
      // console.log(d)
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

    const showAccount = (key) => {
      router.push(`/account/${key}`)
    }

    const showAccountDraw = (id) => {
      router.push(`/accountDraw/${id}/gold`)
    }
    const showAccountDrawed = (id) => {
      router.push(`/accountDrawed/${id}`)
    }

    return {
      ...toRefs(states),
      ShowMessage,
      showModel,
      showAddModel,
      lockIt,
      deleteIt,
      GetData,
      showAccount,
      showModifyModal,
      showAccountDraw,
      showAccountDrawed
    }
  },
})
</script>
