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
            <input type="checkbox" v-model="form.hastatus">
            自定义帐号状态
          </label>
        </div>
        <div class="field"  v-if="form.hastatus">
          <div class="columns is-flex flex-wrap is-flex-wrap-wrap is-align-content-flex-start">
            <div class="field column newP"  v-for="(item) in form.StatusJSON" :key="item.status">
              <div class="field has-addons ">
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
              <div class="control">
                <label class="checkbox mr-4">
                  <input type="checkbox" v-model="item.delete">
                  可删除
                </label>
                <label class="checkbox mr-4">
                  <input type="checkbox" v-model="item.export">
                  可导出
                </label>
                <label class="checkbox mr-4">
                  <input type="checkbox" v-model="item.import">
                  可导入
                </label>
                <label class="checkbox mr-4">
                  <input type="checkbox" v-model="item.pull">
                  可提取
                </label>
                <label class="checkbox">
                  <input type="checkbox" v-model="item.callback">
                  可退回至
                  <input type="text" class="input inputWidth1 is-small" :disabled="item.callback ? false: true" v-model="item.backto">
                  状态
                </label>
                <label class="checkbox mr-4">
                  <input type="checkbox" v-model="item.ignore">
                  忽略主库
                </label>
              </div>
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
  name: 'AddProject',
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
        StatusJSON: [
        {
          "status": "0",
          "title":"未注册状态",
          "delete":   true,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": true,
          "pull": false,
          "ignore": false,
        },{
          "status": "1",
          "title":"注册中状态",
          "delete":   false,
          "callback": true,
          "backto":   "0",
          "export":   false,
          "import": false,
          "pull": false,
          "ignore": false,
        },{
          "status": "2",
          "title":"注册完成状态",
          "delete":   false,
          "callback": false,
          "backto":   "",
          "export":   false,
          "import": true,
          "pull": true,
          "ignore": false,
        },{
          "status": "3",
          "title":"游戏中状态",
          "delete":   false,
          "callback": true,
          "backto":   "2",
          "export":   false,
          "import": false,
          "pull": true,
          "ignore": false,
        },{
          "status": "4",
          "title":"游戏完成状态",
          "delete":   false,
          "callback": true,
          "backto":   "2",
          "export":   false,
          "import": false,
          "pull": true,
          "ignore": false,
        },{
          "status": "5",
          "title":"封号状态",
          "delete":   false,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": false,
          "pull": false,
          "ignore": true,
        },{
          "status": "6",
          "title":"旧帐号状态",
          "delete":   false,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": false,
          "pull": false,
          "ignore": true,
        },{
          "status": "7",
          "title":"备用状态",
          "delete":   false,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": false,
          "pull": false,
          "ignore": true,
        },{
          "status": "8",
          "title":"未使用身份证",
          "delete":   true,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": true,
          "pull": false,
          "ignore": true,
        },{
          "status": "9",
          "title":"使用中身份证",
          "delete":   false,
          "callback": true,
          "backto":   "8",
          "export":   false,
          "import": false,
          "pull": false,
          "ignore": true,
        },{
          "status": "10",
          "title":"已使用身份证",
          "delete":   true,
          "callback": true,
          "backto":   "8",
          "export":   true,
          "import": false,
          "pull": false,
          "ignore": true,
        },{
          "status": "108",
          "title":"提号状态",
          "delete":   false,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": false,
          "pull": false,
          "ignore": true,
        },{
          "status": "11",
          "title":"IP地址",
          "delete":   true,
          "callback": false,
          "backto":   "",
          "export":   true,
          "import": true,
          "pull": false,
          "ignore": true,
        },{
          "status": "12",
          "title":"已使用IP",
          "delete":   false,
          "callback": true,
          "backto":   "11",
          "export":   false,
          "import": false,
          "pull": false,
          "ignore": true,
        }
      ],
        hastatus: false
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
      const hastatus = _data.form.hastatus
      let StatusJSON = ""
      if (hastatus) StatusJSON = JSON.stringify(_data.form.StatusJSON)
      const token = localStorage.getItem("token")
      let data = {
        usersid : String(UserID),
        ProjectsName: ProjectsName,
        username: UserName,
        password: Password,
        AccNumber: AccNumber,
        ColaAPI: String(ColaAPI),
        StatusJSON: StatusJSON
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
      _data.hastatus = false
      _data.form.StatusJSON = []
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
    
    const addStatus = () => {
      const len = _data.form.StatusJSON.length
      _data.form.StatusJSON = [..._data.form.StatusJSON, {
        "status": String(len) - 1,
        "title": "",
        "delete":   false,
        "callback": false,
        "backto":   "",
        "export":   false
      }]
    }

    return {
      ...toRefs(_data),
      props,
      closErr,
      onSubmit,
      checkProjectsName,
      checkPassword,
      checkUserName,
      checkaccNumber,
      addStatus
    }
  }
})
</script>
<style>
.inputWidth {
  width: 144px
}
.inputWidth1 {
  width: 30px;
  height: 21px;
}
.newP {
  padding: 0.45rem 0rem 0.45rem 0.75rem;
  margin-right: 0.5rem;
  background: #f9f9f9;
}
</style>