package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

type ServiceSetup struct {
	ChaincodeId   string
	ChannelClient *channel.Client
	LedgerClient  *ledger.Client
}

//注册链码事件
func registerEvent(client *channel.Client, chaincodeId string, eventId string) (fab.Registration, <-chan *fab.CCEvent) {
	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeId, eventId)
	if err != nil {
		fmt.Printf("注册链码事件发生错误：%s", err)
	}
	return reg, notifier
}

// 执行链码完成后成功了吗?
func eventResult(notifier <-chan *fab.CCEvent, eventId string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件：%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能接受到链码事件：%s\n", eventId)
	}
	return nil
}
