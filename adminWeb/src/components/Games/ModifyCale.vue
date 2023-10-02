<template>
  <div class="modal is-clipped" :class="this.showData.active ? 'is-active': ''" v-if="this.showData.active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">修改价格公式</p>
        <button class="delete" aria-label="close" @click="closErr" v-if="loading ? false : true"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <p class="help has-text-left is-danger">计算公式： 底价 + (每金币数量单价 ÷ 每金币数量) × 金币总数 = 帐号价格</p>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">底价</label>
          </div>
          <div class="field-body">
            <div class="field">
              <p class="control has-icons-left has-icons-right">
                <input class="input" type="number" v-model="form.bprice" placeholder="底价" step="0.01">
                <span class="icon is-small is-left">
                  <i class="fa fa-cny"></i>
                </span>
              </p>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">每*金币单价</label>
          </div>
          <div class="field-body">
            <div class="field">
              <p class="control has-icons-left has-icons-right">
                <input :class="form.upriceErr ? 'input is-danger': 'input'" type="number" v-model="form.uprice" placeholder="每*金币单价" step="0.01" :onBlur="checkuprice">
                <span class="icon is-small is-left">
                  <i class="fa fa-cny"></i>
                </span>
                <span class="icon is-small is-right" v-if="form.upriceErr">
                  <i class="fa fa-exclamation-triangle"></i>
                </span>
              </p>
              <p class="help has-text-left is-danger" v-if="form.upriceErr">{{form.upriceErrMessage}}</p>
            </div>
          </div>
        </div>
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">每*金币数量(亿)</label>
          </div>
          <div class="field-body">
            <div class="field">
              <p class="control has-icons-left has-icons-right">
                <input :class="form.numberErr ? 'input is-danger': 'input'" type="number" v-model="form.number" placeholder="每*金币数量" :onBlur="checknumber">
                <span class="icon is-small is-left">
                  <i class="fa fa-file-excel-o"></i>
                </span>
                <span class="icon is-small is-right" v-if="form.numberErr">
                  <i class="fa fa-exclamation-triangle"></i>
                </span>
              </p>
              <p class="help has-text-left is-danger" v-if="form.numberErr">{{form.numberErrMessage}}</p>
            </div>
          </div>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-info" @click="onSubmit" :disabled="form.uprice <= 0 && form.number <= 0 ? true : false" :class="loading ? 'is-loading' : ''">修改</button>
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
  name: 'ModifyCale',
  props: {
    showData:{
      active:{
        type: Boolean,
        default: false
      },
      data: {
        type: Object
      }
    },
    ShowMessage:Function,
  },
  setup(props){
    const _this = props
    // console.log(_this.showData.data.GameID)
    let _data = reactive({
      loading: false,
      form:{
        id: String(_this.showData.data.GameID),
        bprice: _this.showData.data.BasePrice,
        uprice: _this.showData.data.UnitPrice,
        upriceErr: false,
        upriceErrMessage: '',
        number: _this.showData.data.SingleNumber,
        numberErr: false,
        numberErrMessage: '',
      }
    })
    const closErr = () => {
      const _this = props
      cleanState()
      _this.showData.message = ""
      _this.showData.GameID = 0
      _this.showData.active = false
    }
    const onSubmit = async() => {
      if (!_data.form.upriceErr && !_data.form.numberErr) {
        postData()
      }else{
        _data.form.upriceErr = true
        _data.form.numberErr = true
        _data.form.upriceErrMessage = "单价不能小于等于0"
        _data.form.numberErrMessage = "数量不能小于等于0"
      }
    }
    const postData = async() => {
      const id = _data.form.id
      const bprice = String(_data.form.bprice)
      const uprice = String(_data.form.uprice)
      const number = _data.form.number
      const token = localStorage.getItem("token")
      const data = {
        id: id,
        bprice: bprice,
        uprice: uprice,
        number: number
      }
      const d = await Fetch(Config.Api.ModifyCalc, data, "POST", token)
      if (d.status === 0) {
        cleanState()
        closErr()
        props.ShowMessage({
          active: true,
          message: d.message,
          color: 'is-success',
          data: d.data
        }, 0)
      }else{
        _data.form.upriceErr = true
        _data.loading = false
        _data.form.upriceErrMessage = d.message
      }
    }

    const cleanState = () => {
      _data.form.bprice= 0.00
      _data.form.number= 0
      _data.form.numberErr= false
      _data.form.numberErrMessage= ''
      _data.form.uprice= 0.00
      _data.form.upriceErr= false
      _data.form.upriceErrMessage= ''
      _data.loading = false
    }

    const checkuprice = () => {
      if (_data.form.uprice <= 0) {
        _data.form.upriceErr = true
        _data.form.upriceErrMessage = "单价不能等于0"
      }else{
        _data.form.upriceErr = false
        _data.form.upriceErrMessage = ""
      }
    }
    const checknumber = () => {
      if (_data.form.number <= 0) {
        _data.form.numberErr = true
        _data.form.numberErrMessage = "数量不能等于0"
      }else{
        _data.form.numberErr = false
        _data.form.numberErrMessage = ""
      }
    }

    return {
      ...toRefs(_data),
      closErr,
      onSubmit,
      checkuprice,
      checknumber
    }
  }
})
</script>
