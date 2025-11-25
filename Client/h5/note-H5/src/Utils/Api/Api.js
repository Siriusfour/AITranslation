
import ApiCode from "./ApiCode.js";


const baseUrl = '/noteApi';


const refreshRequest = async (RefreshToken) => {

  const url = baseUrl + 'user/refresh_token';

  try {
    const response = await fetch(url, {
      method: "POST",
      body: RefreshToken || '',
      mode: "cors",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json",
        "Authorization": AccessToken || ''
      },
    });


    const result = await response.json();

    const {code} = result

    if (code === ApiCode.SUCCESS) {
      const AccessToken = result.data.accessToken; // 更新token
      localStorage.setItem("AccessToken", AccessToken);
      return AccessToken;
    }

    return null;
  } catch (error) {
    console.error('Token refresh failed:', error);
    return null;
  }
};

const request = async (url, method, data) => {

  console.log(url, method, data);

  let AccessToken = localStorage.getItem("AccessToken");
  let RefreshToken = localStorage.getItem("RefreshToken");
  const fullUrl = baseUrl + url;

  // 构建请求配置
  const requestConfig = {
    method: method,
    mode: "cors",
    credentials: "same-origin",
    headers: {
      "Authorization": AccessToken || ''
    }
  };

  // 处理URL和请求体
  let requestUrl = fullUrl;
  if (method === 'GET') {
    // GET请求：将参数拼接到URL上，不设置body和Content-Type
    if (data && Object.keys(data).length > 0) {
      const params = new URLSearchParams(data);
      requestUrl += `?${params}`;
    }
  } else {
    // 非GET请求：设置Content-Type和body
    requestConfig.headers["Content-Type"] = "application/json";
    if (data) {
      requestConfig.body = JSON.stringify(data);
    }
  }

  console.log(requestConfig);

  try {
    const response = await fetch(requestUrl, requestConfig);
    const result = await response.json();
    const { code } = result;

    // Token过期处理
    if (code === ApiCode.TOKEN_EXPIRED) {
      const newAccessToken = await refreshRequest(RefreshToken);
      if (newAccessToken) {
        // 更新请求配置中的token
        requestConfig.headers["Authorization"] = newAccessToken;

        const newResponse = await fetch(requestUrl, requestConfig);
        return await newResponse.json();
      } else {
        localStorage.removeItem("AccessToken");
        localStorage.removeItem("RefreshToken");
        throw new Error('Authentication failed');
      }

    }else if(code !== ApiCode.SUCCESS){

      throw new Error(result.message || '请求失败，服务器返回非成功状态码');

    }else{
      console.log(result);
    }


    return result.data;
  } catch (error) {

    throw error

  }
};


const getRequest = (url, data) => request(url, 'GET', data);
const postRequest = (url, data) => request(url, 'POST', data);
const putRequest = (url, data) => request(url, 'PUT', data);
const deleteRequest = (url, data) => request(url, 'DELETE', data);

const Api = {
  get: getRequest,
  post: postRequest,
  put: putRequest,
  delete: deleteRequest
};

export default Api;
