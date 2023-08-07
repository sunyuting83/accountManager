<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
      <div class="com__box" v-if="loading" :style="loading? 'margin-top:5rem':''">
        <LoadIng></LoadIng>
      </div>
      <div v-else>
        <nav class="columns flex-wrap is-flex-wrap-wrap is-justify-content-center pt-3" v-if="GamesList.length !== 0">
          <div class="column is-3 has-text-centered" v-for="item in GamesList" :key="item.ID">
            <div class="is-full">
              <p class="heading">{{item.GameName}}</p>
              <p class="title">{{item.Count}}</p>
            </div>
          </div>
        </nav>
        <nav class="columns flex-wrap is-flex-wrap-wrap is-justify-content-center pt-3" v-if="GamesList.length !== 0">
          <div class="column is-3 has-text-centered" v-for="item in GamesList" :key="item.ID" >
            <div class="is-full">
              <p class="heading">{{item.GameName}}日活号</p>
              <p class="title">{{item.AliveCount}}</p>
            </div>
          </div>
        </nav>
      </div>
    </div>
  </div>
</template>
<script>
import { reactive, toRefs, onMounted, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'


import Fetch from '@/helper/fetch'
import CheckLogin from '@/helper/checkLogin'
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
export default defineComponent({
  name: 'AllDraw',
  components: { ManageHeader, LoadIng },
  setup() {
    let states = reactive({
      GamesList: [],
      loading: false,
      username: "",
    })
    const router = useRouter()
    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-帐号信息`
      
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
      states.GamesList = []
      states.loading = false
    }

    const GetDateList = async() => {
      const token = localStorage.getItem("token")
      states.loading = true
      const d = await Fetch(Config.Api.AllCount, {}, 'GET', token)
      console.log(d)
      if (d.status == 0) {
        states.GamesList = d.gamslist
        states.loading = false
      }else{
        states.GamesList = []
        states.loading = false
      }
    }

    return {
      ...toRefs(states)
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