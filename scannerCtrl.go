package aim_vguang_tools

import (
	"net"
)

//扫码器控制信息
type ScannerCtrl struct {
	bro Broker
}

func NewScannerCtrl(conn net.Conn) *ScannerCtrl {
	return &ScannerCtrl{
		bro: NewBroker(conn),
	}
}

func (sc *ScannerCtrl) ReadQRCode(c chan QRCode) {
	sc.bro.readQRCode(c)
}

//禁止通过显示内容：红灯、错号 音效1
func (sc *ScannerCtrl) ShowPassOrNo(t bool) {
	if t {
		sc.bro.showPassable() //允许通过显示内容：绿灯、对号 音效0
	} else {
		sc.bro.showImpassable() //禁止通过显示内容：红灯、错号 音效1
	}
}

//对扫码器进行归一化处理
func (sc *ScannerCtrl) Normalized() {
	sc.bro.normalized()
}

//获取设备ID
func (sc *ScannerCtrl) GetScannerID() (string, error) {
	return sc.bro.getScannerID()
}

//查询设备状态
func (sc *ScannerCtrl) GetScannerStatus() error {
	return sc.bro.getScannerStatus()
}

//查询设备信息
func (sc *ScannerCtrl) GetScannerInfo() (string, error) {
	return sc.bro.getScannerInfo()
}

//关闭设备连接
func (sc *ScannerCtrl) Close() error {
	return sc.bro.close()
}
