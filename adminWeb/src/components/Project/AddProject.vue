<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">添加项目</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <section class="modal-card-body">
        <div class="field" v-if="data.length !== 0">
          <div class="select is-normal">
            <select v-model="form.UserID">
              <option v-for="(item) in data" :key="item.ID" :value="item.ID">{{item.UserName}}</option>
            </select>
          </div>
        </div>
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input :class="form.ProjectsNameErr ? 'input is-danger': 'input'" type="ProjectsName" v-model="form.ProjectsName" placeholder="项目名称" :onBlur="checkProjectsName">
            <span class="icon is-small is-left">
              <i class="fa fa-user-circle-o"></i>
            </span>
            <span class="icon is-small is-right" v-if="form.ProjectsNameErr">
              <i class="fa fa-exclamation-triangle"></i>
            </span>
          </p>
          <p class="help has-text-left is-danger" v-if="form.ProjectsNameErr">{{form.ProjectsNameErrMessage}}</p>
        </div>
        <div class="field">
          <label class="checkbox">
            <input type="checkbox" v-model="form.cola">
            可乐API
          </label>
        </div>
        <div v-if="form.cola">
          <div class="field">
            <p class="control has-icons-left has-icons-right">
              <input :class="form.UserNameErr ? 'input is-danger': 'input'" type="ProjectsName" v-model="form.UserName" placeholder="平台用户名" :onBlur="checkUserName">
              <span class="icon is-small is-left">
                <i class="fa fa-user-circle-o"></i>
              </span>
              <span class="icon is-small is-right" v-if="form.UserNameErr">
                <i class="fa fa-exclamation-triangle"></i>
              </span>
            </p>
            <p class="help has-text-left is-danger" v-if="form.UserNameErr">{{form.UserNameErrMessage}}</p>
          </div>
          <div class="field">
            <p class="control has-icons-left has-icons-right">
              <input :class="form.PasswordErr ? 'input is-danger': 'input'" type="text" v-model="form.Password" placeholder="API密码" :onBlur="checkPassword">
              <span class="icon is-small is-left">
                <i class="fa fa-user-circle-o"></i>
              </span>
              <span class="icon is-small is-right" v-if="form.PasswordErr">
                <i class="fa fa-exclamation-triangle"></i>
              </span>
            </p>
            <p class="help has-text-left is-danger" v-if="form.PasswordErr">{{form.PasswordErrMessage}}</p>
          </div>
          <div class="field">
            <p class="control has-icons-left has-icons-right">
              <input :class="form.accNumberErr ? 'input is-danger': 'input'" type="number" v-model="form.accNumber" placeholder="帐号数量" :onBlur="checkaccNumber">
              <span class="icon is-small is-left">
                <i class="fa fa-user-circle-o"></i>
              </span>
              <span class="icon is-small is-right" v-if="form.accNumberErr">
                <i class="fa fa-exclamation-triangle"></i>
              </span>
            </p>
            <p class="help has-text-left is-danger" v-if="form.accNumberErr">{{form.accNumberErrMessage}}</p>
          </div>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-info" @click="onSubmit" :class="loading ? 'is-loading' : ''">添加</button>
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
  name: 'AddUser',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      }
    },
    ShowMessage:Function,
    UserData: Array
  },
  setup(props){
    let _data = reactive({
      loading: false,
      userLoading: false,
      data: props.UserData,
      form:{
        ProjectsName: "",
        ProjectsNameErr: false,
        ProjectsNameErrMessage: '',
        UserName: "",
        UserNameErr: false,
        UserNameErrMessage: '',
        Password: "",
        PasswordErr: false,
        PasswordErrMessage: '',
        accNumber: 0,
        accNumberErr: false,
        accNumberErrMessage: '',
        cola: false,
        UserID: 0,
      }
    })
    const closErr = () => {
      const _this = props
      cleanState()
      _this.showData.message = ""
      _this.showData.active = false
      props.ShowMessage({
        active: false,
        message: "",
        color: 'is-success',
        data: []
      }, 4)
    }
    const onSubmit = async() => {
      postData()
    }
    const postData = async() => {
      const ProjectsName = _data.form.ProjectsName
      const UserID = _data.form.UserID
      const ColaAPI = _data.form.cola
      const UserName = _data.form.UserName
      const Password = _data.form.Password
      const AccNumber = _data.form.accNumber
      const token = localStorage.getItem("token")
      let data = {
        usersid : String(UserID),
        ProjectsName: ProjectsName,
        username: UserName,
        password: Password,
        AccNumber: AccNumber,
        ColaAPI: String(ColaAPI),
      }
      const d = await Fetch(Config.Api.addproject, data, "POST", token)
      if (d.status === 0) {
        cleanState()
        closErr()
        props.ShowMessage({
          active: true,
          message: d.message,
          color: 'is-success',
          data: d.data,
          userLoading: false
        }, 1)
      }else{
        _data.form.ProjectsNameErr = true
        _data.loading = false
        _data.form.ProjectsNameErrMessage = d.message
      }
    }

    const cleanState = () => {
      _data.form.ProjectsName = ""
      _data.form.ProjectsNameErr = false
      _data.form.ProjectsNameErrMessage = ""
      _data.form.UserName = ""
      _data.form.UserNameErr = false
      _data.form.UserNameErrMessage = ""
      _data.form.Password = ""
      _data.form.PasswordErr = false
      _data.form.PasswordErrMessage = ""
      _data.form.accNumber = ""
      _data.form.accNumberErr = false
      _data.form.accNumberErrMessage = ""
      _data.form.cola = false
      _data.data = []
      _data.loading = false
    }

    const checkProjectsName = () => {
      if (_data.form.ProjectsName.length < 4) {
        _data.form.ProjectsNameErr = true
        _data.form.ProjectsNameErrMessage = "项目名不能小于4位"
      }else{
        _data.form.ProjectsNameErr = false
        _data.form.ProjectsNameErrMessage = ""
      }
    }
    const checkUserName = () => {
      if (_data.form.UserName.length < 4) {
        _data.form.UserNameErr = true
        _data.form.UserNameErrMessage = "项目名不能小于4位"
      }else{
        _data.form.UserNameErr = false
        _data.form.UserNameErrMessage = ""
      }
    }
    const checkPassword = () => {
      if (_data.form.Password.length < 8) {
        _data.form.PasswordErr = true
        _data.form.PasswordErrMessage = "密码必须大于8位"
      }else{
        _data.form.PasswordErr = false
        _data.form.PasswordErrMessage = ""
      }
    }
    const checkaccNumber = () => {
      if (_data.form.accNumber.length < 8) {
        _data.form.accNumberErr = true
        _data.form.accNumberErrMessage = "密码必须大于8位"
      }else{
        _data.form.accNumberErr = false
        _data.form.accNumberErrMessage = ""
      }
    }
    

    return {
      ...toRefs(_data),
      props,
      closErr,
      onSubmit,
      checkProjectsName,
      checkPassword,
      checkUserName,
      checkaccNumber
    }
  }
})
</script>
