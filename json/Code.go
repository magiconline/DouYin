package json

const (
	SUCCESS        = 0   //成功响应
	ERROR          = 500 //错误响应
	INVALID_PARAMA = 400 //请求参数无效

	ERROR_AUTH          = 10001 //无效的用户
	ERROR_EXIST_AUTH    = 10002 //用户名已存在
	ERROR_TOKEN         = 10003 //token错误
	ERROR_TOKEN_FAIL    = 10004 //token无效
	ERROR_TOKEN_TIMEOUT = 10005 //token超时
	ERROR_PASSWORD      = 10006 //密码错误
)
