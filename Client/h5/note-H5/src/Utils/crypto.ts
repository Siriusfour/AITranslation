import CryptoJS from "crypto-js";

export const base62ToDecimalChars =
  "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";

export const encode10to62 = (num: number) => {
  let result = "";
  do {
    let remainder = num % 62;
    result += base62ToDecimalChars.charAt(remainder);
    num = Math.floor(num / 62);
  } while (num > 0);
  return result.split("").reverse().join("");
};

export const decode62to10 = (str: string) => {
  let result = 0;
  for (let i = 0, len = str.length; i < len; i++) {
    let index = base62ToDecimalChars.indexOf(str[i]);
    let power = len - i - 1;
    result += index * Math.pow(62, power);
  }
  return result;
};

export const encrypt = (plaintext: string, key: string, iv: string) => {
  const keyHex = CryptoJS.enc.Utf8.parse(key);
  const ivHex = CryptoJS.enc.Utf8.parse(iv);
  const ciphertext = CryptoJS.AES.encrypt(plaintext, keyHex, {
    iv: ivHex,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });

  return ciphertext.toString();
};

export const decrypt = (ciphertext: string, key: string, iv: string) => {
  const keyHex = CryptoJS.enc.Utf8.parse(key);
  const ivHex = CryptoJS.enc.Utf8.parse(iv);
  const bytes = CryptoJS.AES.decrypt(ciphertext, keyHex, {
    iv: ivHex,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });

  return bytes.toString(CryptoJS.enc.Utf8);
};
