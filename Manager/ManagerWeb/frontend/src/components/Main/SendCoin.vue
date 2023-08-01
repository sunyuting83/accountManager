<template>
  <a-layout-content :style="{background: '#fff' }">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}">
      <div class="content">
        
      </div>
    </div>
  </a-layout-content>
</template>
<script lang="ts" setup>
import PageHeader from './PageHeader.vue'
import { ref, h } from 'vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { Empty } from 'ant-design-vue';
import { useRouter } from 'vue-router'
const router = useRouter()


interface PageHeaderData {
  title: string;
  subtitle: string;
  routers: Array<
    {
      path: string;
      breadcrumbName: string;
    }
  >
}

const pageHeader = ref<PageHeaderData>({
  title: "发币给用户",
  subtitle: '发币给用户详情',
  routers: [
    {
      path:'main',
      breadcrumbName: '发币给用户'
    },
    {
      path:'main',
      breadcrumbName: '发币给用户'
    }
  ],
})


const errNotification = (text: string) => {
  notification.open({
    message: "发生错误",
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}

interface State {
  loading: boolean;
}
const state = ref<State>({
  loading: false,
})

</script>
