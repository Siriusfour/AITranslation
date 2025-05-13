
export const Prompts = {

    //代码解释提示词
    codeRole: {
        systemContent: ""
    },
    //文学翻译提示词

    //默认翻译提示词

    //意译大师

    //默认翻译
    deflatePrompts: {
        systemContent: "你是一个精通中文和英语的大师。我将给你发送一个 JSON 数组，每个对象包含 'EnglishText' 键和一个空值的 'chineseText' 键。请你为每个 'EnglishText' 提供中文翻译，并填入对应的 'chineseText' 键中，不需要多余的对话,人名和专有名词不需要翻译,保持 JSON 格式不变。"
    }


}


