import { createApp } from 'vue'
import route from './router/index'
import App from './App.vue'

createApp(App).use(route).mount('#app')
