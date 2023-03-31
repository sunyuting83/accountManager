<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            用户列表
          </p>
          <button class="card-header-icon">
            <button class="button is-link is-small" @click="showAddModel">
              <span class="icon is-small">
                <i class="fa fa-plus"></i>
              </span>
              <span>添加用户</span>
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
                <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left	">
                  <thead>
                    <tr>
                      <td width="35%">用户名</td>
                      <td>状态</td>
                      <td v-if="userid == '1'">所属管理员</td>
                      <td>创建时间</td>
                      <td width="30%">操作</td>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item) in data" :key="item.ID">
                      <td>{{item.UserName}}</td>
                      <td>
                        <span class="has-text-success" v-if="item.NewStatus === 0">正常</span>
                        <span class="has-text-danger" v-else>锁定</span>
                      </td>
                      <td  v-if="userid == '1'">{{item.Manager.UserName}}</td>
                      <td><FormaTime :DateTime="item.CreatedAt"></FormaTime></td>
                      <td>
                        <div class="buttons">
                          <button class="button is-success is-small" @click="()=>{showModel(item.UserName)}">修改密码</button>
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
      </div>

      <PaginAtion v-if="data.length >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetData"></PaginAtion>
    </div>
    <ChangePassword
      :showData="openModal"
      :ShowMessage="ShowMessage"
      :Admin="false">
    </ChangePassword>
    <AddUser
      :showData="openAddModal"
      :ShowMessage="ShowMessage">
    </AddUser>
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
import ChangePassword from "@/components/Other/ChangePassword"
import NotIfication from "@/components/Other/Notification"
import AddUser from "@/components/User/AddUser"
import PopoButton from '@/components/Other/PopoButton'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'UserList',
  components: { ManageHeader, LoadIng, EmptyEd, ChangePassword, NotIfication, AddUser, PopoButton, PaginAtion, FormaTime },
  setup() {
    let states = reactive({
      loading: false,
      data: [],
      total: 0,
      username: "",
      userid:"",
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
      document.title = `${Config.GlobalTitle}-用户管理`
      const data = await CheckLogin()
      if (data == 0) {
        const username = localStorage.getItem('user')
        const userid = localStorage.getItem('userid')
        states.username = username
        states.userid = userid
        GetData()
      }else{
        setStorage(false)
        router.push("/")
      }
    })
    const GetData = async(page = 1) => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.userList, {page:page, limit: states.limit}, 'GET', token)
      states.loading = true
      states.pageLoading = true
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
        default:
          break;
      }
    }
    const showModel = (username) => {
      states.openModal.active = true
      states.openModal.username = username
    }
    const showAddModel = () => {
      states.openAddModal.active = true
    }
    const lockIt = async(id) => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.upuserstatus, {id: id}, 'PUT', token)
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
      const d = await Fetch(Config.Api.deluser, {id: id}, 'DELETE', token)
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

    return {
      ...toRefs(states),
      ShowMessage,
      showModel,
      showAddModel,
      lockIt,
      deleteIt,
      GetData
    }
  },
})
</script>
