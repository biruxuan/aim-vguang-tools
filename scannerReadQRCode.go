package aim_vguang_tools

import (
	"bytes"
	"errors"
)

type QRCode struct {
	Msg [][]byte
	Err error
}

//ReadMsg 读取设备发送的报文
func (bro *BrokerImpl) readQRCode(c chan QRCode) {
	//新建缓冲区
	maxLen := 256
	for {
		buf := make([]byte, maxLen+1)
		res := QRCode{
			Err: nil,
			Msg: nil,
		}

		//读取tcp通道内容
		reqLen, err := bro.conn.Read(buf[:maxLen])
		if err != nil {
			//返回错误信息,goroutine退出
			res.Err = err
			c <- res
			return
		}

		//获取请求体
		reqBody := bytes.Split(buf[:reqLen], []byte{0x55, 0xaa})
		if len(reqBody) < 1 {
			res.Err = errors.New("qrcode is empty")
		}

		res.Msg = reqBody
		c <- res
	}
}

//关闭TCP连接
func (bro *BrokerImpl) close() error {
	select {
	case <-bro.done:
		return errors.New("connection is closed")
	default:
	}

	close(bro.done)
	return bro.conn.Close()
}
