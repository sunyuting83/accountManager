<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">添加管理员</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input :class="form.usernameErr ? 'input is-danger': 'input'" type="username" v-model="form.username" placeholder="用户名" :onBlur="checkUsername">
            <span class="icon is-small is-left">
              <i class="fa fa-user-circle-o"></i>
            </span>
            <span class="icon is-small is-right" v-if="form.usernameErr">
              <i class="fa fa-exclamation-triangle"></i>
            </span>
          </p>
          <p class="help has-text-left is-danger" v-if="form.usernameErr">{{form.usernameErrMessage}}</p>
        </div>
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input :class="form.passwordErr ? 'input is-danger': 'input'" type="password" v-model="form.password" placeholder="密码" :onBlur="checkPassword">
            <span class="icon is-small is-left">
              <i class="fa fa-lock"></i>
            </span>
            <span class="icon is-small is-right" v-if="form.passwordErr">
              <i class="fa fa-exclamation-triangle"></i>
            </span>
          </p>
          <p class="help has-text-left is-danger" v-if="form.passwordErr">{{form.passwordErrMessage}}</p>
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
        <button class="button is-info" @click="onSubmit" :disabled="form.password !== form.repassword ? true : false" :class="loading ? 'is-loading' : ''">添加</button>
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
  name: 'AddAdmin',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      }
    },
    ShowMessage:Function
  },
  setup(props){
    let _data = reactive({
      loading: false,
      form:{
        username: "",
        password: "",
        repassword: '',
        usernameErr: false,
        passwordErr: false,
        passErr: false,
        passErrMessage: '',
        passwordErrMessage: '',
        usernameErrMessage: '',
      }
    })
    const closErr = () => {
      const _this = props
      cleanState()
      _this.showData.message = ""
      _this.showData.active = false
    }
    const onSubmit = async() => {
      const pawsd = _data.form.password
      const rpawsd = _data.form.repassword
      if (!_data.form.usernameErr) {
        if (pawsd == rpawsd && pawsd.length >=8 && rpawsd.length >=8) {
          _data.loading = true
          postData()
        }else{
          _data.form.passErr = true
          _data.form.passErrMessage = "密码必须大于8位"
        }
      }else{
        _data.form.usernameErr = true
        _data.form.passErrMessage = "用户名必须大于4位"
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
      const username = _data.form.username
      const token = localStorage.getItem("token")
      const rpawsd = _data.form.repassword
      const data = {
        username: username,
        password: rpawsd
      }
      const d = await Fetch(Config.Api.addadmin, data, "POST", token)
      if (d.status === 0) {
        cleanState()
        closErr()
        props.ShowMessage({
          active: true,
          message: d.message,
          color: 'is-success',
          data: d.data
        }, 1)
      }else{
        _data.form.usernameErr = true
        _data.loading = false
        _data.form.usernameErrMessage = d.message
      }
    }

    const cleanState = () => {
      _data.form.passErr = false
      _data.form.username= ""
      _data.form.password= ""
      _data.form.repassword= ''
      _data.form.usernameErr= false
      _data.form.passwordErr= false
      _data.form.passErr= false
      _data.form.passErrMessage= ''
      _data.form.passwordErrMessage= ''
      _data.form.usernameErrMessage= ''
      _data.loading = false
    }

    const checkUsername = () => {
      if (_data.form.username.length < 4) {
        _data.form.usernameErr = true
        _data.form.usernameErrMessage = "用户名不能小于4位"
      }else{
        _data.form.usernameErr = false
        _data.form.usernameErrMessage = ""
      }
    }
    const checkPassword = () => {
      if (_data.form.password.length < 8) {
        _data.form.passwordErr = true
        _data.form.passwordErrMessage = "密码必须大于8位"
      }else{
        _data.form.passwordErr = false
        _data.form.passwordErrMessage = ""
      }
    }

    return {
      ...toRefs(_data),
      closErr,
      onSubmit,
      check,
      checkUsername,
      checkPassword
    }
  }
})
</script>
