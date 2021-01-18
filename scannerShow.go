package aim_vguang_tools

import (
	"errors"
	"time"
)

//离线状态显示
func (bro *BrokerImpl) showOffLine() error {
	select {
	case <-bro.done:
		return errors.New("connection is closed")
	default:
	}
	//显示 !
	bro.conn.Write(bro.cmd.cmdImg.warn)
	//红灯
	bro.conn.Write(bro.cmd.cmdLight.impassable)
	return nil
}

//允许通过显示内容：绿灯、对号 音效0
func (bro *BrokerImpl) showPassable() error {
	return bro.showCtrlTemplate(true)
}

//禁止通过显示内容：红灯、错号 音效1
func (bro *BrokerImpl) showImpassable() error {
	return bro.showCtrlTemplate(false)
}

//允许、禁止通过操作模板
func (bro *BrokerImpl) showCtrlTemplate(res bool) error {
	select {
	case <-bro.done:
		return errors.New("connection is closed")
	default:
	}

	//关灯
	bro.conn.Write(bro.cmd.cmdLight.close)

	var cmdLight []byte
	var cmdImg []byte
	var cmdSound []byte

	if res {
		cmdLight = bro.cmd.cmdLight.passable
		cmdImg = bro.cmd.cmdImg.passable
		cmdSound = bro.cmd.cmdSound.passable
		// 继电器动作1s
		cmd:=[]byte{0x55, 0xaa ,0x2a ,0x02 ,0x00 ,0x01 ,0x02 ,0xd4}
		bro.conn.Write(cmd)
	} else {
		cmdLight = bro.cmd.cmdLight.impassable
		cmdImg = bro.cmd.cmdImg.impassable
		cmdSound = bro.cmd.cmdSound.impassable
	}

	bro.conn.Write(cmdLight)
	bro.conn.Write(cmdImg)
	bro.conn.Write(cmdSound)

	//图片显示1.5s，重新回到主窗口
	ticker := time.NewTicker(1500 * time.Millisecond)
	<-ticker.C

	//显示第一张图片
	bro.conn.Write(bro.cmd.cmdImg.home)

	//关灯
	bro.conn.Write(bro.cmd.cmdLight.close)
	return nil
}
