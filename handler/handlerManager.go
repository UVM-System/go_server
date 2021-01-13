// 管理所有的 Handler 函数
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/config"
	"go_server/photo"
	"go_server/token"
	"gocv.io/x/gocv"
	"log"
	"net/http"
)

// 处理来自树莓派的照片数据
func ImageHandler(c *gin.Context) {
	fmt.Println("============= A POST request comes =============")
	m, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Sorry this POST request is not MultipartForm")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Sorry this POST request is not MultipartForm",
		})
		return
	}
	// 检查 token
	if !token.IsLegalToken(m.Value["machineid"][0], m.Value["token"][0]) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Sorry the token of this request is illegal",
		})
		return
	}
	// 封装图片
	images := make([]photo.ImageBytes, 0)
	data := make([]byte, 8 * 1024 * 128) //文件的信息可以读取进一个[]byte切片
	fileHeaders := m.File["image"]
	for i := 0; i < len(fileHeaders); i++ {  // 每次只传输一张图片，即 len(fileHeaders) = 1
		fileHeader := fileHeaders[i]
		file, _ := fileHeader.Open()
		n, err :=file.Read(data)
		if err != nil {
			log.Println("got erro", err)
			log.Println("read byte number:", n)
		}
		// 创建图像结构
		images = append(images, photo.CreateImageBytes(data, fileHeader.Filename))
		// 写入图像中
		pic, _ := gocv.IMDecode(data, gocv.IMReadColor)
		gocv.IMWrite("./pictures/original/" + fileHeader.Filename, pic)
	}
	//// 传到远程服务器进行检测
	//respBody := photo.PostImage(images)
	//// 根据左右摄像头，处理检测结果
	//goodsList, receiveJson := receiveJsonHandle(respBody, m.Value["sequence"][0])
	//// TODO: 更新购物车
	//if Result.LeftResult1 != nil && Result.RightResult1 != nil && Result.LeftResult2 != nil && Result.RightResult2 != nil {
	//	goodsAll := addResults([]map[string] int{Result.LeftResult1, Result.RightResult1, Result.LeftResult2, Result.RightResult2})
	//	change := counter.UpdateCounter(m.Value["machineid"][0], m.Value["state"][0], goodsAll)
	//	Result.LeftResult1, Result.RightResult1, Result.LeftResult2, Result.RightResult2 = nil, nil, nil, nil        // 更新存储的结果
	//	fmt.Println(m.Value["state"][0])
	//	c.JSON(http.StatusOK, gin.H{
	//		"status":  "success",
	//		"get photo number": 2,
	//		"goodList": goodsAll,
	//		"changed": change,
	//	})
	//} else {
	//	fmt.Println(m.Value["state"][0])
	//	c.JSON(http.StatusOK, gin.H{
	//		"status":  "success",
	//		"get photo number": 1,
	//		"goodList": goodsList,
	//	})
	//}
	//drawRectangle(fileHeaders[0].Filename, receiveJson, m.Value["sequence"][0])
	c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"get photo number": 2,
		})
}

// 处理来自树莓派端的 token 请求信息
func TokenHandler(c *gin.Context) {
	fmt.Println("============= A GET request comes =============")
	machineId := c.Query("machineid")
	password := c.Query("password")
	// 检查 machineid 是否为空
	if len(machineId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  c.Request.Method,
			"message": "Failed, request token failed, the machineid is not found",
		})
	}
	// 检查 password 是否正确
	for _, account := range config.Config.Accounts {
		if account.MachineId == machineId && account.Password == password {
			machineToken := token.GetToken(machineId)
			if machineToken == "" {
				machineToken = token.CreateToken(machineId)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "success",
				"token": machineToken,
			})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  c.Request.Method,
		"message": "Failed, password is wrong or the server doesn't have this machine information",
	})
}
