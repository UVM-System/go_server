// 管理每个机器的货柜信息
package counter

import (
	"go_server/config"
	"log"
	"strings"
)


type Goods struct {
	Name string
	Number int
}

var Counters map[string][]Goods

func init() {
	Counters = make(map[string][]Goods)
}

// UpdateCounter 接收检测的结果，如果是第一次检测，则存储信息，返回 nil；否则，更新冰箱商品数量，并返回商品数量的改变量
func UpdateCounter(machineid string, state string, goodsList []Goods) map[string]int {
	// 开门初始化
	if strings.EqualFold(state,"start") {
		Counters[machineid] = goodsList
		return nil
	} else if strings.EqualFold(state,"end") {  // 拿出结算
		goodsListPre := Counters[machineid]
		change := calChangedNumber(goodsListPre, goodsList)
		//// 更新 Counters
		//Counters[machineid] = goodsList
		return change
	}
	log.Println("receive illegal message")
	return nil
}
// 创建一个商品列表
func createGoodsList(detectResult map[string]int) []Goods {
	goodsList := make([]Goods, 0)
	for _, goodsName := range config.Config.Goods {
		num := 0
		// 检查是否存在该商品
		if _, ok := detectResult[goodsName]; ok {
			num = detectResult[goodsName]
		}
		goods := Goods{
			Name:   goodsName,
			Number: num,
		}
		goodsList = append(goodsList, goods)
	}
	return goodsList
}

// 计算商品改变数目，没变的不计数
func calChangedNumber(goodsPre []Goods, goodsNow []Goods) map[string]int {
	change := make(map[string]int)
	for i := 0; i < len(goodsNow); i++ {
		num := goodsNow[i].Number - goodsPre[i].Number
		if num != 0 {
			change[goodsNow[i].Name] = num
		}
	}
	return change
}
