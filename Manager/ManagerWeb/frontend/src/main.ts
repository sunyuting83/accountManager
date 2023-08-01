import { createApp } from 'vue'
import route from './router/index'
import Antd from 'ant-design-vue'
import App from './App.vue'
import "./style.css"
if (!('go' in window)) location.replace('/')

const app = createApp(App).use(route)

app.use(Antd).mount('#app')