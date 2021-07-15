// 管理所有的 Handler 函数
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
	"go_server/counter"
	"go_server/photo"
	"gocv.io/x/gocv"
	"io/ioutil"
	"log"
	"net/http"
)

type ProductInfo struct {
	productId int `json:"productId"`
	ImageUrl  string `json:"imageUrl"`
	Name     string `json:"name"`
	Price float64 `json:"price"`
	Number int `json:"number"`
}

// ImageHandler 处理来自树莓派的照片数据
func ImageHandler(c *gin.Context) {
	fmt.Println("============= A POST request comes =============")
	m, err := c.MultipartForm()
	if err != nil {
		fmt.Println("Sorry this POST request is not MultipartForm")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Sorry this POST request is not MultipartForm",
		})
		return
	}
	// 封装图片
	images := make([]photo.ImageBytes, 0)
	data := make([]byte, 8*1024*128) //文件的信息可以读取进一个[]byte切片
	fileHeaders := m.File["image"]
	for i := 0; i < len(fileHeaders); i++ { // 每次只传输一张图片，即 len(fileHeaders) = 1
		fileHeader := fileHeaders[i]
		file, _ := fileHeader.Open()
		n, err := file.Read(data)
		if err != nil {
			log.Println("got error", err)
			log.Println("read byte number:", n)
		}
		// 创建图像结构
		images = append(images, photo.CreateImageBytes(data, fileHeader.Filename))
		// 写入图像中
		pic, _ := gocv.IMDecode(data, gocv.IMReadColor)
		gocv.IMWrite(originalFilepath+fileHeader.Filename, pic)
		log.Println(fileHeader.Filename)
		log.Println(m.Value["sequence"][0])
	}
	// 传到远程服务器进行检测
	respBody := photo.PostImage(images)
	// 根据左右摄像头，处理检测结果
	goodsList, receiveJson := receiveJsonHandle(respBody, m.Value["sequence"][0])
	// TODO: 更新购物车
	if Result.LeftResult1 != nil && Result.RightResult1 != nil && Result.LeftResult2 != nil && Result.RightResult2 != nil {
		goodsAll := addResults([]map[string] int{Result.LeftResult1, Result.RightResult1, Result.LeftResult2, Result.RightResult2})
		change := counter.UpdateCounter(m.Value["machineid"][0], m.Value["state"][0], goodsAll)
		Result.LeftResult1, Result.RightResult1, Result.LeftResult2, Result.RightResult2 = nil, nil, nil, nil        // 更新存储的结果
		fmt.Println(m.Value["state"][0])
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"get photo number": 4,
			"goodList": goodsAll,
			"changed": change,
		})
	} else {
		fmt.Println(m.Value["state"][0])
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"get photo number": 1,
			"goodList": goodsList,
		})
	}
	drawRectangle(fileHeaders[0].Filename, receiveJson, m.Value["sequence"][0])
}

// ResultHandler 处理来自客户端的请求
func ResultHandler(c *gin.Context)  {
	var change []ProductInfo
	for k, v := range counter.Change {
		resp, err := http.Get("http://10.249.47.213:8000/business/product/getInfoByEN?EnglishName=" + k)
		if err != nil {
			log.Println(err)
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		jsonq := gojsonq.New().FromString(string(body))
		p := int(jsonq.Find("data.product.id").(float64))
		jsonq.Reset()
		i := jsonq.Find("data.product.image_url").(string)
		jsonq.Reset()
		n := jsonq.Find("data.product.name").(string)
		jsonq.Reset()
		c := jsonq.Find("data.product.price").(float64)
		change = append(change, ProductInfo{
			productId: p,
			ImageUrl: i,
			Name: n,
			Price: c,
			Number: v,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"change": change,
	})
}
