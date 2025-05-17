import { RequestAPI } from '../utils/request'
import { AITranslation } from '../utils/request'
import { handleJsonFromString } from '../utils/handle'


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
function findElementsWithText(text, allElements) {
  return Array.from(allElements).find((item) => {
    return item.textContent.trim() === text;
  });
}

(async () => {
  // 获取翻译结果
  console.log("ok!");
  const resultJson = await AITranslation(finalTextArrayJson)


  const handleJson = handleJsonFromString(resultJson)
  console.log(handleJson);
  

  try {
    var result = JSON.parse(handleJson)
    console.log(result);
  } catch (err) {
    console.log(err);
  }

console.log(result);
  result.forEach((item) => {
    console.log(item.EnglishText);
    try{
      const matches = findElementsWithText(item.EnglishText, allElements);
      matches.textContent += item.chineseText;
      matchingElements.push(matches)
    }catch(err){
console.log(err);
    }
  })

  console.log(result);

})()












