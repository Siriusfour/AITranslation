package swagger

type JwtTokenExpired struct {
	Code    int    `json:"code" example:"400100"`
	Message string `json:"message" example:"Token Expired"`
}

type JwtTokenInvalid struct {
	Code    int    `json:"code" example:"400101"`
	Message string `json:"message" example:"Token is invalid"`
}

type PasswordIsFail struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Password is fail"`
}

type GetChallengeIsFail struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"获取随机数失败：...."`
}
