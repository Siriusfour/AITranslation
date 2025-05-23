import { createApp } from 'vue'
import './style.css'
import Antd from 'ant-design-vue';
import App from './App.vue';
import 'ant-design-vue/dist/reset.css';
import {router} from './components/router'

console.log(import.meta.env.VITE_API_KEY);


createApp(App).use(Antd).use(router).mount('#app')


