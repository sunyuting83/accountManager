<template>
  <div class="container">
    <div class="card">
      <header class="card-header">
        <p class="card-header-title">
          导入数据
        </p>
      </header>
      <div class="card-content">
        <div class="content">
          <div class="control">
            分隔符选择：
            <label class="radio">
              <input type="radio" v-model="splitstr" value="0" name="splitstr">
              Tab
            </label>
            <label class="radio">
              <input type="radio" v-model="splitstr" value="1" name="splitstr">
              ----
            </label>
            <label class="radio">
              <input type="radio" v-model="splitstr" value="2" name="splitstr">
              空格
            </label>
          </div>
          <div class="field has-addons">
            <textarea class="textarea" rows="10" v-model="seleteData" placeholder="粘贴要查询的帐号到这里"></textarea>
          </div>
          <div class="buttons">
            <button class="button is-primary" :class="loading?'is-loading':''" :disabled="userlist.length > 0?false:true" @click="PostData">导入</button>
            <button class="button is-light" @click="CleanSeleteData">清空</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { reactive, toRefs, defineComponent, onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'

import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent({
  name: 'AccList',
  components: { },
  setup() {
    const router = useRouter()
    let states = reactive({
      loading: false,
      seleteData: "",
      splitstr: "0",
      title: "",
    })

    const onSubmit = async() => {
      let url = Config.Api.AddUser
      states.loading = true
      const params = {
        username: states.createUser
      }
      const d = await Fetch(url, params, 'POST')
      if (d.status == 0) {
        states.loading = false
        states.createUser = ""
      }else{
        states.createUser = ""
        states.loading = false
      }
    }

    const PostData = async() => {
      let url = Config.Api.PostAccount
      states.loading = true
      const params = {
        data: states.seleteData,
        splitstr: states.splitstr
      }
      const d = await Fetch(url, params, 'POST')
      if (d.status == 1) color = "is-danger"
      states.loading = false
    }
    const CleanSeleteData = () => {
      states.seleteData = ""
    }

    return {
      ...toRefs(states),
      onSubmit,
      CleanSeleteData,
      PostData
    }
  },
})
</script>