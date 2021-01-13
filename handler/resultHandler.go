package handler

import (
	"encoding/json"
	"go_server/config"
	"go_server/counter"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
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
	LeftBound int = 1200
	RightBound int = 660
)

var Result ImagesResult  // 存储左右两张图片的检测结果

func init() {
	Result = ImagesResult{
		LeftResult1:  nil,
		RightResult1: nil,
		LeftResult2:  nil,
		RightResult2: nil,
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
		if sequence == "[1]left" && (detectResultJson.X1 + detectResultJson.X2) / 2 > LeftBound {
			continue
		}
		if sequence == "[1]right" && (detectResultJson.X1 + detectResultJson.X2) / 2 < RightBound {
			continue
		}
		if sequence == "[2]left" && (detectResultJson.X1 + detectResultJson.X2) / 2 > LeftBound {
			continue
		}
		if sequence == "[2]right" && (detectResultJson.X1 + detectResultJson.X2) / 2 < RightBound {
			continue
		}
		if _, ok := goodsList[goodsNmae]; ok {  // 检测是否已经加入货架中
			goodsList[goodsNmae]++
		} else {  // 如果不存在，初始化为 1
			goodsList[goodsNmae] = 1
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
	img := gocv.IMRead(fileName, gocv.IMReadColor)
	for _, result := range receiveJson.DetectResult{
		// 检查坐标
		if sequence == "[1]left" && (result.X1 + result.X2) / 2 > LeftBound {
			continue
		} else if sequence == "[1]right" && (result.X1 + result.X2) / 2 < RightBound {
			continue
		}
		// TODO: 第二层摄像头
		if sequence == "[2]left" && (result.X1 + result.X2) / 2 > LeftBound {
			continue
		} else if sequence == "[2]right" && (result.X1 + result.X2) / 2 < RightBound {
			continue
		}
		r := image.Rect(result.X1, result.Y1, result.X2, result.Y2)
		c := color.RGBA{0, 255, 33, 255}
		gocv.Rectangle(&img, r, c, 2)
	}
	gocv.IMWrite(fileName, img)
}
