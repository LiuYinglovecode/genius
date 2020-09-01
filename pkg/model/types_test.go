package model

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Time ADCTime   `json:"time"`
	ID   ADCInt64  `json:"id"`
	Name ADCString `json:"name"`
}

func TestTypesDesSer(t *testing.T) {
	tt := NewADCTime("2019-02-27 18:49:15")
	assert.NotNil(t, tt)

	id := NewADCInt64(100)
	name := NewADCString("sam")
	obj := testObj{Time: tt, ID: id, Name: name}

	res, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
	}

	var objBack testObj
	err = json.Unmarshal(res, &objBack)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, objBack.Time, obj.Time)
	assert.Equal(t, objBack.ID, obj.ID)
	assert.Equal(t, objBack.Name, obj.Name)
}

func TestSync(t *testing.T) {
	// 测试channel的使用方法
	// 1. 初始化一个无缓冲的chan
	// 2. 在goroutine initFun中关闭
	// 3. 在主线程读chan 因为是无缓冲的所以主线程会一直等待，除非该chan被关闭。
	// 所以在关闭chan后,wait自然也就结束了。主线程可以正常退出。
	sync := make(chan struct{})
	initFun := func() {
		fmt.Println("start init func")
		time.Sleep(5 * time.Second)
		fmt.Println("close sync")
		close(sync)
	}

	go initFun()

	fmt.Println("wait init fun")

	<-sync

	fmt.Println("complete waiting")

	fmt.Println("====== start another test")
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("监控退出，停止了...")
				return
			default:
				fmt.Println("goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	stop <- true
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}
