<template>
  <div class="modal" :class="showData.active ? 'is-active': ''" v-if="showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">帐号信息</p>
        <button class="delete" aria-label="close" @click="closErr"></button>
      </header>
      <section class="modal-card-body">
        <p>帐号：{{showData.account.accountgs}}</p>
        <p>密码：{{showData.account.password}}</p>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-success" @click="copyAccount">复制到剪切板</button>
        <button class="button" @click="closErr">关闭</button>
      </footer>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import useClipboard from 'vue-clipboard3'
const { toClipboard } = useClipboard()
export default defineComponent ({
  name: 'RenewalCard',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      },
      account: {
        accountgs: {
        type: String
        },
        password: {
          type: String
        },
      },
      atest : {
        type: String
      }
    },
    ShowMessage:Function
  },
  setup(props){
    const closErr = () => {
      const _this = props
      _this.showData.active = false
    }
    const copyAccount = async() => {
      const p = props
      await toClipboard(p.showData.atest)
      const d = {
        active: true,
        message: "复制帐号成功",
        color: "is-success"
      }
      props.ShowMessage(d)
    }
    
    return {
      closErr,
      copyAccount
    }
  }
})
</script>
