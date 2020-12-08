package aim_vguang_tools

import "net"

type Broker interface {
	//对扫码器进行归一化处理
	normalized() error
	//获取设备ID
	getScannerID() (string, error)
	//查询设备状态
	getScannerStatus() error
	//获取设备信息
	getScannerInfo() (string, error)
	//设备离线显示内容: 红灯、!
	showOffLine() error
	//允许通过显示内容：绿灯、对号 音效0
	showPassable() error
	//禁止通过显示内容：红灯、错号 音效1
	showImpassable() error
	//允许、禁止通过操作模板
	//showCtrlTemplate(res string) error
	//读取TCP通道中二维码内容
	readQRCode(c chan QRCode)
	//关闭TCP连接
	close() error
}

//扫码器控制信息
type BrokerImpl struct {
	conn  net.Conn //tcp连接句柄
	cmd   *allCmd
	macID string
	done  chan struct{}
}

func NewBroker(conn net.Conn) *BrokerImpl {
	return &BrokerImpl{
		conn: conn,
		done: make(chan struct{}),
		cmd:  newAllCmd(),
	}
}
