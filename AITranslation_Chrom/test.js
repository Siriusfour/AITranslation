

fetch('https://open.bigmodel.cn/api/paas/v4/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization':'14cc2eb752714fba9b55a681793edfde.m0yOWQwvz8psxqZK'
    },
    body: JSON.stringify({ 
        model:"glm-4-flash-250414",
        messages:[
            {"role": "user", "content": "作为一名营销专家，请为我的产品创作一个吸引人的口号"},
            {"role": "assistant", "content": "当然，要创作一个吸引人的口号，请告诉我一些关于您产品的信息"},
            {"role": "user", "content": "智谱AI开放平台"},
            {"role": "assistant", "content": "点燃未来，智谱AI绘制无限，让创新触手可及！"},
            {"role": "user", "content": "创作一个更精准且吸引人的口号"}
        ]
     })
  })
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error('Error:', error));
  