<template>
  <a-layout-content :style="{background: '#fff' }" v-if="userState.status == 0">
    <PageHeader :data="pageHeader" />
    <div :style="{ padding: '24px'}">
      <div class="content">
        <div class="main">
          <a-descriptions size="small" :column="2">
            <a-descriptions-item label="用户名">{{userState.users.UserName}}</a-descriptions-item>
            <a-descriptions-item label="钱包地址">{{userState.users.WalletAddress}}</a-descriptions-item>
            <a-descriptions-item label="E-mail">{{userState.users.Email}}</a-descriptions-item>
            <a-descriptions-item label="电话">{{userState.users.PhoneNumber}}</a-descriptions-item>
            <a-descriptions-item label="最后登陆时间">{{foramTime(userState.users.UpdatedAt)}}</a-descriptions-item>
            <a-descriptions-item label="最后登陆IP">{{userState.users.IPAddress}}</a-descriptions-item>
            <a-descriptions-item label="所在地">
              {{userState.users.LocalAddress}}
            </a-descriptions-item>
          </a-descriptions>
        </div>
        <div class="extra">
          <div
            :style="{
              display: 'flex',
              width: 'max-content',
              justifyContent: 'flex-end',
            }"
          >
            <a-statistic
              title="状态"
              :value="userState.users.NewStatus == 0 ? '正常' : '锁定'"
              :style="{
                marginRight: '32px',
              }"
            />
            <a-statistic title="余额" prefix="¥" :value="userState.users.Coin" />
          </div>
        </div>
      </div>
      <a-descriptions title="" :style="{'margin-top': '1rem'}">
        <a-descriptions-item label="">
          <a-space>
            <a-button type="primary" @click="showChangePassword">修改密码</a-button>
            <a-button type="primary" @click="showTransfer">转账给Ta人</a-button>
          </a-space>
        </a-descriptions-item>
      </a-descriptions>
    </div>
  </a-layout-content>
  <a-modal v-model:visible="openChangePassword" title="修改密码" @ok="changePassword" :maskClosable="false">
    <template #footer>
      <a-button key="back" @click="handleCancel">取消修改</a-button>
      <a-button key="submit" type="primary" :loading="loading" :disabled="disabled" @click="changePassword">提交修改</a-button>
    </template>
    <a-form
      ref="formRef"
      name="custom-validation"
      :model="formState"
      :rules="rules"
      v-bind="layout"
    >
      <a-form-item has-feedback label="新密码" name="password">
        <a-input v-model:value="formState.password" type="password" autocomplete="off" />
      </a-form-item>
      <a-form-item has-feedback label="重复密码" name="repassword">
        <a-input v-model:value="formState.repassword" type="password" autocomplete="off" />
      </a-form-item>
    </a-form>
  </a-modal>
  <a-modal v-model:visible="openTransfer" title="转账给Ta人" @ok="Transfer" :maskClosable="false">
    <template #footer>
      <a-button key="back" @click="handleCancelT">取消</a-button>
      <a-button key="submit" type="primary" :loading="loading" :disabled="disabledt" @click="Transfer">转账</a-button>
    </template>
    <a-form
      ref="TransferRef"
      name="custom-validation"
      :model="TransferState"
      :rules="trules"
      v-bind="layout"
    >
      <a-form-item has-feedback label="钱包地址" name="wallet">
        <a-input v-model:value="TransferState.wallet" autocomplete="off" />
      </a-form-item>
      <a-form-item has-feedback label="金额" name="coin">
        <a-input v-model:value="TransferState.coin" type="number" autocomplete="off" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
<script lang="ts" setup>
import PageHeader from './PageHeader.vue'
import { onMounted, ref, h, reactive, computed } from 'vue'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'
import { GetUsers, RePassword, TransferUseWallet } from '../../../wailsjs/go/main/App'
import type { Rule } from 'ant-design-vue/es/form';

import { useRouter } from 'vue-router'
const router = useRouter()

interface FormState {
  password: string;
  repassword: string;
}
const formRef = ref<any>()
const formState = reactive<FormState>({
  password: '',
  repassword: '',
})

interface TransferState {
  coin: number;
  wallet: string;
}
const TransferRef = ref<any>()
const TransferState = reactive<TransferState>({
  coin: 0.0,
  wallet: '',
})

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
  title: "用户中心",
  subtitle: '用户详细信息',
  routers: [
    {
      path:'main',
      breadcrumbName: '用户中心'
    },
    {
      path:'main',
      breadcrumbName: '用户中心'
    }
  ],
});

interface Users {
  Coin: number;
  CreatedAt: number;
  Email: string;
  ID: number;
  IPAddress: string;
  LocalAddress: string;
  NewStatus: number;
  Password: string;
  PhoneNumber: string;
  UpdatedAt: number;
  UserName: string;
  WalletAddress: string;
}

interface User {
  status: number;
  message: string;
  users: Users;
}

const userState = ref<User>({
  status: 0,
  message: '',
  users: {
    Coin: 0,
    CreatedAt: 0,
    Email: '',
    ID: 0,
    IPAddress: '',
    LocalAddress: '',
    NewStatus: 0,
    Password: '',
    PhoneNumber: '',
    UpdatedAt: 0,
    UserName: '',
    WalletAddress: '',
  },
});
onMounted(async() => {
  const data = await GetUsers()
  if (data.status == 0) {
    userState.value = data as User
    // console.log(userState.value)
  }else {
    openNotification(data.message)
    if (data.message == "403") {
      router.push({
        'name': 'login',
      })
    }
  }
})
const openNotification = (text: string) => {
  notification.open({
    message: '发生错误',
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}
const foramTime = (d: number) => {
  const date = new Date(d)
  const Y = date.getFullYear()
  const M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1)
  const D = date.getDate()
  const H = date.getHours() < 10 ? `0${date.getHours()}` : date.getHours()
  const minute = date.getMinutes() < 10 ? `0${date.getMinutes()}` : date.getMinutes()
  const S = date.getSeconds() < 10 ? `0${date.getSeconds()}` : date.getSeconds()
  return `${Y}-${M}-${D} ${H}:${minute}:${S}`
}
const loading = ref<boolean>(false);
const openChangePassword = ref<boolean>(false)
const openTransfer = ref<boolean>(false)

const showChangePassword = () => {
  openChangePassword.value = true;
}
const showTransfer = () => {
  openTransfer.value = true;
}

const changePassword = async() => {
  loading.value = true
  const data = await RePassword(formState)
  if (data.status == 0) {
    loading.value = false
    openChangePassword.value = false
    formState.password = ""
    formState.repassword = ""
    sucNotification(data.message)
  }else{
    loading.value = false
    openChangePassword.value = false
    formState.password = ""
    formState.repassword = ""
    errNotification(data.message)
  }
}

const Transfer = async() => {
  loading.value = true
  // console.log(TransferState)
  const data = await TransferUseWallet(TransferState)
  if (data.status == 0) {
    loading.value = false
    openTransfer.value = false
    TransferState.coin = 0.0
    TransferState.wallet = ""
    sucNotification(data.message)
  }else{
    loading.value = false
    openTransfer.value = false
    TransferState.coin = 0.0
    TransferState.wallet = ""
    errNotification(data.message)
  }
}

const handleCancel = () => {
  openChangePassword.value = false
  formState.password = ""
  formState.repassword = ""
};
const handleCancelT = () => {
  openTransfer.value = false
  TransferState.wallet = ""
};


let validatePass = async (_rule: Rule, value: string) => {
  if (value === '') {
    return Promise.reject('请输入密码');
  } else {
    if (value.length < 6 || value.length > 16) {
      return Promise.reject('密码必须大于6位或小于16位');
    } else {
      if (formState.repassword !== '') {
        formRef.value.validateFields('repassword');
      }
      return Promise.resolve();
    }
  }
}
let validatePass2 = async (_rule: Rule, value: string) => {
  if (value === '') {
    return Promise.reject('请再次输入密码');
  } else if (value !== formState.password) {
    return Promise.reject("两次密码不同!");
  } else {
    return Promise.resolve();
  }
}

let checkWallet = async (_rule: Rule, value: string) => {
  if (!value) {
    return Promise.reject('请输入对方钱包地址');
  }
  if (value.length != 16) {
    return Promise.reject('钱包地址必须等于16位');
  } else {
    return Promise.resolve();
  }
}
let checkCoin = async (_rule: Rule, value: number) => {
  if (!value) {
    return Promise.reject('金额不能为空');
  }
  if (value <= 0) {
    return Promise.reject('金额不能<=0');
  } else {
    return Promise.resolve();
  }
}

const rules: Record<string, Rule[]> = {
  password: [{ required: true, validator: validatePass, trigger: 'change' }],
  repassword: [{ validator: validatePass2, trigger: 'change' }],
}

const trules: Record<string, Rule[]> = {
  wallet: [{ required: true, validator: checkWallet, trigger: 'change' }],
  coin: [{ required: true, validator: checkCoin, trigger: 'change' }],
}
const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 14 },
}

const disabled = computed(() => {
  return !(formState.password && formState.password == formState.repassword)
})

const disabledt = computed(() => {
  return !(TransferState.wallet && TransferState.coin && TransferState.wallet.length == 16 && TransferState.coin > 0)
})


const errNotification = (text: string) => {
  notification.open({
    message: "发生错误",
    description:
    text,
    icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
  });
}

const sucNotification = (text: string) => {
  notification.open({
    message: "成功",
    description:
    text,
    icon: () => h(CheckCircleOutlined, { style: 'color: #389e0d' }),
  });
}

</script>
<style scoped>
.content {
  display: flex;
}
</style>
