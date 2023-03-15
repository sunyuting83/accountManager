<template>
  <button class="button is-small" :class="buttonLoading?'is-loading' + styles : styles" @click="ExportAccount">{{title}}</button>
</template>
<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'DownloadFile',
  props: {
    uri:{
      type: String,
      default: ""
    },
    status:{
      type: String,
      default: ""
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
    }
  },
  setup(props) {
    const _this = props
    const url = _this.uri
    const status = _this.status
    const ExportAccount = async() => {
      const d = await exportFile()
      download(d)
    }

    const exportFile = () => {
      const token = localStorage.getItem("token")

      let requestConfig = {
        method: "put",
        responseType: "blob"
      }
      Object.defineProperty(requestConfig, 'body', {
          value: JSON.stringify({
          status: status,
        })
      })
      requestConfig.headers = new Headers({
        Accept: '*/*',
      })
      requestConfig.headers.append("Content-Type","application/json;charset=UTF-8")
      requestConfig.headers.append('Authorization',`Bearer ${token}`)
      return new Promise((resolve) => {
        fetch(url, requestConfig)
          .then(res => {
            if(res.ok) {
              resolve(res.text())
            }else {
              resolve({
                status: 1,
                message: "访问出错"
              })
            }
          })
          .catch((err) => {
            resolve({
              status: 1,
              message: err.message
            })
          })
      })
    }
    const download = (data) => {
        if (!data) {
            return
        }
        // const contentType = data.type
        // const fileName = contentType.split('filename=')[1]
        let url = window.URL.createObjectURL(new Blob([data]))
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
            fileName = `${String(Y)}${String(M)}${String(D)}${String(h)}${String(m)}${String(s)}.txt`
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