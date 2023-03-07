<template>
  <img :src="url" :class="imgStyle" />
</template>

<script>
import { reactive, toRefs, defineComponent } from 'vue'

export default defineComponent({
  name: 'EmptyEd',
  props: {
    imgUrl:{
      type: String,
      default: ""
    },
    imgStyle: {
      type: String,
      default: ""
    }
  },
  setup(props) {
    const _this = props
    const rootimg = 'https://r5k.oss-cn-hongkong.aliyuncs.com/book/z2_88_88407/img01.jpg'
    const defaultimg = `https://c2.im5i.com/2022/11/24/9vzRt.jpg`
    const loadingimg = `${rootimg}?x-oss-process=image/resize,w_15`
    let states = reactive({
      url: loadingimg
    })
    var newImg = new Image()
    newImg.src = _this.imgUrl
    newImg.onerror = () => { // 图片加载错误时的替换图片
      newImg.src = defaultimg
    }
    newImg.onload = async() => { // 图片加载成功后把地址给原来的img
      if(newImg.height < 1) {
        states.url = defaultimg
      }else {
        states.url = newImg.src
      }
      // this.url = await Api.Img2Base64(this.imgUrl, this.bookId)
    }
    return {
      ...toRefs(states)
    }
  },
})
</script>