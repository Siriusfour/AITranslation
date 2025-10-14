
import Api  from "./src/Utils/Api/Api"


//判断store是否有token，判断请求一个接口判断现存的token是否有效
const AccessToken = localStorage.getItem("AccessToken")
const RefreshToken = localStorage.getItem("RefreshToken")

if (AccessToken!=null && RefreshToken!=null ){

    await Api.get("/product/user_tags").then(()=>{
        location.replace("http://localhost:5173/src/pages/main/main.html");
    }).catch(()=>{  throw new Error('请重新登录');})

}

// 登录处理
document.getElementById('loginBtn').addEventListener('click', function() {
    const formData = getFormData();

    // 调用登录API
     login(formData).then(()=>{}).catch((error)=>{alert(error.message || '登录失败，请重试');});
});

const login = async (loginInfo) => {
    try {
        const loginResponse = await fetchLogin(loginInfo); // 等待异步结果

        if (loginResponse) {
            location.replace("http://localhost:5173/src/pages/main/main.html");
        }

    } catch (error) {
        return error
    }
};



function getFormData() {
    const accountInput = document.getElementById('accountInput');
    const passwordInput = document.getElementById('passwordInput');

    return {
        username: accountInput.value,
        password: passwordInput.value
    };
}

