<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">修改密码</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input class="input" type="password" v-model="form.password" placeholder="密码">
            <span class="icon is-small is-left">
              <i class="fa fa-lock"></i>
            </span>
          </p>
        </div>
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input :class="form.passErr ? 'input is-danger': 'input'" type="password" v-model="form.repassword" @input="check" placeholder="重复密码">
            <span class="icon is-small is-left">
              <i class="fa fa-lock"></i>
            </span>
            <span class="icon is-small is-right" v-if="form.passErr">
              <i class="fa fa-exclamation-triangle"></i>
            </span>
          </p>
          <p class="help has-text-left is-danger" v-if="form.passErr">{{form.passErrMessage}}</p>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-info" @click="onSubmit" :disabled="form.password !== form.repassword ? true : false" :class="loading ? 'is-loading' : ''">修改</button>
        <button class="button" @click="closErr" :disabled="loading ? true : false" :class="loading ? 'is-loading' : ''">取消</button>
      </footer>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, defineComponent } from 'vue'
import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent ({
  name: 'ChangePassword',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      },
      username: {
        type: String,
        default: ""
      }
    },
    ShowMessage:Function,
    Admin: Boolean,
  },
  setup(props){
    let _data = reactive({
      loading: false,
      form:{
        password: '',
        repassword: '',
        passErr: false,
        passErrMessage: ''
      }
    })
    const closErr = () => {
      const _this = props
      _data.form.passErr = false
      _data.form.passErr = ""
      _data.form.password = ""
      _data.form.repassword = ""
      _data.loading = false
      _this.showData.message = ""
      _this.showData.active = false
    }
    const onSubmit = async() => {
      const pawsd = _data.form.password
      const rpawsd = _data.form.repassword
      if (pawsd == rpawsd && pawsd.length >=8 && rpawsd.length >=8) {
        _data.loading = true
        postData()
      }else{
        _data.form.passErr = true
        _data.form.passErrMessage = "密码必须大于8位"
      }
    }
    const check = () => {
      const pawsd = _data.form.password
      const rpawsd = _data.form.repassword
      if (pawsd !== rpawsd) {
        _data.form.passErr = true
        _data.form.passErrMessage = "两次密码不一致"
      }else{
        _data.form.passErr = false
        _data.form.passErr = ""
      }
    }
    const postData = async() => {
      const username = props.showData.username
      const token = localStorage.getItem("token")
      const rpawsd = _data.form.repassword
      const data = {
        username: username,
        password: rpawsd
      }
      const _this = props
      let url = Config.Api.reUserPassword
      if (_this.Admin) url = Config.Api.repassword
      const d = await Fetch(url, data, "put", token)
      if (d.status === 0) {
        _data.form.passErr = false
        _data.form.passErr = ""
        _data.form.password = ""
        _data.form.repassword = ""
        _data.loading = false
        closErr()
        props.ShowMessage({
          active: true,
          message: d.message,
          color: 'is-success'
        })
      }else{
        _data.form.passErr = true
        _data.loading = false
        _data.form.passErrMessage = d.message
      }
    }
    return {
      ...toRefs(_data),
      closErr,
      onSubmit,
      check
    }
  }
})
</script>
