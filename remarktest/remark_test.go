package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)


func  TestView_video_remark(t *testing.T) {
	geturl :="http://localhost:8080/douyin/comment/list/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjU0Njg5NjIxLCJpc3MiOiJ6amN5In0.PH9g1kWMMJq0sgUqqCHSEt8exEUcL6ZV_LnrTPwo8Lg&video_id=1"
	resp,err:=http.Get(geturl)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestLeave_remark(t *testing.T)  {

	posturl:="http://localhost:8080/douyin/comment/action/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjU0Njg5NjIxLCJpc3MiOiJ6amN5In0.PH9g1kWMMJq0sgUqqCHSEt8exEUcL6ZV_LnrTPwo8Lg&video_id=1&action_type=1&comment_text=aaaa"
	for i :=0; i<5;i++{
		resp,err:=http.Post(posturl,"application/x-www-form-urlencoded",strings.NewReader("id=1"))
		if err != nil{
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}

}

func TestLeave_remark1(t *testing.T)  {
	posturl:="http://localhost:8080/douyin/comment/action/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjU0Njg5NjIxLCJpc3MiOiJ6amN5In0.PH9g1kWMMJq0sgUqqCHSEt8exEUcL6ZV_LnrTPwo8Lg&video_id=1&action_type=2&comment_text=asfafs"
	resp,err:=http.Post(posturl,"application/x-www-form-urlencoded",strings.NewReader("id=1"))
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}