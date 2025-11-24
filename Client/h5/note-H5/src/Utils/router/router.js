import { createRouter, createWebHistory } from 'vue-router'
import Callback from '../../../component/callBack/callBack.vue'
import index from '../../../component/index/index.vue'


const routes = [
    // 根路径直接跳到 /index
    {
        path: '/',
        redirect: '/index'
    },
    {
        path: '/index',
        name: 'Index',
        component: index
    },
    {
        path: '/callback',
        name: 'Callback',
        component: Callback
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 👉 调试用：启动时在控制台打印所有路由，看是不是这份 router 在工作
console.log(
    '%c[router routes]',
    'color: #42b983',
    router.getRoutes().map(r => r.path)
)

export default router
