<template>
  <div>
    <ManageHeader></ManageHeader>
    <nav class="level is-mobile pt-3" v-if="GamesList.length !== 0">
      <div class="level-item has-text-centered" v-for="item in GamesList" :key="item.ID" :class="item.Count != 0 ? '' : 'hiddenit'">
        <div>
          <p class="heading">{{item.GameName}}</p>
          <p class="title">{{item.Count}}</p>
        </div>
      </div>
    </nav>
    <div class="container.is-fullhd">
      <div class="card events-card">
        <div class="card-content">
          <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
            <div class="com__box" v-if="loading" :style="loading? 'margin-top:5rem':''">
              <LoadIng></LoadIng>
            </div>
            <div v-else>
              <div class="field mt-5">
                <div class="columns flex-wrap is-flex-wrap-wrap">
                  <div class="field is-horizontal ml-3 mr-2">
                    <div class="field-body">
                      <div class="field is-narrow">
                        <div class="control">
                          <div class="select is-small">
                            <select v-model="Filter.GameID">
                              <option :value=0>选择游戏</option>
                              <option v-for="(item) in GamesList" :key="item.ID" :value="item.ID">{{item.GameName}}</option>
                            </select>
                          </div>
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
                    <div class="buttons">
                      <PopoButton
                        message="按条件搜索" 
                        color="is-link" 
                        :callBack="GetDrawData" 
                        v-if="Filter.GameID !== 0"
                        ></PopoButton>
                      <PopoButton
                        message="按条件提取" 
                        color="is-primary" 
                        :callBack="pullData"
                        v-if="checkTemp.length !== 0"
                        ></PopoButton>
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="data.length === 0">
                <EmptyEd></EmptyEd>
              </div>
              <div v-else>
                <div class="container.is-fullhd mb-3" v-for="project in data" :key="project.ID">
                  <div class="card events-card">
                    <header class="card-header">
                      <p class="card-header-title">
                        用户：{{project.Remarks}}--{{project.UserName}} 项目名称：{{project.ProjectsName}}
                      </p>
                      <button class="card-header-icon">
                        帐号总数: {{project.Count}}
                      </button>
                    </header>
                    <div class="card-content">
                      <div class="table-container">
                        <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                          <thead class="is-size-7">
                            <tr>
                              <td>
                                <label class="checkbox"><input type="checkbox" @click="(e)=>{checkall(e, project.ID)}" />全选</label>
                              </td>
                              <td>序号</td>
                              <td>帐号</td>
                              <td v-if="project.Accounts[0].PhoneNumber.length > 0">手机号</td>
                              <td>今日金币</td>
                              <td>昨日金币</td>
                              <td>炮台</td>
                              <td>钻石</td>
                              <td>狂暴</td>
                              <td>冰冻</td>
                              <td>瞄准</td>
                              <td v-if="project.Accounts[0].Price.length > 0">价格</td>
                              <td>过期时间</td>
                              <td>更新时间</td>
                            </tr>
                          </thead>
                          <tbody class="is-size-7">
                            <tr v-for="(item, index) in project.Accounts" :key="item.ID" class="hasimg" :class="checkTemp.indexOf(item.ID) !== -1  ? 'hasClick' : ''">
                              <td>
                                <label class="checkbox">
                                  <input type="checkbox" v-model="item.check" @click="(e)=>checkBox(e, item.ID, project.ID)">
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
                              <td class="potd">
                                <FormaTime :DateTime="item.UpdatedAt" />
                                <div v-if="item.Cover.length > 0" class="poimg">
                                  <img :src="IMGUri+item.Cover" />
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
            </div>
          </div>
        </div>
      </div>
    </div>
    <NotIfication
      :showData="openerr">
    </NotIfication>
    <RenewalCard
      :showData="openModal"
      :Close="closeModal"
      :ShowMessage="ShowMessage" />
  </div>
</template>
<script>
import { reactive, toRefs, onMounted, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'
import EmptyEd from '@/components/Other/Empty'
import NotIfication from "@/components/Other/Notification"
import FormaTime from '@/components/Other/FormaTime'
import FormaNumber from '@/components/Other/FormaNumber'
import ExpTime from '@/components/Other/ExpTime'
import RenewalCard from '@/components/Other/Renewal'
import PopoButton from '@/components/Other/PopoButton'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AllDraw',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, FormaTime, FormaNumber, ExpTime, RenewalCard, PopoButton },
  setup() {
    let states = reactive({
      AccountKey: "",
      AccountType: "",
      GamesList: [],
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
      openModal:{
        active: false,
        message: "",
        title: "",
        data: "",
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
        GameID: 0,
      },
      IMGUri: Config.IMGUri
    })
    const router = useRouter()
    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-帐号管理`
      
      const data = await CheckLogin()
      if (data == 0) {
        const username = localStorage.getItem('user')
        states.username = username
        GetDateList()
      }else{
        setStorage(false)
        CleanData()
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


    const getDateData = (date) => {
      states.CurrentDate = date
      GetDateList()
    }

    const GetDateList = async() => {
      const token = localStorage.getItem("token")
      states.loading = true
      const d = await Fetch(Config.Api.AllCount, {}, 'GET', token)
      if (d.status == 0) {
        states.GamesList = d.gamslist
        states.pageLoading = true
        states.loading = false
      }else{
        states.GamesList = []
        states.data = []
        states.temp = []
        states.total = 0
        states.page = []
        states.projects = {}
        states.loading = false
      }
    }

    const GetDrawData = async() => {
      const token = localStorage.getItem("token")
      const Filter = states.Filter
      const data = {
        mingold: Filter.mingold * 100000000,
        maxgold: Filter.maxgold * 100000000,
        multiple: Filter.multiple * 10000,
        diamond: Filter.diamond,
        crazy: Filter.crazy,
        cold: Filter.cold,
        precise: Filter.precise,
        gameid: Filter.GameID
      }
      states.loading = true
      const d = await Fetch(Config.Api.drawSelect, data, 'PUT', token)
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
        const e = {
          active: true,
          message: d.message,
          color: "is-danger",
        }
        ShowMessage(e)
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

    const checkBox = (e, account, id) => {
      if (e.target.checked) {
        states.checkTemp = [...states.checkTemp, {'pid':id, 'aid':account}]
      }else{
        states.checkTemp = states.checkTemp.filter((el) => {
          if (el.aid !== account) return el
        })
      }
      // console.log(states.checkTemp)
    }
    const checkall = (e,id) => {
      if (e.target.checked) {
        states.data.map((el) => {
          if (el.ID == id) {
            el.Accounts.map((es) => {
              states.checkTemp =  [...states.checkTemp, {'pid': id, 'aid': es.ID}]
              es.check = true
              return es
            })
          }
        })
        states.checkTemp = [...new Set(states.checkTemp.map(item => 
            JSON.stringify(item)
        ))].map(val => JSON.parse(val))
      }else{
        states.data = states.data.map((el) => {
          if (el.ID == id) {
            el.Accounts.map((es) => {
              es.check =  false
              states.checkTemp = states.checkTemp.filter((ct)=> {
                return ct.aid != es.ID
              })
              return es
            })
          }
          return el
        })
      }
      // console.log(states.checkTemp)
    }
    const makeNumberINT = (n) =>{
      let x = "0"
      if ((n+"").length >= 9 && n !== 0) {
        const a = Math.floor(n / 100000000)
        x = `${a}`
      }else if (n === 123) {
        x = "识别错误"
      }else{
        if (n !== 0 ) {
          const a = Math.floor(n / 10000)
          x = `${a}`
        }
      }
      return x
    }

    const makeData = (data) => {
      let d = []
      data.forEach((el) => {
        el.Accounts.forEach((a) => d = [...d, `${a.UserName}\t${a.Password}\t${makeNumberINT(a.TodayGold)}\t${a.Multiple}\t${el.Projects.ProjectsName}\t${el.Projects.Users.Remarks}\t${el.Projects.Users.UserName}`])
      })
      const x = d.join("\r\n")
      return x
    }

    const closeModal = () => {
      states.checkTemp = []
      GetDrawData()
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
          list: JSON.stringify(states.checkTemp)
        }
        states.buttonLoading = true
        states.pageLoading = false
        states.loading = true
        const d = await Fetch(Config.Api.drawSelectPull, data, 'PUT', token)
        if (d.status == 0) {
          states.loading = false
          states.buttonLoading = false
          // here
          const data = makeData(d.data)
          states.openModal.active = true
          states.openModal.title = "提取成功"
          states.openModal.message = "成功提取帐号，点击复制到剪切板再关闭此弹窗。"
          states.openModal.data = data
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

    const pushRouterToDrawed = () => {
      const AccountKey = router.currentRoute._value.params.key
      router.push(`/accountDrawed/${AccountKey}`)
    }

    return {
      ...toRefs(states),
      backRouter,
      GetDrawData,
      getDateData,
      GetDateList,
      checkBox,
      checkall,
      pullData,
      pushRouterToDrawed,
      closeModal,
      ShowMessage
    }
  },
})
</script>
<style scoped>
.f-1 {
  margin-left: -1px;
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
.hasClick {
  background: #c5e0fd !important
}
.hiddenit {
  display: none !important;
}
</style>