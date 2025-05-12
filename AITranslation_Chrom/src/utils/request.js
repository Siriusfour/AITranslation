import OpenAI from "openai";
import { Prompts } from "./Prompt";

const getRequest = (url, data) => request('GET', url, data)
const postRequest = (url, data) => request('POST', url, data)
const putRequest = (url, data) => request('PUT', url, data)
const deleteRequest = (url, data) => request('DELETE', url, data)
const BaseUrl = "127.0.0.1:7950/AITranslation/"
const api_key=import.meta.env.VITE_API_KEY


//向目标接口发送请求
async function request(type, URL, data) {

  const response = await fetch(BaseUrl + URL, {
    method: type,
    body: typeof data === 'object' ? JSON.stringify(data) : data,
  }).then((res) => {
    console.log(res);

  }).catch((err) => {
    console.log(err);
  })
}

//不需要登陆的翻译接口
function AITranslation(){
fetch('https://open.bigmodel.cn/api/paas/v4/chat/completions', {
  method: 'POST',     
  headers: {
    'Content-Type': 'application/json',
    'Authorization': api_key
  },
  body: JSON.stringify({
    model: "glm-4-flash-250414",
    messages: [
      { "role": "system", "content":  },
      { "role": "user", "content": "请为我的产品创作一个吸引人的口号" },
    ]
  })
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
}


export { BaseUrl }
export {AITranslation }