
import api from "./src/Utils/Api/Api.js";
import decodeBase64URL, {arrayBufferToBase64URL} from "./src/Utils/encode/base64.js"
import {int64ToBytes} from "./src/Utils/encode/base64.js"

//判断store是否有token，判断请求一个接口判断现存的token是否有效
// const AccessToken = localStorage.getItem("AccessToken")
// const RefreshToken = localStorage.getItem("RefreshToken")

// if (AccessToken!=null && RefreshToken!=null ){
//
//     await Api.get("/product/user_tags").then(()=>{
//         location.replace("http://localhost:5173/src/pages/main/main.html");
//     }).catch(()=>{  throw new Error('请重新登录');})
//
// }

// 登录处理
document.getElementById('loginBtn').addEventListener('click', function() {
    const formData = getFormData();

    // 调用登录API
     login(formData).then(()=>{}).catch((error)=>{alert(error.message || '登录失败，请重试');});
});

const login = async (loginInfo) => {

    api.post("/Auth/Login",loginInfo).then((res)=>{

        localStorage.setItem("AccessToken",res.AccessToken)
        localStorage.setItem("RefreshToken",res.RefreshToken)

        const ok =confirm("你需要为该账号注册一个秘钥吗？")
        if(ok){
            ApplicationWebAuthn()
        }

    }).catch((error)=>{alert(error.message || "登录失败！")});

};

function getFormData() {

    const accountInput = document.getElementById('accountInput');
    const passwordInput = document.getElementById('passwordInput');

    return {
        Email: accountInput.value,
        Password: passwordInput.value
    };
}

function ApplicationWebAuthn(){

    // 检查设备是否支持平台认证器
    if (window.PublicKeyCredential) {
        PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable()
            .then(available => {
                console.log("平台认证器可用:", available);
                if (!available) {
                    alert("你的设备不支持平台认证器（如 Windows Hello、Touch ID）");
                }
            });
    }


    api.get("/Auth/ApplicationWebAuthn").then((res)=>{

        const algs = [

            {type: "public-key", alg: -257},
            {type:"public-key", alg:-35},
            {type:"public-key", alg:-36},
            {type:"public-key", alg: -7},
            {type:"public-key", alg: -8},]
        const publicKeyCredentialCreationOptions = {
            challenge:decodeBase64URL(res.Challenge),
            rp: {
                name: "Susi",
                id: window.location.hostname,
            },
            user: {
                id: int64ToBytes(res.Config.UserID),
                name: res.Config.UserName ,
                displayName:res.Config.UserName,
            },
            pubKeyCredParams:algs,
            authenticatorSelection: {
                authenticatorAttachment: res.Config.Attachment,
            },
            timeout: 60000,
            attestation: res.Config.Attestation
        }

          navigator.credentials.create({
            publicKey: publicKeyCredentialCreationOptions
        }).then((res)=>{
            console.log(res)

              //把字节序列通过base64编码
              const encodedCredential = {
                      id: arrayBufferToBase64URL(res.rawId),
                      rawId: arrayBufferToBase64URL(res.rawId),
                      type: res.type,

                          clientDataJSON: arrayBufferToBase64URL(
                              res.response.clientDataJSON
                          ),
                          attestationObject: arrayBufferToBase64URL(
                              res.response.attestationObject
                          ),
                          transports: res.response.getTransports?.()

                  };


            api.post("/Auth/VerifyWebAuthn", encodedCredential).then((res)=>{
                console.log(res)
            })
        })

    }).catch((error)=>{console.log(error)});

}

