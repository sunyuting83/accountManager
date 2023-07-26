<template>
  <div :style="`background:url('${background}'); background-size:100% 100%;width: 100%;height: 100%`">
    <a-row type="flex" class="pd-t30vh">
      <a-col :span="6"></a-col>
      <a-col :span="12">
        <a-form
          ref="formRef"
          name="custom-validation"
          :model="formState"
          :rules="rules"
          v-bind="layout"
          @finish="handleFinish"
        >
          <a-form-item has-feedback label="用户名" name="username">
            <a-input v-model:value="formState.username" />
          </a-form-item>
          <a-form-item has-feedback label="密码" name="password">
            <a-input v-model:value="formState.password" type="password" autocomplete="off" />
          </a-form-item>
          <a-form-item has-feedback label="重复密码" name="repassword">
            <a-input v-model:value="formState.repassword" type="password" autocomplete="off" />
          </a-form-item>
          <a-form-item has-feedback label="验证码" name="vcode">
            <a-input v-model:value="formState.vcode" />
            <img :src="vcodeImg" @click="getCaptcha" />
          </a-form-item>
          <a-form-item label="推荐人" name="referrer">
            <a-input v-model:value="referrer" />
          </a-form-item>
          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-button type="primary" :disabled="disabled" html-type="submit">立即注册</a-button>
            <a-button style="margin-left: 10px" @click="resetForm">重置表单</a-button>
          </a-form-item>
          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            已有帐号？ <a-button type="link" @click="pushLogin">立即登陆</a-button>
          </a-form-item>
        </a-form>
      </a-col>
    </a-row>
  </div>
</template>
<script lang="ts" setup>
import { reactive, ref, computed, h } from 'vue';
import type { Rule } from 'ant-design-vue/es/form';
import { Captcha } from '../../../wailsjs/go/main/App'
import { ClipboardGetText } from '../../../wailsjs/runtime/runtime'
import background from '../../assets/images/bbck.jpg'
import { useRouter } from 'vue-router'
import { Regedit } from '../../../wailsjs/go/main/App'
import { InfoCircleOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { notification } from 'ant-design-vue'

const router = useRouter()
interface FormState {
  password: string;
  repassword: string;
  username: string;
  vcode: string;
  referrer: string;
}
const formRef = ref<any>()
const formState = reactive<FormState>({
  password: '',
  repassword: '',
  username: '',
  vcode: '',
  referrer: '',
})

ClipboardGetText().then((e: string) => {
  if (e !== '') {
    if (e.indexOf('&referrer=') !== -1) {
      const refStr = e.split('&referrer=')[1]
      referrer.value = refStr.substring(0,16)
    }
  }
})

const disabled = computed(() => {
  return !(formState.username.length >= 3 && formState.username.length <= 12 && formState.password && formState.password == formState.repassword)
})

let checkUsername = async (_rule: Rule, value: string) => {
  if (!value) {
    return Promise.reject('请输入用户名');
  }
  if (value.length < 3 || value.length > 12) {
    return Promise.reject('用户名必须大于3位或小于12位');
  } else {
    return Promise.resolve();
  }
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
};
let validatePass2 = async (_rule: Rule, value: string) => {
  if (value === '') {
    return Promise.reject('请再次输入密码');
  } else if (value !== formState.password) {
    return Promise.reject("两次密码不同!");
  } else {
    return Promise.resolve();
  }
}

const checkVcode = async (_rule: Rule, value: string) => {
  if (value === '') {
    return Promise.reject('请输入验证码')
  }if (value.length !== 6) {
    return Promise.reject('验证码必须等于6位');
  } else {
    return Promise.resolve();
  }
}

const rules: Record<string, Rule[]> = {
  password: [{ required: true, validator: validatePass, trigger: 'change' }],
  repassword: [{ validator: validatePass2, trigger: 'change' }],
  username: [{ validator: checkUsername, trigger: 'change' }],
  vcode: [{ validator: checkVcode, trigger: 'change' }],
};
const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 14 },
};
const handleFinish = async(values: FormState) => {
  const data = await Regedit(values)
  if (data.status == 0) {
    sucNotification("注册成功，请登陆")
    pushLogin()
  }else{
    errNotification(data.message)
  }
}

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

const vcodeImg = ref<string>()
const referrer = ref<string>()

const getCaptcha = async() => {
  const data = await Captcha()
  vcodeImg.value = `data:image/png;base64,${data}`
}

const resetForm = () => {
  formRef.value.resetFields();
}

const pushLogin = () => {
  router.push({
    'name': 'login',
  })
}
getCaptcha()
</script>
<style>
.login-form {
  max-width: 300px;
}
.login-form-forgot {
  float: right;
}
.login-form-button {
  width: 100%;
}
.pd-t30vh {
  padding-top: 30vh
}
</style>