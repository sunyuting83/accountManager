<template>
  <div v-if="page.length > 0" class="mb-5">
    <hr />
    <nav class="pagination is-small" role="navigation" aria-label="pagination">
      <button class="button pagination-previous" @click="setPreviousPage" title="This is the first page" :disabled="current == 1?true:false">上一页</button>
      <button class="button pagination-next" @click="setNextPage" :disabled="current !== page.length && page.length > 1?false:true">下一页</button>
      <ul class="pagination-list">
        <li v-for="(item) in page" :key="item">
          <a class="pagination-link" @click="()=>{setPage(item)}" :class="item === current ? 'is-current': ''" :aria-label="'Page '+item" :aria-current="item === current ? 'page': ''">{{item}}</a>
        </li>
      </ul>
    </nav>
  </div>
</template>
<script>
import { reactive, toRefs, defineComponent } from 'vue'
export default defineComponent({
  name: 'PaginAtion',
  props: {
    total: Number,
    number: Number,
    GetData: Function
  },
  setup(props) {
    let states = reactive({
      current: 1,
      page: []
    })

    const makePage = (t) => {
      let x = []
      const p = Math.ceil(t/props.number)
      for (let i = 0; i < p; i++) {
        x = [...x, i + 1]
      }
      return x
    }
    states.page = makePage(props.total)

    const setPage = (p) =>{
      if (p !== states.current && p >= 1) {
        states.current = p
        props.GetData(p)
      }
    }
    const setNextPage = () =>{
      if (states.current !== states.page.length && states.current < states.page.length) {
        states.current = states.current + 1
        props.GetData(states.current)
      }
    }
    const setPreviousPage = () =>{
      if (states.current > 1) {
        states.current = states.current - 1
        props.GetData(states.current)
      }
    }

    return {
      ...toRefs(states),
      setPage,
      setNextPage,
      setPreviousPage
    }
  },
})
</script>
