
  // 函数：遍历所有 DOM 元素，找出包含文本的节点
  export  function findMatchingElements(textArray) {
    const elements = []; // 用于存储匹配的元素

    // 递归查找所有 DOM 节点
    function traverse(node) {
      // 如果节点是文本节点
      if (node.nodeType === Node.TEXT_NODE) {
        // 去除前后空白字符，并和 textArray 中的每一项进行匹配
        const text = node.textContent.trim();
        if (text && textArray.includes(text)) {
          elements.push(node.parentNode); // 找到匹配的文本，返回该文本节点的父元素
        }
      }

      // 遍历所有子节点
      node.childNodes.forEach(childNode => traverse(childNode));
    }

    // 从 body 开始遍历
    traverse(document.body);

    return elements;
  }
