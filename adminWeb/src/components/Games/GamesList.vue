<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container.is-fullhd">
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            游戏列表
          </p>
          <button class="card-header-icon">
            <button class="button is-link is-small" @click="showAddModel">
              <span class="icon is-small">
                <i class="fa fa-plus"></i>
              </span>
              <span>添加游戏</span>
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
                      <td width="35%">游戏名</td>
                      <td>创建时间</td>
                      <td width="30%">操作</td>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item) in data" :key="item.ID">
                      <td>{{item.GameName}}</td>
                      <td><FormaTime :DateTime="item.CreatedAt"></FormaTime></td>
                      <td>
                        <div class="buttons">
                          <button class="button is-success is-small" @click="()=>{showCaleModel(item.ID)}">价格公式</button>
                          <PopoButton message="删除" color="is-danger" :callBack="()=>{deleteIt(item.ID)}"></PopoButton>
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
    <AddGame
      :showData="openAddModal"
      :ShowMessage="ShowMessage">
    </AddGame>
    <ModifyCale
      v-if="modifyStatus"
      :showData="openCaleModal"
      :ShowMessage="ShowMessage">
    </ModifyCale>
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
import AddGame from "@/components/Games/AddGame"
import ModifyCale from "@/components/Games/ModifyCale"
import PopoButton from '@/components/Other/PopoButton'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'GamesList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, AddGame, PopoButton, PaginAtion, FormaTime, ModifyCale },
  setup() {
    let states = reactive({
      loading: false,
      data: [],
      total: 0,
      username: "",
      openModal:{
        active: false,
        username: ""
      },
      modifyStatus: false,
      openAddModal:{
        active: false
      },
      openCaleModal:{
        active: false,
        data: {
          GameID: 0,
          BasePrice:0,
          UnitPrice:0,
          SingleNumber: 0
        }
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
      document.title = `${Config.GlobalTitle}-游戏管理`
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
      const d = await Fetch(Config.Api.GamesList, {page:page, limit: states.limit}, 'GET', token)
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
          states.modifyStatus = false
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
          states.modifyStatus = false
          break;
        case 3:
          states.data = states.data.filter((e)=>{
            return e.ID !== id
          })
          states.modifyStatus = false
          break;
        default:
          states.modifyStatus = false
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
    const showCaleModel = (id) => {
      states.data.forEach((e)=>{
        if(e.ID == id) {
          states.openCaleModal.data.BasePrice = e.BasePrice
          states.openCaleModal.data.UnitPrice = e.UnitPrice
          states.openCaleModal.data.SingleNumber = e.SingleNumber
        }
      })
      states.openCaleModal.data.GameID = id
      states.openCaleModal.active = true
      states.modifyStatus = true
    }
    const deleteIt = async(id) => {
      const token = localStorage.getItem("token")
      const d = await Fetch(Config.Api.deladmin, {id: id}, 'DELETE', token)
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
      showCaleModel,
      deleteIt,
      GetData
    }
  },
})
</script>
