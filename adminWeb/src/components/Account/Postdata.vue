<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">导入{{showData.postParams.title}}数据</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <div class="modal-card-body">
        <div class="field">
          <div class="content">
            <div class="mb-5">
              <div class="icon-text" @click="showNote">
                <span class="icon has-text-danger">
                  <i class="fa fa-exclamation-triangle"></i>
                </span>
                <span>点击{{shownote?'查看':'关闭'}}导入说明</span>
              </div>
              <p class="block" v-if="shownote">
                上传过程中不要关闭页面，并耐心等待<br />
                一次数据导入数量限制在3万条以内，超过3万的数据请分批导入。
                使用文本导入，数量限制在1万条以内，超过1万条的数据请使用文件导入。
                可以选择分隔符或自定义分隔符，根据文本自行选择。<br />
                没有密码的帐号 密码必须为空格或其他任意字符。<br />
                数据中带有手机号的会自动判断并保存，手机号分割符后的数据自动保存手机密码<br />
                空格密码的例子：<br />
                byyUBPPCCJWJDQRERLEEY----password<br />
                带手机号的例子：<br />
                byyUBPPCCJWJDQRERLEEY----13613231234----password<br />
              </p>
            </div>
            <div class="control mb-3">
              <label class="checkbox mr-4">
                <input type="checkbox" v-model="repeated">
                自动去重复
              </label>
              <label class="checkbox mr-4">
                <input type="checkbox" v-model="hasmore">
                格式化数据导入
              </label>
            </div>
            <div class="mb-5" v-if="repeated">
              <div class="icon-text">
                <span class="icon has-text-danger">
                  <i class="fa fa-exclamation-triangle"></i>
                </span>
                <span>自动去重使用说明</span>
              </div>
              <p class="block">
                只会过滤第一列的数据（帐号数据）。<br />
                会自动去重已上传过的数据。(封号、备用、旧帐号等状态不在过滤范围内)<br />
                使用此功能会加大服务器压力，去重过程很漫长，上传过程中不要关闭页面，并耐心等待
              </p>
            </div>
            <div class="mb-5" v-if="hasmore">
              <div class="icon-text">
                <span class="icon has-text-danger">
                  <i class="fa fa-exclamation-triangle"></i>
                </span>
                <span>格式化数据使用说明</span>
              </div>
              <p class="block">
                导入的数据必须严格按照数据格式导入，否则会导致数据错乱、无法统计等不可预测的问题。此功能用于帐号拥有金币、炮台、狂暴、瞄准等游戏数据的情况。<br />
                数据格式为：<br />
                帐号----密码----金币----炮台----钻石----狂暴----瞄准----冰冻----其他数据<br />
                “----”（不包含引号）是分隔符，依然可以选择分隔符或自定义分隔符，根据文本自行选择。<br />
                其他数据 会自动保存到 帐号其他选项中 依然以文本形式保存<br />
                空格密码的例子：<br />
                byyUBPPCCJWJDQRERLEEY---- ----405亿----200000----200----35----50----80<br />
              </p>
            </div>
            <div class="control">
              分隔符选择：
              <label class="radio mr-2">
                <input type="radio" v-model="splitstr" value="0" name="splitstr">
                Tab
              </label>
              <label class="radio mr-2">
                <input type="radio" v-model="splitstr" value="1" name="splitstr">
                ----
              </label>
              <label class="radio mr-2">
                <input type="radio" v-model="splitstr" value="2" name="splitstr">
                空格
              </label>
              <div class="field has-addons mb-2">
                <p class="control mr-1">
                  <label class="checkbox">
                    <input type="radio" name="splitstr">
                    自定义分隔符
                  </label>
                </p>
                <p class="control">
                  <input type="text" class="input inputWidth1 is-small" v-model="splitstr">
                </p>
              </div>
            </div>
            <div class="mb-5" v-if="errStatus">
                <div class="icon-text">
                  <span class="icon has-text-danger">
                    <i class="fa fa-ban"></i>
                  </span>
                  <span>{{errMessage}}</span>
                </div>
            </div>
            <div class="tabs is-large">
              <ul class="ml-0">
                <li :class="hasfile?'':'is-active'" @click="selectTab"><a>文本上传</a></li>
                <li :class="hasfile?'is-active':''" @click="selectTab"><a>文件上传</a></li>
              </ul>
            </div>
            <div class="field has-addons" v-if="!hasfile">
              <textarea class="textarea" rows="10" v-model="seleteData" placeholder="粘贴要查询的帐号到这里"></textarea>
            </div>
            <div class="file is-large is-boxed has-name mb-3" v-else>
              <label class="uploaders file-label">
                <input class="file-input" type="file" name="files" @change="(e)=> {uploaders(e)}">
                <span class="file-cta">
                  <span class="file-icon">
                    <i class="fa fa-upload"></i>
                  </span>
                  <span class="file-label">
                    选择文件…
                  </span>
                </span>
                <span class="file-name">
                  请点击选择文件
                </span>
              </label>
            </div>
            <div class="buttons">
              <button class="button is-primary" :disabled="isfile.size > 0 && isfile.type == 'text/plain'?false:true" :class="loading?'is-loading':''" v-if="hasfile" @click="onSubmit">上传</button>
              <button class="button is-primary"  v-if="!hasfile"  :disabled="seleteData.length > 0?false:true" :class="loading?'is-loading':''" @click="onSubmit">导入</button>
              <button class="button is-light" v-if="!hasfile" @click="CleanSeleteData">清空</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { reactive, toRefs, defineComponent } from 'vue'
import { useRouter } from 'vue-router'

import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent({
  name: 'PostData',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      },
      postParams: {
        type: Object,
      }
    },
    ShowMessage:Function
  },
  components: { },
  setup(props) {
    let states = reactive({
      loading: false,
      seleteData: "",
      splitstr: "0",
      hasmore: false,
      repeated: false,
      hasfile: false,
      shownote: false,
      isfile: {},
      errStatus: false,
      errMessage: ""
    })
    const router = useRouter()
    const closErr = () => {
      const _this = props
      cleanState()
      _this.showData.message = ""
      _this.showData.active = false
      _this.showData.postParams = {}
      _this.ShowMessage({
        active: false,
        message: "",
        color: 'is-success',
        data: []
      })
    }
    const onSubmit = async() => {
      const token = localStorage.getItem("token")
      const AccountKey = router.currentRoute._value.params.key
      const _this = props
      states.loading = true
      const url = `${Config.RootUrl}${AccountKey}/PostAccount`
      
      let params = {
        data: states.seleteData,
        splitstr: states.splitstr,
        status: _this.showData.postParams.status,
        hasmore: String(states.hasmore),
        hasfile: String(states.hasfile),
        repeated: String(states.repeated),
      }
      if (states.hasfile) {
        params['data'] = ""
        params['files'] = states.isfile
      }
      const d = await Fetch(url, params, 'POST', token)
      if (d.status == 0) {
        states.loading = false
        cleanState()
        _this.showData.message = ""
        _this.showData.active = false
        _this.showData.postParams = {}
        _this.ShowMessage({
          active: false,
          message: "",
          color: 'is-success'
        }, 4)
      }else{
        states.loading = false
        states.errStatus = true
        states.errMessage = d.message
      }
    }

    const CleanSeleteData = () => {
      states.seleteData = ""
    }

    const cleanState = () => {
      states.loading = false
      states.seleteData = ""
      states.splitstr = "0"
      states.hasmore = false
      states.repeated = false
      states.hasfile = false
      states.shownote = false
      states.isfile = {}
      states.errStatus = false
      states.errMessage = ""
    }

    const selectTab = () => {
      if (states.hasfile) {
        states.seleteData = ""
      }else {
        states.isfile = {}
      }
      states.hasfile = !states.hasfile
    }

    const showNote = () => {
      states.shownote = !states.shownote
    }

    const uploaders = (e) => {
      if (e.target.files.length > 0) {
        if (e.target.files[0].type == "text/plain") {
          states.isfile = e.target.files[0]
        }else {
          states.errStatus = true
          states.errMessage = "必须是.txt文件"
        }
        // console.log(states.isfile)
        // console.log(states.isfile.size)
        // console.log(states.isfile.type)
      }
    }

    return {
      ...toRefs(states),
      closErr,
      onSubmit,
      CleanSeleteData,
      selectTab,
      showNote,
      uploaders
    }
  },
})
</script>
<style scoped>
.uploaders {width: 100% !important}
.uploaders .file-label {width: 100% !important; text-align: center;}
.uploaders .file-name {width: 100% !important; max-width: 100% !important;text-align: center;}
</style>