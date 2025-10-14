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
async function AITranslation(finalTextArrayJson) {
  try {
    const textArray = JSON.stringify(finalTextArrayJson);
    console.log('Input to AITranslation:', textArray);

    const response = await fetch('https://open.bigmodel.cn/api/paas/v4/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': api_key
      },
      body: JSON.stringify({
        model: "glm-4-flash-250414",
        messages: [
          { role: "system", content: Prompts.deflatePrompts.systemContent },
          { role: "user", content: textArray }
        ]
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();


    if (!data.choices || !data.choices[0] || !data.choices[0].message) {
      throw new Error('Invalid API response: choices or message missing');
    }

  

    const result = data.choices[0].message.content;

    return result;
   
  } catch (error) {
    console.log("err")
    console.error('AITranslation error:', error);
    throw error;
  }
}


export { BaseUrl }
export {AITranslation }