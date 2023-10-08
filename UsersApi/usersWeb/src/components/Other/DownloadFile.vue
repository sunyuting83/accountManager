<template>
  <button class="button is-small" :class="buttonLoading?'is-loading' + styles : styles" @click="ExportAccount">{{title}}</button>
</template>
<script>
import { defineComponent } from 'vue'
import Fetch from '@/helper/fetch'

export default defineComponent({
  name: 'DownloadFile',
  props: {
    uri:{
      type: String,
      default: ""
    },
    Data:{
      type: Object,
      default: () => ({})
    },
    title:{
      type: String,
      default: ""
    },
    buttonLoading:{
      type: Boolean,
      default: false
    },
    styles:{
      type: String,
      default: ""
    },
    ext:{
      type: String,
      default: ".txt"
    }
  },
  setup(props) {
    const token = localStorage.getItem("token")
    const _this = props
    const url = _this.uri
    let Data = _this.Data
    const ext = _this.ext
    const ExportAccount = async() => {
      let excel = true
      if (ext == ".txt") excel = false
      console.log(excel)
      let d = await Fetch(url, Data, 'GET', token, true, excel)
      download(d)
    }

    const download = (data) => {
        if (!data) {
            return
        }
        // const contentType = data.type
        // const fileName = contentType.split('filename=')[1]
        let url = window.URL.createObjectURL(new Blob([data]))
        if (ext == ".xlsx") url = window.URL.createObjectURL(new Blob([data], {type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}))
        let link = document.createElement('a')
        link.style.display = 'none'
        link.href = url
        link.id='Adownload'
        const date = new Date(),
            Y = date.getFullYear(),
            M = date.getMonth(),
            D = date.getDate(),
            h = date.getHours(),
            m = date.getMinutes(),
            s = date.getSeconds(),
            fileName = `${String(Y)}${String(M)}${String(D)}${String(h)}${String(m)}${String(s)}${ext}`
        // console.log(fileName)
        link.setAttribute('download', fileName)
        
        document.body.appendChild(link)
        link.click()
        document.getElementById('Adownload').remove();
    }
    return {
      ExportAccount
    }
  },
})
</script>