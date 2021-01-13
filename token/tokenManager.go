/**
管理 Token
1. 每个机器发送请求 token 时，生成一个长度为 tokenlength 的 token
2. 每个 token 的保质期是 60 分钟，每隔 10s 检查一次 token 的时间
*/
package token

import (
	"crypto/rand"
	"fmt"
	"go_server/config"
	"sync"
	"time"
)

type Token struct {
	MachineId string
	MachineToken string
	CreateTime time.Time
}

var (
	TokenList []Token
	mutex sync.Mutex
)
func init()  {
	TokenList = make([]Token, 0)
	go deleteTimeOutToken()
}

// 随机生成 2*num 长度的 token
func CreateToken(machineId string) string {
	b := make([]byte, config.Config.TokenLength / 2)
	rand.Read(b)
	machineToken := fmt.Sprintf("%x", b)
	t := Token{
		MachineId:    machineId,
		MachineToken: machineToken,
		CreateTime:   time.Now(),
	}
	mutex.Lock()
	TokenList = append(TokenList, t)
	mutex.Unlock()
	return machineToken
}

// 检测 token 是否合法
func IsLegalToken(machineid string, machinetoken string) bool {
	for _, t := range TokenList {
		if t.MachineId == machineid && t.MachineToken == machinetoken {
			return true
		}
	}
	return false
}

// 根据 machineid 获取 token
func GetToken(machineId string) string {
	for _, t := range TokenList {
		if t.MachineId == machineId {
			return t.MachineToken
		}
	}
	return ""
}

// 检测是否 token 是否过期
func deleteTimeOutToken()  {
	// 每十秒检查一下 token 是否过期
	for {
		mutex.Lock()
		for i := 0; i < len(TokenList); {
			if time.Now().Sub(TokenList[i].CreateTime).Minutes() >= 60 {
				fmt.Println("remove the time out token ", TokenList[i])
				TokenList = append(TokenList[:i], TokenList[i+1:]...)
				fmt.Println(TokenList)
			} else {
				i++
			}
		}
		mutex.Unlock()
		time.Sleep(10 * time.Second)
	}
}
