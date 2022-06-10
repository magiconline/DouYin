package benchmark

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func BenchmarkFavoriteList(b *testing.B) {
	b.ResetTimer()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	for i := 0; i < b.N; i++ {
		//模拟一个get提交请求
		resp, err := http.Get("http://127.0.0.1:8080/douyin/favorite/list/?token=" + token + "&user_id=1")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		//读取报文中所有内容
		body, err := ioutil.ReadAll(resp.Body)
		//输出内容
		fmt.Println(string(body))
	}
	b.StopTimer()
}

func BenchmarkStar(b *testing.B) {
	b.ResetTimer()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	for i := 0; i < b.N; i++ {
		//模拟一个Post提交请求
		resp, err := http.Post("http://127.0.0.1:8080/douyin/favorite/action/?token="+token+"&video_id=1&action_type=1", "application/x-www-form-urlencoded", nil)
		if err != nil {
			fmt.Println(err)
		}
		//关闭连接
		defer resp.Body.Close()
		//读取报文中所有内容
		body, err := ioutil.ReadAll(resp.Body)
		//输出内容
		fmt.Println(string(body))
	}
	b.StopTimer()
}
