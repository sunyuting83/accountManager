<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">添加游戏</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input :class="form.gamenameErr ? 'input is-danger': 'input'" type="gamename" v-model="form.gamename" placeholder="用户名" :onBlur="checkgamename">
            <span class="icon is-small is-left">
              <i class="fa fa-user-circle-o"></i>
            </span>
            <span class="icon is-small is-right" v-if="form.gamenameErr">
              <i class="fa fa-exclamation-triangle"></i>
            </span>
          </p>
          <p class="help has-text-left is-danger" v-if="form.gamenameErr">{{form.gamenameErrMessage}}</p>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-info" @click="onSubmit" :disabled="form.gamename.length <= 4 ? true : false" :class="loading ? 'is-loading' : ''">添加</button>
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
  name: 'AddGame',
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
        gamename: "",
        gamenameErr: false,
        gamenameErrMessage: '',
      }
    })
    const closErr = () => {
      const _this = props
      cleanState()
      _this.showData.message = ""
      _this.showData.active = false
    }
    const onSubmit = async() => {
      if (!_data.form.gamenameErr) {
        postData()
      }else{
        _data.form.gamenameErr = true
        _data.form.passErrMessage = "游戏名必须大于4位"
      }
    }
    const postData = async() => {
      const gamename = _data.form.gamename
      const token = localStorage.getItem("token")
      const data = {
        gamename: gamename
      }
      const d = await Fetch(Config.Api.AddGame, data, "POST", token)
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
        _data.form.gamenameErr = true
        _data.loading = false
        _data.form.gamenameErrMessage = d.message
      }
    }

    const cleanState = () => {
      _data.form.gamename= ""
      _data.form.gamenameErr= false
      _data.form.gamenameErrMessage= ''
      _data.loading = false
    }

    const checkgamename = () => {
      if (_data.form.gamename.length < 4) {
        _data.form.gamenameErr = true
        _data.form.gamenameErrMessage = "用户名不能小于4位"
      }else{
        _data.form.gamenameErr = false
        _data.form.gamenameErrMessage = ""
      }
    }

    return {
      ...toRefs(_data),
      closErr,
      onSubmit,
      checkgamename
    }
  }
})
</script>
