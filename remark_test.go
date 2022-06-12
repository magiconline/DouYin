package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)


// 测试remark没有token
func TestReamrkWithoutToken(t *testing.T) {
	//请求值的参数
	videoid := 1
	actiontype := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?video_id=%v&action_type=%v&comment_text=tokentest", videoid,actiontype), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	if int(body["status_code"].(float64)) != 0 {
		//t.Errorf("status_code: %v != 0, status_msg: %v", body["status_code"], body["status_msg"])
		//t.FailNow()
		assert.Equal(t,body["status_msg"],"用户未登录")
	}

}


// 测试错误token
func TestReamrkWithWrongToken(t *testing.T) {
	//请求值的参数
	token := "132"
	videoid := 1
	actiontype := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=tokentest", token,videoid,actiontype), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "token has some error")
}

//测试过期token
func TestRemarkWithExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	videoid := 1
	actiontype := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=tokentest", token,videoid,actiontype), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)

	result, err := regexp.Match("token has some error", []byte(body["status_msg"].(string)))
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if result != true {
		t.Fatalf("status_msg: %v != token has some error", body["status_msg"].(string))
		t.FailNow()
	}
}

//测试插入评论
func TestInsertRemark(t *testing.T)  {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk2OTE0LCJpc3MiOiJ6amN5In0.WeN4fZgkitj_ETYIvwAP-nvIPewWMIRBT4tIbX_mTYY"
	videoid  :=1
	actiontype  :=1
	commenttext := "test1"

	response :=httptest.NewRecorder()
	request, err := http.NewRequest("POST",fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=%v", token,videoid,actiontype,commenttext),nil)
	assert.Equal(t, err, nil)
	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	//如果插入失败
	if int(body["status_code"].(float64))==1 {
		assert.Equal(t,body["status_msg"].(string),"delete error")
	}else {
		assert.Equal(t, int(body["status_code"].(float64)), 0)
	}
}

//测试删除评论
func TestDeleteRemark(t *testing.T)  {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk2OTE0LCJpc3MiOiJ6amN5In0.WeN4fZgkitj_ETYIvwAP-nvIPewWMIRBT4tIbX_mTYY"
	videoid  :=1
	actiontype  :=2
	commenttext := "testdelete"
	commentid := 111

	response :=httptest.NewRecorder()
	request, err := http.NewRequest("POST",fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=%v&comment_id=%v", token,videoid,actiontype,commenttext,commentid),nil)
	assert.Equal(t, err, nil)
	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)


	//如果删除失败
	if int(body["status_code"].(float64))==1 {
		assert.Equal(t,body["status_msg"].(string),"delete error")
	}else {
		assert.Equal(t, int(body["status_code"].(float64)), 0)
	}

}
//测试查看所有评论列表
func TestRemarkList(t *testing.T)  {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk2OTE0LCJpc3MiOiJ6amN5In0.WeN4fZgkitj_ETYIvwAP-nvIPewWMIRBT4tIbX_mTYY"
	videoid  :=1

	response :=httptest.NewRecorder()
	request, err := http.NewRequest("GET",fmt.Sprintf("/douyin/comment/list/?token=%v&video_id=%v", token,videoid),nil)
	assert.Equal(t, err, nil)
	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

}
