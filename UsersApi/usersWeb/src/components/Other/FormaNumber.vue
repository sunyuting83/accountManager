<template>
  <span>{{FormaNumber(Numbers)}}</span>
</template>
<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'FormaNumber',
  props: {
    Numbers: {
      type: Number
    }
  },
  setup() {
    const getnum = (num) =>  {
      let result = ""
      const s = num.toString()
      if (s.indexOf(".") !== -1) {
        result = s.substring(0,s.indexOf(".")+2)
      }else{
        result = s
      }
      return result
    }

    const FormaNumber = (n) =>{
      let x = "0"
      if (n >= 100000000 && n !== 0) {
        let a = n / 100000000
        x = `${getnum(a)}亿`
      }else if (n === 123) {
        x = "识别错误"
      }else{
        if (n !== 0 ) {
          let a = n / 10000
          x = `${getnum(a)}万`
        }
      }
      return x
    }
    return {
      FormaNumber
    }
  },
})
</script>