<script setup>
import { useRoute } from 'vue-router'
import { watchEffect } from 'vue'
const route = useRoute()

// 调试：实时打印当前路由
watchEffect(() => {
  console.log('🔥 current route =>', route.fullPath)
})
</script>

<template>
  <RouterView />
</template>
