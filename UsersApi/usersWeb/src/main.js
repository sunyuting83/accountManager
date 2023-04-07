import { createApp } from 'vue'
import route from './router/index'
import App from './App.vue'
import tooltip from '@/directives/tooltip' 
import "@/directives/tooltip.css";

createApp(App).use(route).directive("tooltip",tooltip).mount('#app')
