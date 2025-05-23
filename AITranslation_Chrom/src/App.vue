<script setup>
import { MailOutlined, AppstoreOutlined, SettingOutlined } from '@ant-design/icons-vue';
import { h, ref } from 'vue';
import { router } from './components/router';

const current = ref(['mail']);
const items = ref([
  {
    key: 'Translation',
    keyPath: '/',
    icon: () => h(MailOutlined),
    label: '翻译',
    title: 'Translation',

  },
  {
    key: 'Chat',
    keyPath: 'talk',
    icon: () => h(AppstoreOutlined),
    label: '对话',
    title: 'Chat',
  },
  {
    key: 'more',

    icon: () => h(AppstoreOutlined),
    label: '更多',
    title: 'Chat',
    children: [
      {
        type: 'group',
        label: '生活',
        children: [
          {
            label: '旅游攻略',
            key: 'setting:1',
          },
          {
            label: 'Option 2',
            key: 'setting:2',
          },
        ],
      },
      {
        type: 'group',
        label: '科技',
        children: [
          {
            label: '一键爬虫',
            key: 'setting:3',
          },
          {
            label: 'Option 4',
            key: 'setting:4',
          },
        ],
      },
    ],
  },
  {
    key: 'My',
    keyPath: '/Mine',
    icon: () => h(SettingOutlined),
    label: '我的',
    title: '',
  },

]);

const changeComponent = ({ item }) => {
  console.log(item.keyPath);
  router.push(item.keyPath);
};


</script>

<template>
  <div style="margin-top: -20px; ">
    <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" @click="changeComponent" />
  </div>

  <RouterView />
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}

.Translation_btn {
  width: 50%;
  height: 50px;
  color: crimson;
}
</style>
