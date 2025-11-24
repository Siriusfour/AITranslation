
import { createApp } from 'vue'
import Antd from 'ant-design-vue';
import App from '/component/app/app.vue'
import 'ant-design-vue/dist/reset.css';
import router from "./src/Utils/router/router.js";

createApp(App).use(router).use(Antd).mount('#app')