<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">修改项目</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <p class="control has-icons-left has-icons-right">
            <input type="text" :placeholder="form.ProjectName" disabled>
            <input type="hidden" :placeholder="form.ID" disabled>
          </p>
        </div>
        <div class="field">
          <label class="checkbox">
            <input type="checkbox" v-model="form.hastatus">
            自定义帐号状态
          </label>
        </div>
        <div class="field"  v-if="form.hastatus">
          <div class="columns is-flex flex-wrap is-flex-wrap-wrap is-align-content-flex-start">
            <div class="field has-addons column newP"  v-for="(item) in form.StatusJSON" :key="item.status">
              <p class="control">
                <input class="input inputWidth is-small" type="hidden" v-model="item.status">
                <input class="input inputWidth is-small" type="text" v-model="item.title">
              </p>
              <p class="control">
                <span class="button is-static is-small">
                  {{item.status}}
                </span>
              </p>
            </div>
            <div class="column">
              <button class="button is-small is-info" @click="addStatus">
                添加新状态
              </button>
            </div>
          </div>
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
        <button class="button is-info" @click="onSubmit" :class="loading ? 'is-loading' : ''">修改</button>
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
  name: 'ModifyProject',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      },
      Project: {
        type: Object
      }
    },
    ShowMessage:Function
  },
  setup(props){
    let StatusJSON = JSON.parse(props.showData.Project.StatusJSON)
    console.log(StatusJSON)
    let haStatus = false
    if (StatusJSON.length > 0) haStatus = true
    let _data = reactive({
      loading: false,
      form:{
        ID: props.showData.Project.ID,
        ProjectName: props.showData.Project.ProjectsName,
        UserName: props.showData.Project.UserName,
        UserNameErr: false,
        UserNameErrMessage: '',
        Password: props.showData.Project.Password,
        PasswordErr: false,
        PasswordErrMessage: '',
        accNumber: props.showData.Project.AccNumber,
        accNumberErr: false,
        accNumberErrMessage: '',
        cola: props.showData.Project.ColaAPI,
        StatusJSON: StatusJSON,
        hastatus: haStatus
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
      }, 5)
    }
    const onSubmit = async() => {
      postData()
    }
    const postData = async() => {
      const ColaAPI = _data.form.cola
      const UserName = _data.form.UserName
      const Password = _data.form.Password
      const AccNumber = _data.form.accNumber
      const StatusJSON = JSON.stringify(_data.form.StatusJSON)
      const token = localStorage.getItem("token")
      let data = {
        id : String(_data.form.ID),
        username: UserName,
        password: Password,
        AccNumber: AccNumber,
        ColaAPI: String(ColaAPI),
        StatusJSON: StatusJSON
      }
      const d = await Fetch(Config.Api.UpdateProjects, data, "PUT", token)
      if (d.status === 0) {
        cleanState()
        closErr()
        props.ShowMessage({
          active: true,
          message: d.message,
          color: 'is-success',
          data: d.data,
          modifyStatus: false
        }, 6)
      }else{
        _data.form.ProjectsNameErr = true
        _data.loading = false
        _data.form.ProjectsNameErrMessage = d.message
      }
    }

    const cleanState = () => {
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
      _data.form.StatusJSON = []
    }

    const checkUserName = () => {
      if (_data.form.UserName.length < 4) {
        _data.form.UserNameErr = true
        _data.form.UserNameErrMessage = "用户名不能小于4位"
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
    
    const addStatus = () => {
      const len = _data.form.StatusJSON.length + 1
      _data.form.StatusJSON = [..._data.form.StatusJSON, {
        "status": String(len),
        "title": ""
      }]
    }

    return {
      ...toRefs(_data),
      props,
      closErr,
      onSubmit,
      checkPassword,
      checkUserName,
      checkaccNumber,
      addStatus
    }
  }
})
</script>
<style scoped>
.inputWidth {
  width: 80px
}
.newP {
  padding: 0.45rem 0.5rem 0.45rem 0.75rem;
}
</style>