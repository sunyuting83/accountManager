<template>
  <div :style="`background:url('${formState.background}'); background-size:100% 100%;width: 100%;height: 100%`">
    <a-row type="flex" class="pd-t30vh">
      <a-col :span="8"></a-col>
      <a-col :span="6">
        <a-form
          :model="formState"
          name="normal_login"
          class="login-form"
          @finish="onFinish"
          @finishFailed="onFinishFailed"
        >
          <a-form-item
            label="用户名"
            name="username"
            labelAlign="right"
            :rules="[{ required: true, message: '请输入用户名!', min: 3, max: 13 }]"
          >
            <a-input v-model:value="formState.username">
              <template #prefix>
                <UserOutlined class="site-form-item-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item
            label="密　码"
            name="password"
            labelAlign="right"
            :rules="[{ required: true, message: '请输入密码!', min: 5, max: 16 }]"
          >
            <a-input-password v-model:value="formState.password">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <a-form-item name="remember" no-style>
              <a-checkbox v-model:checked="formState.remember">记住我</a-checkbox>
            </a-form-item>
            <a class="login-form-forgot" href="">忘记密码？</a>
          </a-form-item>

          <a-form-item>
            <a-button :disabled="disabled" type="primary" html-type="submit" class="login-form-button">
              登录
            </a-button>
          </a-form-item>
          <a-form-item>
            Or
            <a href="">立即注册</a>
          </a-form-item>
        </a-form>
      </a-col>
    </a-row>
  </div>
</template>
<script lang="ts">
import { notification } from 'ant-design-vue'
import { defineComponent, reactive, computed, h } from 'vue'
import { UserOutlined, LockOutlined, InfoCircleOutlined } from '@ant-design/icons-vue'
import background from '../../assets/images/bbck.jpg'
import { Login } from '../../../wailsjs/go/main/App'
import { useRouter } from 'vue-router'
interface FormState {
  username: string;
  password: string;
  remember: boolean;
  background: string;
}
export default defineComponent({
  components: {
    UserOutlined,
    LockOutlined,
  },
  setup() {
    const router = useRouter()
    const formState = reactive<FormState>({
      username: '',
      password: '',
      remember: true,
      background: background,
    });
    const onFinish = async(values: any) => {
      const params = {
        username: values.username,
        password: values.password,
      }
      const data = await Login(params)
      if (data.status == 0) {
        router.push({
          'name': 'user',
        })
      }else {
        openNotification(data.message)
      }
    };

    const onFinishFailed = (errorInfo: any) => {
      openNotification("errorInfo")
    };
    const disabled = computed(() => {
      return !(formState.username && formState.password);
    });
    const openNotification = (text: string) => {
      notification.open({
        message: '发生错误',
        description:
        text,
        icon: () => h(InfoCircleOutlined, { style: 'color: #ff1855' }),
      });
    };
    return {
      formState,
      onFinish,
      onFinishFailed,
      disabled,
    };
  },
});
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