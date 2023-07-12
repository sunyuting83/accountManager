<template>
  <a-layout-header style="background: #fff; padding: 0">
    <menu-unfold-outlined
      v-if="props.Trigger"
      class="trigger"
      @click="() => {if (typeof props.ChangeIt === 'function') {
  props.ChangeIt(); // 调用传递的函数
}}"
    />
    <menu-fold-outlined v-else class="trigger" @click="() => {if (typeof props.ChangeIt === 'function') {
  props.ChangeIt(); // 调用传递的函数
}}" />
  </a-layout-header>
</template>
<script lang="ts" setup>
import { defineProps } from 'vue';
import { MenuUnfoldOutlined, MenuFoldOutlined } from '@ant-design/icons-vue'
const props = defineProps({
  ChangeIt: Function,
  Trigger: Boolean,
})
</script>
<style>
</style>
