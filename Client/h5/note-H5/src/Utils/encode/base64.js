
//Base64URL-->字节数组 解码函数
function Base64URL_To_arrayBuffer(base64url) {
    if (typeof base64url !== 'string') {
        throw new TypeError('Base64URL_To_arrayBuffer输入必须是数组');
    }
    // 1. 替换 URL 安全字符
    let base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');

    // 2. 补齐填充（Base64 长度必须是 4 的倍数）
    const padLength = (4 - (base64.length % 4)) % 4;
    base64 += '='.repeat(padLength);

    // 3. 用 atob 解码为二进制字符串
    const binaryString = atob(base64);

    // 4. 转成 Uint8Array（字节序列）
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
}


// 字节数组-->Base64URL 编码函数
function arrayBufferToBase64URL(buffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    // 标准 Base64
    const base64 = btoa(binary);
    // 转换为 Base64URL（替换字符和去掉填充）
    return base64
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
}

 function int64ToBytes(int64) {
    const buffer = new ArrayBuffer(8); // int64 是 8 个字节
    const view = new DataView(buffer);
    view.setBigUint64(0, BigInt(int64)); // 按大端序存储
    return new Uint8Array(buffer);
}

function toPublicKeyOptions(serverData) {
    return {
        publicKey:{
            challenge: Base64URL_To_arrayBuffer(serverData.Challenge),
/*            allowCredentials: serverData.AllowCreds.map(c => ({
                type: c.Type,
                id: Base64URL_To_arrayBuffer(c.CredentialID),
                transports: ["internal"],
            })),*/
            userVerification: "required",
            timeout: 1200000
        }
    };
}

function publicKeyCredentialToJSON(cred) {
    if (cred instanceof ArrayBuffer) {
        // 处理所有二进制字段
        return arrayBufferToBase64URL(cred);
    } else if (Array.isArray(cred)) {
        // 处理数组
        return cred.map(x => publicKeyCredentialToJSON(x));
    } else if (cred && typeof cred === 'object') {
        // 处理对象
        const obj = {};
        for (const key in cred) {
            obj[key] = publicKeyCredentialToJSON(cred[key]);
        }
        return obj;
    } else {
        // string / number / boolean / null / undefined 原样返回
        return cred;
    }
}

export    {int64ToBytes,arrayBufferToBase64URL,Base64URL_To_arrayBuffer,toPublicKeyOptions,publicKeyCredentialToJSON}
