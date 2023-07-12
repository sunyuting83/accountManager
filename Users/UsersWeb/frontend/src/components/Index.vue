<template>
  <div>
    <a-row type="flex" class="pd-t30vh">
      <a-col :span="8"></a-col>
      <a-col :span="6">
        <a-space>
          <a-spin size="large" />
        </a-space>
      </a-col>
    </a-row>
  </div>
</template>
<script lang="ts" setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { CheckLogin } from '../../wailsjs/go/main/App'
const router = useRouter()
onMounted(async() => {
  const data = await CheckLogin()
  // console.log(data)
  if (data === null || data.status === 1) {
    router.push("login")
  }else{
    router.push("main/user")
  }
})
</script>

<style scoped>
.pd-t30vh {
  padding-top: 30vh;
  text-align: center;
}
</style>