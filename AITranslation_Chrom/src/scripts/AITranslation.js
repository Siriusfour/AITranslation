import { RequestAPI } from '../utils/request'
import { AITranslation } from '../utils/request'


console.log("===========================================");

const allElements = document.body.querySelectorAll('*');
const textContent = document.body.innerText;


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


const finalTextArrayJson = []
finalTextArray.forEach((item) => {
  finalTextArrayJson.push({ "EnglishText": item, "chineseText": "" })
})


//  遍历allElements,找到传入字符串匹配的元素
function findElementsWithText(text) {
  // 使用 XPath 查询包含完全匹配文本的文本节点
  const xpath = `//text()[.='${text.replace(/'/g, "\\'")}']`;
  const result = document.evaluate(
    xpath,
    document,
    null,
    XPathResult.ORDERED_NODE_SNAPSHOT_TYPE,
    null
  );

  // 存储匹配的 DOM 元素
  const elements = [];
  for (let i = 0; i < result.snapshotLength; i++) {
    const textNode = result.snapshotItem(i);
    // 返回文本节点的父元素（通常是我们想要的 DOM 元素）
    elements.push(textNode.parentNode);
  }

  return elements;
}

(async () => {
  // 获取翻译结果
  console.log("ok!");
  const resultJson = await AITranslation(finalTextArrayJson)
  const result = JSON.parse(resultJson)

  console.log(result);
  
  
  result.forEach((item) => {
    console.log(item.EnglishText);
    matchingElements.push(findElementsWithText(item.EnglishText))
  })

console.log('1');

  console.log(matchingElements,"ok!");

})()












