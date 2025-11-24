
import api from "../../src/Utils/Api/Api.js";
import {
    arrayBufferToBase64URL,
    int64ToBytes,
    Base64URL_To_arrayBuffer,
    toPublicKeyOptions
} from "../../src/Utils/encode/base64.js"



// 登录处理
export function Login(){
    const formData = getFormData();

    // 调用登录API
    login(formData).then(()=>{

        if (confirm("要为该账号注册一条安全密钥吗")){
            ApplicationWebAuthn()
        }

    }).catch((error)=>{alert(error.message || '登录失败，请重试');});
}

export function LoginByWenXin(){}

export function LoginByQQ(){}

export function LoginByEmail(){}

export function LoginByWebAuthn(){

    api.get("/Auth/GetUserAllCredential").then(res=>{

       let PublicKeyOptions = toPublicKeyOptions(res)

        navigator.credentials.get(PublicKeyOptions).then((PublicKeyCredential )=>{
            console.log(PublicKeyCredential)

            api.post("/Auth/LoginByWebAuthn",PublicKeyCredential,).then((result)=>{

                alert("webAuthn登录成功！")
                //跳转到主业务页面
                //调用获取其他业务数据的接口

            })

        }).catch((err)=>{alert(err)})

    }).catch(err=>{alert(err)});

}


function login(loginInfo){

   return  api.post("/Auth/Login",loginInfo).then((res)=> {
            localStorage.setItem("AccessToken", res.AccessToken)
            localStorage.setItem("RefreshToken", res.RefreshToken)
    })
}
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
            challenge:Base64URL_To_arrayBuffer(res.Challenge),
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
                residentKey: "required",
                userVerification: "preferred",
                authenticatorAttachment: res.Config.Attachment,
            },
            timeout: 60000,
            attestation: res.Config.Attestation
        }

        navigator.credentials.create({
            publicKey: publicKeyCredentialCreationOptions
        }).then((res)=>{

            //把字节序列通过base64编码
            const encodedCredential = {
                id: arrayBufferToBase64URL(res.rawId),
                rawId: arrayBufferToBase64URL(res.rawId),
                type: res.type,
                response:{
                     clientDataJSON: arrayBufferToBase64URL(res.response.clientDataJSON),
                     attestationObject: arrayBufferToBase64URL(res.response.attestationObject),
                },

                transports: res.response.getTransports?.()
            };
            console.log(encodedCredential);
            api.post("/Auth/RegisterWebAuthn", encodedCredential).then((res)=>{
                console.log(res)
            })
        })

    }).catch((error)=>{console.log(error)});

}
export function LoginByGithub(){

api.get("/Auth/GetChallenge?OAuth_provider=Github").then((res)=>{

    window.location.href =
        "https://github.com/login/oauth/authorize" +
        "?client_id=Ov23lis7yW3qDZV1ARrr" +
        "&redirect_uri=http://localhost:5174/callback" +
        "&scope=read:user%20user:email" +
        "&state=" + res;

}).catch((error)=>{console.log(error)});

}