/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-05-28 11:56:43
 * @LastEditTime: 2023-05-29 16:32:46
 */
package serialbit

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	sl "go.bug.st/serial"
)

const (
	HARDWARE_CMD byte = 0x03
	SINGLE_CMD   byte = 0x22
	EPC_TID_CMD  byte = 0x39
	WRITE_CMD    byte = 0x49
	SELECT_CMD   byte = 0x12
)

type SerialCom struct {
}
type FrameCmd struct {
	Header    byte
	Type      byte
	Command   byte
	MSB       byte
	LSB       byte
	Parameter byte
	Checksum  byte
	End       byte
}

type EpcObj struct {
	PC  string
	EPC string
	TID string
	CRC string
}

func FindCom() []string {
	ports, err := sl.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("未找到任何串口!")
	}
	fmt.Println("找到串口：", ports)
	return ports
}

func OpenCom(com string) sl.Port {
	mode := &sl.Mode{
		BaudRate: 115200,
	}
	port, err := sl.Open(com, mode)
	if err != nil {
		fmt.Println("打开端口失败！", err)
		return nil
	}
	port.SetReadTimeout(10 * time.Second)
	return port
}

func WriteCmd(port sl.Port, cmd string) EpcObj {
	epcObj := &EpcObj{}
	bytes, _ := hex.DecodeString(cmd)
	port.Write(bytes)
	time.Sleep(100 * time.Millisecond)
	// 读取数据
	buf := make([]byte, 128)
	n, err := port.Read(buf)
	if err != nil {
		fmt.Println("读取失败", err)
	}
	fmt.Println("返回长度", n)
	byteHex := buf[:n]
	log.Println("Received:", byteHex)

	bytehexStr := hex.EncodeToString(byteHex)
	fmt.Println("HEX:", bytehexStr)
	res := ParseByte(byteHex)
	fmt.Println("解析结果：", res)
	if res.Command == 0x3 {
		l := len(byteHex) - 2
		info := byteHex[6:l]
		fmt.Println("硬件信息：", hex.EncodeToString(info))
		fmt.Println("硬件信息：", string(info))
	} else if res.Command == SINGLE_CMD {
		l := len(byteHex) - 4
		pc := strings.ToUpper(hex.EncodeToString(byteHex[6:8]))
		crc := strings.ToUpper(hex.EncodeToString(byteHex[l : l+2]))
		epc := strings.ToUpper(hex.EncodeToString(byteHex[8:l]))
		fmt.Println("PC HEX ", pc)
		fmt.Println("EPC HEX ", epc)
		fmt.Println("CRC HEX ", crc)
		epcObj.PC = pc
		epcObj.EPC = epc
		epcObj.CRC = crc
	} else if res.Command == EPC_TID_CMD {
		//crc :=
		pc := strings.ToUpper(hex.EncodeToString(byteHex[6:8]))
		epc := strings.ToUpper(hex.EncodeToString(byteHex[8:20]))
		tid := strings.ToUpper(hex.EncodeToString(byteHex[20 : len(byteHex)-2]))
		fmt.Println("PC HEX ", pc)
		fmt.Println("EPC HEX ", epc)
		fmt.Println("TID HEX", tid)
		epcObj.PC = pc
		epcObj.EPC = epc
		epcObj.TID = tid
	} else if res.Command == WRITE_CMD {
		fmt.Println("写入命令")
	} else if res.Command == SELECT_CMD {
		fmt.Println("设置Select")
	} else {
		fmt.Println("结果有误")
	}
	return *epcObj
}

func ReadFunc(port sl.Port) {
	//var wg sync.WaitGroup
	// go func(port sl.Port) {
	for true {
		time.Sleep(100 * time.Millisecond)
		// 读取数据
		buf := make([]byte, 128)
		n, err := port.Read(buf)
		if err != nil {
			fmt.Println("读取失败", err)
		}
		fmt.Println("返回长度", n)
		byteHex := buf[:n]
		log.Println("Received:", byteHex)

		bytehexStr := hex.EncodeToString(byteHex)
		fmt.Println("HEX:", bytehexStr)
		res := ParseByte(byteHex)
		fmt.Println("解析结果：", res)
	}

	// }(port)
	// wg.Wait()
}

// 解析
func ParseByte(buf []byte) FrameCmd {
	fc := &FrameCmd{}

	i := 0
	for i < len(buf) {
		if buf[i] == 0xaa {
			fc.Header = buf[i]
			if len(buf) >= 3 {
				fc.Type = buf[i+1]
				fc.Command = buf[i+2]
			}
			return *fc
		}
		if buf[i] == 0xdd {

		}
		i++
	}
	return *fc
}

func ReadEPC(port sl.Port) {

}

func CheckSumCalc(hexByte []byte) string {
	sum := 0
	for _, h := range hexByte {
		sum += int(h)
	}
	lowByte := byte(sum & 0xFF)

	hexStr := fmt.Sprintf("%02X", lowByte)

	return hexStr
}
