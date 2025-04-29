import { findMatchingElements } from "../util/domHandle.js" 
import RequestAPI, { url } from "../util/request.js"
import { createApp } from 'vue'
const app = createApp({
  data() {
    return {
      count: 0
    }
  }
})
console.log("==========================",app);


app.mount('#app')


// 获取按钮实例
let changeColor = document.getElementById('changeColor');


// 点击按钮
changeColor.addEventListener('click', async () => {

  console.log("==========================",app);


  let [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

  // 向目标页面里注入js方法
  chrome.scripting.executeScript({
    target: { tabId: tab.id },
    function: setPageBackgroundColor
  });


});


// 注入的方法
function setPageBackgroundColor() {
  const allElements = document.body.querySelectorAll('*');
  const textContent = document.body.innerText;
  console.log(textContent);
  const textArray = textContent.split('\n').map(line => line.trim()).filter((line) => {

    return line.length > 0 && line !== "";
  });

  let matchingElements = []


  //筛选出所有中文
  const chineseRegex = /[\u4e00-\u9fa5]/;
  let textFilterArray = textArray.filter((text) => { return !chineseRegex.test(text) });
  let finalTextArray = textFilterArray.filter((text) => {
    console.log(text);

    return text.length > 0 && text !== ""
  })
  console.log(finalTextArray);

  //  遍历allElements,找到传入字符串匹配的1元素
  function findElementsWithText(text, allElements) {
    return Array.from(allElements).find((item) => {
      return item.textContent.trim() === text;
    });
  }

  //遍历所有
  textFilterArray.forEach((item) => {
    matchingElements.push(findElementsWithText(item, allElements))
  })

   const final =  matchingElements.filter((item)=>{return item!==undefined})

  // 输出结果
  console.log(final);


  //2=======================文本数组发送到服务器，获取response
  RequestAPI.post(url,"123")



  //把返回值插入到Dom页面

};







