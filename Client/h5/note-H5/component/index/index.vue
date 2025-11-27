<script setup>
import 'remixicon/fonts/remixicon.css'
import 'ant-design-vue/dist/reset.css' // 路径按你项目实际调整
import {Login, LoginByGithub, LoginByWebAuthn} from './index.js'
import {onMounted} from "vue";

onMounted( () => {
  const token = localStorage.getItem('AccessToken')
  if (token) {
    Login()
  }
})

</script>

<template>
  <div id="page" class="site">
    <div class="container">
      <div class="login">
        <div class="hero">
          <h1>博康海洋生物有限公司<br>BoKang</h1>
          <p>如果你没有账号<br>可以<a href="#">点击这里</a>进行申请</p>
        </div>
        <div class="main">
          <form action="" >
            <p>
              <input placeholder="账号/account"  id="accountInput"  name="account" >
            </p>
            <p class="password">
              <input placeholder="密码/password" id="passwordInput"  name="password" type="password" >
              <i class="ri-eye-off-line"></i>

            </p>
            <p>
              <input type="button" class="submit" value="登录" id="loginBtn" @click="Login()">
            </p>
          </form>
          <div class="options">
            <div class="separator">
              <p>使用其他方式登录</p>
            </div>
            <ul>
              <li><a href="#"><i class="ri-wechat-fill ri-2x" ></i></a></li>
              <li><a href="#"><i class="ri-qq-fill ri-2x"></i></a></li>
              <li><a href="#"><i class="	ri-mail-fill ri-2x"  @click="LoginByGithub"></i></a></li>
              <li><a href="#"><i class="ri-key-fill ri-2x"     @click="LoginByWebAuthn"></i></a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>

<style scoped src="./index.css"></style>