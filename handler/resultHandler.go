package handler

import (
	"encoding/json"
	"go_server/config"
	"go_server/counter"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
)

type DetectResultJson struct {
	Label string `json:"label"`
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
	Confidence float64 `json:"confidence"`
}

type ReceiveJson struct {
	FileName string `json:"filename"`
	Width int `json:"width"`
	Height int `json:"height"`
	DetectResult []DetectResultJson `json:"detectresult"`
}

type ImagesResult struct {
	LeftResult1 map[string]int
	RightResult1 map[string]int
	LeftResult2 map[string]int
	RightResult2 map[string]int
}

const(
	LeftBound int = 1320
	RightBound int = 415
)

var (
	Result ImagesResult  // 存储左右两张图片的检测结果
	originalFilepath = "./pictures/original/"
	detectedFilepath = "./pictures/detected/"
)

func init() {
	Result = ImagesResult{
		LeftResult1:  nil,
		RightResult1: nil,
		LeftResult2:  nil,
		RightResult2: nil,
	}
	err := os.MkdirAll(originalFilepath, os.ModePerm)
	if err != nil {
		log.Println(originalFilepath, " mkdir failed!!!")
	}
	err = os.MkdirAll(detectedFilepath, os.ModePerm)
	if err != nil {
		log.Println(detectedFilepath, " mkdir failed!!!")
	}
}

// 处理接收远程服务器传来的检测信息
func receiveJsonHandle(data []byte, sequence string) (map[string]int, ReceiveJson) {
	var receiveJson ReceiveJson
	err := json.Unmarshal(data, &receiveJson)
	if err != nil {
		log.Println("json unmarshal failed")
		return nil, receiveJson
	}
	log.Println("json unmarshal success")
	log.Println(receiveJson.DetectResult)
	goodsList := make(map[string]int)
	for _, detectResultJson := range receiveJson.DetectResult{
		goodsNmae := detectResultJson.Label
		//// 检查置信度，如果小于 0.5， 不计数
		//if detectResultJson.Confidence < 0.5 {
		//	continue
		//}
		// 检查坐标
		position := (detectResultJson.X1 + detectResultJson.X2) / 2
		if strings.Contains(sequence, "left") && position < LeftBound ||
			strings.Contains(sequence, "right") && position > RightBound {
			if _, ok := goodsList[goodsNmae]; ok {  // 检测是否已经加入货架中
				goodsList[goodsNmae]++
			} else {  // 如果不存在，初始化为 1
				goodsList[goodsNmae] = 1
			}
		}
	}
	if sequence == "[1]left" {
		Result.LeftResult1 = goodsList
	} else if sequence == "[1]right" {
		Result.RightResult1 = goodsList
	} else if sequence == "[2]left" {
		Result.LeftResult2 = goodsList
	} else if sequence == "[2]right" {
		Result.RightResult2 = goodsList
	}
	return goodsList, receiveJson
}

// 将多个摄像头的结果加起来
func addResults(results []map[string]int) []counter.Goods {
	goodsList := make([]counter.Goods, 0)
	for _, goodsName := range config.Config.Goods {
		num := 0
		for _, result := range results{
			if _, ok := result[goodsName]; ok {
				num += result[goodsName]
			}
		}
		goods := counter.Goods{
			Name:   goodsName,
			Number: num,
		}
		goodsList = append(goodsList, goods)
	}
	return goodsList
}

// 测试在图像上画检测出来的框
func drawRectangle(fileName string, receiveJson ReceiveJson, sequence string) {
	img := gocv.IMRead(originalFilepath + fileName, gocv.IMReadColor)
	if img.Empty() {
		log.Println("read image file failed, please check the file path")
		log.Println("filepath is: ", originalFilepath + fileName)
	}
	for _, result := range receiveJson.DetectResult{
		// 检查坐标
		position := (result.X1 + result.X2) / 2
		if strings.Contains(sequence, "left") && position < LeftBound ||
			strings.Contains(sequence, "right") && position > RightBound {
			r := image.Rect(result.X1, result.Y1, result.X2, result.Y2)
			c := color.RGBA{0, 255, 33, 255}
			gocv.Rectangle(&img, r, c, 2)
		}
	}
	gocv.IMWrite(detectedFilepath + fileName, img)
}
