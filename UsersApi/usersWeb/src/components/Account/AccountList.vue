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
                  
                <div class="buttons are-small has-addons is-justify-content-flex-end mb-3">
                  <span v-if="total !== 0" class="is-size-7 mr-3">帐号总数 <span class="has-text-danger ml-1">{{total}}</span></span>
                  <button class="button is-small is-success is-light" :class="buttonLoading?'is-loading':''" v-if="CurrentStatus.import" @click="showPostModal">
                    导入{{CurrentStatus.title}}帐号
                  </button>
                  <button class="button is-small is-link is-light" :class="buttonLoading?'is-loading':''"  v-if="CurrentStatus.callback && data.length > 0" @click="backTo">
                    退回{{CurrentStatus.title}}帐号
                  </button>
                  <PopoButton :message="`删除${CurrentStatus.title}帐号`" color="is-danger"  :loading="buttonLoading" :callBack="deleteAccount" v-if="CurrentStatus.delete && data.length > 0"></PopoButton>
                  <button class="button is-small is-warning is-light" :class="buttonLoading?'is-loading':''" @click="ExportAccount" v-if="CurrentStatus.export && data.length > 0">
                    导出{{CurrentStatus.title}}帐号
                  </button>
                </div>
                </div>
              </div>
              <div v-if="data.length <= 0">
                <EmptyEd></EmptyEd>
              </div>
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left" v-else>
                <thead class="is-size-7">
                  <tr>
                    <td>序号</td>
                    <td>帐号</td>
                    <td v-if="data[0].Password.length > 0">密码</td>
                    <td v-if="data[0].PhoneNumber.length > 0">手机号</td>
                    <td v-if="data[0].PhonePassword.length > 0">手机密码</td>
                    <td>今日金币</td>
                    <td>昨日金币</td>
                    <td>炮台</td>
                    <td>钻石</td>
                    <td>狂暴</td>
                    <td>冰冻</td>
                    <td>瞄准</td>
                    <td v-if="data[0].Remarks.length > 0">其他</td>
                    <td v-if="data[0].Price.length > 0">价格</td>
                    <td>过期时间</td>
                    <td>创建时间</td>
                    <td>更新时间</td>
                  </tr>
                </thead>
                <tbody class="is-size-7">
                  <tr v-for="(item, index) in data" :key="item.ID">
                    <td>{{index}}</td>
                    <td>{{item.UserName}}</td>
                    <td v-if="item.Password.length > 0">{{item.Password}}</td>
                    <td v-if="item.PhoneNumber.length > 0">{{item.PhoneNumber}}</td>
                    <td v-if="item.PhonePassword.length > 0">{{item.PhonePassword}}</td>
                    <td><FormaNumber :Numbers="item.TodayGold" /></td>
                    <td><FormaNumber :Numbers="item.YesterdayGold" /></td>
                    <td><FormaNumber :Numbers="item.Multiple" /></td>
                    <td>{{item.Diamond}}</td>
                    <td>{{item.Crazy}}</td>
                    <td>{{item.Cold}}</td>
                    <td>{{item.Precise}}</td>
                    <td v-if="item.Remarks.length > 0" class="w165">{{item.Remarks}}</td>
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
      <PaginAtion v-if="total >= limit && pageLoading === true" :total="total" :number="limit" :GetData="GetData"></PaginAtion>
    </div>
    <PostData
      v-if="postStatus"
      :showData="openPostModal"
      :ShowMessage="ShowMessage">
    </PostData>
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
import PostData from "@/components/Account/Postdata"
import NotIfication from "@/components/Other/Notification"
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'
import FormaNumber from '@/components/Other/FormaNumber'
import PopoButton from '@/components/Other/PopoButton'
import ExpTime from '@/components/Other/ExpTime.vue'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AccountList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, PaginAtion, FormaTime, PostData, PopoButton, FormaNumber, ExpTime },
  setup() {
    let states = reactive({
      AccountKey: "",
      CurrentStatus: {},
      statusList: [],
      projects: {},
      loading: false,
      data: [],
      total: 0,
      username: "",
      buttonLoading: false,
      postStatus: false,
      openPostModal:{
        active: false,
        postParams: {}
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
        GetData(1,true)
      }else{
        setStorage(false)
        router.push("/")
      }
    })

    const CleanData = () => {
      states.data = []
      states.total = 0
      states.page = []
      states.projects = {}
      states.pageLoading = true
      states.loading = false
      states.buttonLoading = false
    }

    const GetData = async(page = 1, first = false) => {
      const token = localStorage.getItem("token")
      let status = states.CurrentStatus.status
      if (first) status = "0"
      const data = {
        page:page, 
        limit: states.limit,
        status: status,
      }
      const url = `${Config.RootUrl}${states.AccountKey}/AccountList`
      states.loading = true
      const d = await Fetch(url, data, 'GET', token)
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
        states.pageLoading = false
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
        case 4:
          GetData()
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
      states.data = []
      states.total = 0
      states.pageLoading = false
      GetData()
    }

    const showPostModal = () => {
      states.openPostModal.active = true
      // console.log(states.CurrentStatus)
      states.openPostModal.postParams = states.CurrentStatus
      states.postStatus = true
    }

    const deleteAccount = async() => {
      const token = localStorage.getItem("token")
      let status = states.CurrentStatus.status
      const data = {
        status: status,
      }
      const url = `${Config.RootUrl}${states.AccountKey}/DeleteAccount`
      states.loading = true
      states.pageLoading = false
      states.buttonLoading = true
      const d = await Fetch(url, data, 'DELETE', token)
      if (d.status == 0) {
        CleanData()
      }else{
        CleanData()
      }
    }

    const backTo = async() => {
      const token = localStorage.getItem("token")
      let status = states.CurrentStatus.status
      const data = {
        status: status,
      }
      const url = `${Config.RootUrl}${states.AccountKey}/GoBackAccount`
      const d = await Fetch(url, data, 'PUT', token)
      states.loading = true
      states.pageLoading = false
      states.buttonLoading = true
      if (d.status == 0) {
        CleanData()
      }else{
        CleanData()
      }
    }

    const ExportAccount = async() => {
      const d = await exportFile()
      download(d)
    }

    const exportFile = () => {
      const token = localStorage.getItem("token")
      let status = states.CurrentStatus.status

      const url = `${Config.RootUrl}${states.AccountKey}/ExportAccount`
      let requestConfig = {
        method: "put",
        responseType: "blob"
      }
      Object.defineProperty(requestConfig, 'body', {
          value: JSON.stringify({
          status: status,
        })
      })
      requestConfig.headers = new Headers({
        Accept: '*/*',
      })
      requestConfig.headers.append("Content-Type","application/json;charset=UTF-8")
      requestConfig.headers.append('Authorization',`Bearer ${token}`)
      return new Promise((resolve) => {
        fetch(url, requestConfig)
          .then(res => {
            if(res.ok) {
              resolve(res.text())
            }else {
              resolve({
                status: 1,
                message: "访问出错"
              })
            }
          })
          .catch((err) => {
            resolve({
              status: 1,
              message: err.message
            })
          })
      })
    }
    const download = (data) => {
        if (!data) {
            return
        }
        // const contentType = data.type
        // const fileName = contentType.split('filename=')[1]
        let url = window.URL.createObjectURL(new Blob([data]))
        let link = document.createElement('a')
        link.style.display = 'none'
        link.href = url
        link.id='Adownload'
        const date = new Date(),
            Y = date.getFullYear(),
            M = date.getMonth(),
            D = date.getDate(),
            h = date.getHours(),
            m = date.getMinutes(),
            s = date.getSeconds(),
            fileName = `${String(Y)}${String(M)}${String(D)}${String(h)}${String(m)}${String(s)}.txt`
        // console.log(fileName)
        link.setAttribute('download', fileName)
        
        document.body.appendChild(link)
        link.click()
        document.getElementById('Adownload').remove();
    }

    return {
      ...toRefs(states),
      ShowMessage,
      showModel,
      GetData,
      backRouter,
      pushToData,
      showPostModal,
      deleteAccount,
      backTo,
      ExportAccount
    }
  },
})
</script>
<style scoped>
.f-1 {
  margin-left: -1px;
}
.w165 {
  width: 165px;
}
</style>