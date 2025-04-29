// 定义颜色
let color = '#3aa757';

// 首次安装插件、插件更新、chrome浏览器更新时触发
chrome.runtime.onInstalled.addListener(() => {
  chrome.storage.sync.set({ color });

  
});