
import api from "../../src/Utils/Api/Api.js";
import {
    arrayBufferToBase64URL,
    int64ToBytes,
    Base64URL_To_arrayBuffer,
    toPublicKeyOptions, publicKeyCredentialToJSON
} from "../../src/Utils/encode/base64.js"
import {onMounted} from "vue";



onMounted(() => {
    if (localStorage.getItem("AccessToken")!==undefined ) {
        Login()
    }

});



// 登录处理
export function Login(){
    const formData = getFormData();
    // 调用登录API
        api.post("/NotAuth/Login",formData).then((res)=> {
            if (confirm("要为该账号注册一条安全密钥吗")){
                ApplicationWebAuthn()
            }
        }).catch((error)=>{alert(error.message || '登录失败，请重试');});
}

export function LoginByWenXin(){}

export function LoginByQQ(){}

export function LoginByEmail(){}

export function LoginByWebAuthn(){

    api.get("/NotAuth/LoginGetWebAuthnInfo").then(res=>{

        console.log(res)
        const PublicKeyOptions = toPublicKeyOptions(res)
        console.log(PublicKeyOptions)


        navigator.credentials.get(PublicKeyOptions).then((credential )=>{

            console.log(publicKeyCredentialToJSON(credential));

            api.post("/NotAuth/LoginByWebAuthn", publicKeyCredentialToJSON(credential),).then((result)=>{

                alert("webAuthn登录成功！")
                //跳转到主业务页面
                //调用获取其他业务数据的接口

            })

        }).catch((err)=>{
            console.error(err)})

    }).catch(err=>console.error(err));

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
    api.get("/Auth/RegisterGetWebAuthnInfo").then((res)=>{

        console.log(res)

        const algs = [

            {type: "public-key", alg: -257},
            {type:"public-key", alg:-35},
            {type:"public-key", alg:-36},
            {type:"public-key", alg: -7},
            {type:"public-key", alg: -8},]

        const publicKeyCredentialCreationOptions = {
            challenge:Base64URL_To_arrayBuffer(res.Challenge),
            rp: {
                name: "Suis",
                id: window.location.hostname,
            },
            user: {
                id: int64ToBytes(res.UserID),
                name: res.UserName,
                displayName:res.UserName,
            },
            pubKeyCredParams:algs,
            authenticatorSelection: {
                residentKey: "required",
                userVerification: "preferred",
                authenticatorAttachment: res.WebAuthn.Attachment,
            },
            timeout: 60000,
            attestation: res.WebAuthn.Attestation
        }

        navigator.credentials.create({
            publicKey: publicKeyCredentialCreationOptions
        }).then((res)=>{

            //把字节序列通过base64url编码
            const encodedCredential = {

                id: arrayBufferToBase64URL(res.rawId),
                rawId: arrayBufferToBase64URL(res.rawId),
                type: res.type,
                response:{
                     clientDataJSON: arrayBufferToBase64URL(res.response.clientDataJSON),
                     attestationObject: arrayBufferToBase64URL(res.response.attestationObject),
                },
            };

            api.post("/Auth/RegisterWebAuthn", encodedCredential).then((res)=>{
                console.log(res)
            })
        })

    }).catch((error)=>{console.log(error)});

}
export function LoginByGithub(){


api.get("/NotAuth/GetChallenge?OAuth_provider=Github").then((res)=>{

    window.location.href =
        "https://github.com/login/oauth/authorize" +
        "?client_id=Iv23limF64jySyryK3Kx" +
        "&redirect_uri=http://localhost:5174/callback" +
        "&scope=read:user%20user:email" +
        "&state=" + res.Challenge
        "&token=" + res.Token;

}).catch((error)=>{console.log(error)});

}