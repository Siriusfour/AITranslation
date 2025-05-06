import RequestAPI from 'request'


export function setPageBackgroundColor() {
  console.log("===========================================");
  
  const allElements = document.body.querySelectorAll('*');
  const textContent = document.body.innerText;
console.log(textContent);

//分割textContent，形成数组
  const textArray = textContent.split('\n').map(line => line.trim()).filter((line) => {
    return line.length > 0 && line !== "";
  });


  let matchingElements = []


  //筛选出所有中文，空字符
  const chineseRegex = /[\u4e00-\u9fa5]/;
  let textFilterArray = textArray.filter((text) => { return !chineseRegex.test(text) });
  let finalTextArray = textFilterArray.filter((text) => {
    return text.length > 0 && text !== ""
  })
  console.log(finalTextArray);

  //  遍历allElements,找到传入字符串匹配的元素
  function findElementsWithText(text, allElements) {
    return Array.from(allElements).find((item) => {
      return item.textContent.trim() === text;
    });
  }

  //遍历所有
  textFilterArray.forEach((item) => {
    matchingElements.push(findElementsWithText(item, allElements))
  })

  console.log(matchingElements);
  
  const final = matchingElements.filter((item) => { return item !== undefined })

  // 输出结果
  console.log(final);


  //2=======================文本数组发送到服务器，获取response
  RequestAPI.post()
  


  //把返回值插入到Dom页面

};


