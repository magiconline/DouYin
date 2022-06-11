package json

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMA: "请求参数错误",

	ERROR_AUTH:          "用户不存在",
	ERROR_EXIST_AUTH:    "用户已存在",
	ERROR_TOKEN:         "token错误",
	ERROR_TOKEN_FAIL:    "token鉴权失败",
	ERROR_TOKEN_TIMEOUT: "token鉴权超时",
	ERROR_PASSWORD:      "密码错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
