<script setup>
import {  ref } from 'vue';
import { DownOutlined } from '@ant-design/icons-vue';

const modeList = [
  { modeName: 'DeepSeek', Link: "", img: "deepseek-color.svg" }, { modeName: 'chatGPT=4.0', Link: "openai.svg", img: "openai.svg" }, { modeName: '智谱GLM-4-Flash', Link: "", img: "/chatglm.svg" }
]
const modeImg = ref("/chatglm.svg")
const modeName = ref('智谱GLM-4-Flash')
const value4 = ref('');
const visible = ref(true);
const status = ref("翻译页面")


function changeStatus() {
  status.value == "翻译页面" ? status.value = "显示原文" : status.value = "翻译页面"
  console.log(status.value);
}


function changeMode({ key }) {
  console.log(`Click on item ${key}`);
  console.log(typeof key);
  modeName.value = modeList[key].modeName
  modeImg.value = modeList[key].img

}
// 点击按钮触发
async function greet() {
  changeStatus()

  let [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

  // 向目标页面里注入js方法
  chrome.scripting.executeScript({
    target: { tabId: tab.id },
    files: ['/assets/AITranslation.js']
  });
}


</script>
<template>

    <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" />


  <div style="display: flex; flex-direction: column; align-items: center; justify-content: center;">

    <a href="https://vite.dev" target="_blank">
      <img :src="modeImg" class="logo" alt="Vite logo" />
    </a>
    <div style="margin-bottom: 30px;">
      <a-dropdown overlayStyle="margin-top:50rpx">
        <a class="ant-dropdown-link" @click.prevent style="color: black;">
          {{ modeName }}
          <DownOutlined />
        </a>
        <template #overlay>
          <a-menu @click="changeMode">
            <a-menu-item v-for="(item, index) in modeList" :key="index">{{ item.modeName }}</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>


    <a-space>
      <a-input-password v-model:value="value4" v-model:visible="visible" placeholder="input password" />
      <a-button @click="visible = !visible">{{ visible ? '翻译' : '还原' }}</a-button>
    </a-space>

    <button @click="greet" style="margin-top: 30px;">{{ status }}</button>
  </div>
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