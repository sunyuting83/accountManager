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
          @validate="handleValidate"
          @finishFailed="handleFinishFailed"
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
          <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-button type="primary" html-type="submit">立即注册</a-button>
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
import { reactive, ref } from 'vue';
import type { Rule } from 'ant-design-vue/es/form';
import type { FormInstance } from 'ant-design-vue';
import background from '../../assets/images/bbck.jpg'
import { useRouter } from 'vue-router'
const router = useRouter()
interface FormState {
  password: string;
  repassword: string;
  username: string;
}
const formRef = ref<FormInstance>();
const formState = reactive<FormState>({
  password: '',
  repassword: '',
  username: '',
});
let checkAge = async (_rule: Rule, value: string) => {
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
    if (formState.repassword !== '') {
      formRef.value.validateFields('repassword');
    }
    return Promise.resolve();
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
};

const rules: Record<string, Rule[]> = {
  password: [{ required: true, validator: validatePass, trigger: 'change' }],
  repassword: [{ validator: validatePass2, trigger: 'change' }],
  username: [{ validator: checkAge, trigger: 'change' }],
};
const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 14 },
};
const handleFinish = (values: FormState) => {
  console.log(values, formState);
};
const handleFinishFailed = errors => {
  console.log(errors);
};
const resetForm = () => {
  formRef.value.resetFields();
};
const handleValidate = (...args) => {
  console.log(args);
}

const pushLogin = () => {
  router.push({
    'name': 'login',
  })
}
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