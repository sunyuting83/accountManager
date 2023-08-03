<template>
  <span>{{PriceCalc(Price,Calc)}}</span>
</template>
<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'PriceCalc',
  props: {
    Price: {
      type: Number
    },
    Calc: {
      type: Object
    }
  },
  setup() {
    const PriceCalc = (price, calc) => {
      if (price <= 0) {
        return 0
      }else {
        const {BasePrice,UnitPrice,SingleNumber} = calc
        // console.log(UnitPrice / (SingleNumber*100000000))
        // console.log(price)
        const newPrice = BasePrice + ((UnitPrice / (SingleNumber*100000000)) * price)
        return newPrice.toFixed(2)
        // return 0.1
      }
    }
    return {
      PriceCalc
    }
  },
})
</script>