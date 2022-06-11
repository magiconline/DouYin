package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

// 测试没有token
func TestFeedWithoutToken(t *testing.T) {
	timeStamp := time.Now().UnixMilli()
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", timeStamp))
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	if int(body["status_code"].(float64)) != 0 {
		t.Errorf("status_code: %v != 0, status_msg: %v", body["status_code"], body["status_msg"])
		t.FailNow()
	}

}

// 测试错误token
func TestFeedWithWrongToken(t *testing.T) {
	timeStamp := time.Now().UnixMilli()
	token := "123"

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v&token=%v", timeStamp, token))

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "token contains an invalid number of segments")
}

// 测试过期token
func TestFeedWithExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	timeStamp := time.Now().UnixMilli()

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v&token=%v", timeStamp, token))

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)

	result, err := regexp.Match("token is expired by.*", []byte(body["status_msg"].(string)))
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if result != true {
		t.Fatalf("status_msg: %v != token is expired by*", body["status_msg"].(string))
		t.FailNow()
	}
}

// 测试访问完所有视频
func TestFeedWithEndTime(t *testing.T) {
	timeStamp := 0
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", timeStamp))

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

	l := len(body["video_list"].([]interface{}))
	assert.NotEqual(t, l, 0)
	assert.NotEqual(t, body["next_time"], timeStamp)
}

// func BenchmarkFeed(b *testing.B) {
// 	// 初始化请求
// 	timeStamp := time.Now().UnixMilli()
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", fmt.Sprintf("/douyin/feed/?latest_time=%v", timeStamp), nil)

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			// r.ServeHTTP(w, req)
// 		}

// 	})
// }

// func TestRegister(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", fmt.Sprintf("/douyin/register/?username=test02&password=123456"), nil)
// 	// r.ServeHTTP(w, req)
// 	assert.Equal(t, 200, w.Code)
// }

// func BenchmarkRedister(b *testing.B) {
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", fmt.Sprintf("/douyin/register/?username=test03&password=123456"), nil)
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			// r.ServeHTTP(w, req)
// 		}

// 	})
// }

// --------------------benchmark----------------------------------------

func BenchmarkFeed(b *testing.B) {
	req := fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", time.Now().UnixMilli())

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			http.Get(req)
		}
	})
}

/*func TestFavorite(t *testing.T) {
	// 初始化请求
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMDM2NTkzLCJpc3MiOiJ6amN5In0.HWHm5JzbcIBeiXVOWXyKV6uQNB1po6CyK8bPQRMGvSc"
	videoId := 2
	actionType := 2
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/douyin/favorite/action/?token=%v&video_id=%v&action_type=%v", token, videoId, actionType), nil)
	r.ServeHTTP(w, req)
}
func BenchmarkFavorite(b *testing.B) {
	// 初始化请求
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMDM2NTkzLCJpc3MiOiJ6amN5In0.HWHm5JzbcIBeiXVOWXyKV6uQNB1po6CyK8bPQRMGvSc"
	videoId := 2
	actionType := 2
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/douyin/favorite/action/?token=%v&video_id=%v&action_type=%v", token, videoId, actionType), nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.ServeHTTP(w, req)
		}
	})
}*/
