const url = "http://127.0.0.1:7950" // 添加http://协议

const getRequest = (url, data) => request('GET', url, data)
const postRequest = (url, data) => request('POST', url, data) // 修正拼写错误
const putRequest = (url, data) => request('PUT', url, data)
const deleteRequest = (url, data) => request('DELETE', url, data)

function request(type, URL, data) {
  // 返回Promise以便调用者可以使用.then()
  return fetch(URL, {
    method: type,
    headers: {
      'Content-Type': 'application/json',
    },
    body: typeof data === 'object' ? JSON.stringify(data) : data,
  })
    .then(response => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return response.json() // 解析JSON响应
    })
    .then(data => {
      console.log('Success:', data)
      return data // 返回数据以便链式调用
    })
    .catch(error => {
      console.error('Error:', error)
      throw error // 重新抛出错误以便调用者捕获
    })
}

const RequestAPI = {
  get: getRequest,
  post: postRequest,
  put: putRequest,
  delete: deleteRequest
}

export default RequestAPI
export { url }