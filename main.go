/**
与树莓派和目标检测服务器连接的中间服务器
1. 处理来自树莓派的照片数据
2. 处理来自远程服务器的检测数据（json格式），根据传来的信息，判断客户买了什么商品
3. 向用户设备发送购物车数据（具体功能待定）
*/
package main

import (
	"go_server/handler"
	"sync"

	"github.com/gin-gonic/gin"
)

type Counter struct {
	MachineId string
	GoodsList []Goods
}

type Goods struct {
	Name   string
	Number int
}

var Counters []Counter

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// 接收照片数据 /token待定
	router := gin.Default()
	// 调用 POST 方法，传入路由参数和路由函数
	router.POST("/photo", handler.ImageHandler)
	// 监听 Get 请求，
	router.GET("/result", handler.ResultHandler)
	// 监听端口 8000
	router.Run(":8000")
	wg.Wait()
}
