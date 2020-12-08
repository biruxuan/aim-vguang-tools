package aim_vguang_tools

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

//对扫码器进行归一化处理
func (bro *BrokerImpl) normalized() error {
	select {
	case <-bro.done:
		return errors.New("connection is closed")
	default:
	}

	//关灯
	bro.conn.Write(bro.cmd.cmdLight.close)
	//音效
	bro.conn.Write(bro.cmd.cmdSound.passable)
	//图片
	bro.conn.Write(bro.cmd.cmdImg.home)
	//间隔扫码模式
	bro.conn.Write(bro.cmd.cmdScan.workMode)

	return nil
}

//获取设备ID
func (bro *BrokerImpl) getScannerID() (string, error) {
	select {
	case <-bro.done:
		return "", errors.New("connection is closed")
	default:
	}

	//发送获取设备号报文
	cmd := []byte{0x55, 0xAA, 0x02, 0x00, 0x00, 0xFD}
	bro.conn.Write(cmd)

	//获取设备号等待时间
	ticker := time.NewTicker(5 * 1000 * time.Millisecond)
	defer ticker.Stop()

	bufChan := make(chan []byte)

	go func() {
		for {
			//创建缓冲区
			maxLen := 16
			buf := make([]byte, maxLen+1)
			//读取设备响应报文
			_, err := bro.conn.Read(buf[:maxLen])
			if err != nil {
				//return "", err
				return
			}

			//判断报文类型
			if buf[2] != 0x02 {
				//msgType := fmt.Sprintf("%d", buf[2])
				//return "", errors.New("msg type is " + msgType)
				continue
			}
			//得到设备号报文
			bufChan <- buf
			return
		}
	}()

	select {
	case <-ticker.C:
		return "", errors.New("wait for the macID time out")
	case buf := <-bufChan:
		//判断设备是否正常
		if buf[3] != 0x00 {
			return "", errors.New("scanner is abnormal")
		}

		//转为10进制
		macID := fmt.Sprintf("%d", buf[6])
		macID = strings.Replace(macID, " ", "", -1)
		bro.macID = macID
		return bro.macID, nil
	}
}

//查询设备状态
func (bro *BrokerImpl) getScannerStatus() error {
	select {
	case <-bro.done:
		return errors.New("connection is closed")
	default:
	}

	cmd := []byte{0x55, 0xAA, 0x01, 0x00, 0x00, 0xFE}
	_, err := bro.conn.Write(cmd)

	return err
}

//查询设备信息
func (bro *BrokerImpl) getScannerInfo() (string, error) {
	select {
	case <-bro.done:
		return "", errors.New("connection is closed")
	default:
	}
	return bro.conn.RemoteAddr().String() + " " + bro.macID + " ", nil
}
