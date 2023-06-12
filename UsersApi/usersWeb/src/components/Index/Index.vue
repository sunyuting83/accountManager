<template>
  <div class="hero is-fullheight is-primary" :style="`background:url('${background}'); background-size:100% 100%;`">
    <div class="hero-body">
      <div class="container has-text-centered">
        <div class="column is-8 is-offset-2">
          <h3 class="title has-text-white">后台管理系统</h3>
          <hr class="login-hr">
          <p class="subtitle has-text-white">感谢今天正在努力的自己!</p>
          <div class="box">
            <div class="box">
              <img :src="logo">
            </div>
            <div class="title has-text-grey is-5">请输入用户名和密码.</div>
          </div>
          <div class="field">
            <p class="control has-icons-left has-icons-right">
              <input class="input" type="username" v-model="form.username" placeholder="用户名" :onBlur="checkUsername">
              <span class="icon is-small is-left">
                <i class="fa fa-user-circle-o"></i>
              </span>
              <span class="icon is-small is-right" v-if="form.usernameErr !== '1' && form.usernameErr.length > 1">
                <i class="fa fa-exclamation-triangle"></i>
              </span>
            </p>
            <p class="help has-text-left is-white" v-if="form.usernameErr !== '1' && form.usernameErr.length > 1">{{form.usernameErr}}</p>
          </div>
          <div class="field">
            <p class="control has-icons-left has-icons-right">
              <input class="input" type="password" v-model="form.password" placeholder="密码" :onBlur="checkPassword">
              <span class="icon is-small is-left">
                <i class="fa fa-lock"></i>
              </span>
              <span class="icon is-small is-right" v-if="form.passwordErr !== '1' && form.passwordErr.length > 1">
                <i class="fa fa-exclamation-triangle"></i>
              </span>
            </p>
            <p class="help has-text-left is-white" v-if="form.passwordErr !== '1' && form.passwordErr.length > 1">{{form.passwordErr}}</p>
          </div>
          <div class="field has-addons">
            <div class="control is-expanded">
              <div class="control has-icons-left has-icons-right">
                <input class="input" type="code" maxlength="4" v-model="form.code" placeholder="请输入验证码" @input="validateCode" />
                <span class="icon is-small is-left">
                  <i class="fa fa-qrcode"></i>
                </span>
                <span class="icon is-small is-right" v-if="form.codeErr !== '1' && form.codeErr.length == 4">
                  <i class="fa fa-exclamation-triangle"></i>
                </span>
              </div>
              <p class="help has-text-left is-white" v-if="form.codeErr !== '1' && form.codeErr.length > 1">{{form.codeErr}}</p>
            </div>
            <div class="control">
              <Identify :identifyCode="identifyCode" @click="refreshCode" />
            </div>
          </div>
          <button class="button is-block is-info is-large is-fullwidth" :disabled="form.usernameErr !== '1' && form.usernameErr.length === 0  && form.passwordErr !== '1' && form.passwordErr.length === 0  && form.codeErr !== '1' && form.codeErr.length === 0 ? false : true" @click="onSubmit">登陆</button>
        </div>
      </div>
    </div>
    <NotIfication
      :showData="openerr">
    </NotIfication>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { reactive, toRefs, onMounted, defineComponent } from 'vue'

import Config from '@/helper/config'
import Fetch from '@/helper/fetch'
import setStorage from '@/helper/setStorage'
import CheckLogin from '@/helper/checkLogin'
import Identify from "@/components/Other/Identify"
import NotIfication from "@/components/Other/Notification"
export default defineComponent({
  name: 'AppIndex',
  components: { NotIfication, Identify },
  setup() {
    let states = reactive({
      identifyCodes: "1234567890",
      identifyCode: "",
      logo: Config.images[1],
      background: Config.images[4],
      form: {
        username: "",
        password: "",
        code: "",
        usernameErr: "1",
        passwordErr: "1",
        codeErr: "1",
      },
      openerr: {
        active: false,
        message: "",
        color: "is-danger"
      }
    })
    const router = useRouter()

    onMounted(async() => {
      document.title = `${Config.GlobalTitle}-登陆`
      const data = await CheckLogin()
      if (data === 1) {
        states.identifyCode = ""
        makeCode(states.identifyCodes, 4)
      }else {
        router.push("project")
      }
    })

    const onSubmit = async() => {
      const { username, password } = states.form
      if (username.length >= 4 && password.length >= 8 && validateCode()) {
        const data = {
          username: username,
          password: password
        }
        const d = await Fetch(Config.Api.login, data, "post")
        if (d.status === 0) {
          // console.log(d)
          setStorage(true, d.token, d.username, String(d.userid))
          router.push("project")
        }else{
          states.openerr.message = d.message
          states.openerr.active = true
          states.form.codeErr = "1"
          states.form.code = ""
          setStorage(false)
          refreshCode()
        }
      }
    }
    const checkUsername = () =>{
      const username = states.form.username
      if (username.length < 4) {
        states.form.usernameErr = "用户名不能小于4位"
      }else{
        states.form.usernameErr = ""
      }
    }
    const checkPassword = () =>{
      const password = states.form.password
      if (password.length < 8) {
        states.form.passwordErr = "用户名不能小于4位"
      }else{
        states.form.passwordErr = ""
      }
    }
    const validateCode = () => {
      const value = states.form.code
      if (value === "") {
        states.form.codeErr = "验证码为空"
        return false;
      } else if (value !== states.identifyCode) {
        states.form.codeErr = "验证码不正确"
        return false;
      } else {
        states.form.codeErr = ""
        return true
      }
    }
        // 验证码相关
    const randomNum = (min, max) => {
      return Math.floor(Math.random() * (max - min) + min);
    }
        // 点击刷新验证码
    const refreshCode = () => {
      states.identifyCode = "";
      makeCode(states.identifyCodes, 4);
    }
    
        // 生成验证码
    const makeCode = (o, l) => {
      for (let i = 0; i < l; i++) {
        states.identifyCode +=
          states.identifyCodes[randomNum(0, states.identifyCodes.length)];
      }
    }

    return {
      ...toRefs(states),
      onSubmit,
      checkUsername,
      checkPassword,
      validateCode,
      refreshCode
    }
  }
})
</script>