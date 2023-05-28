/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-05-28 11:56:43
 * @LastEditTime: 2023-05-28 15:17:09
 */
package serialbit

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	sl "go.bug.st/serial"
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

	return port
}

func ReadReaderInfo(port sl.Port) {
	hexStr := "AA000300010004DD"
	bytes, _ := hex.DecodeString(hexStr)
	port.Write(bytes)

	// // 读取数据
	// buf := make([]byte, 128)
	// n, err := port.Read(buf)
	// if err != nil {
	// 	fmt.Println("读取失败", err)
	// }
	// fmt.Println("返回长度", n)
	// byteHex := buf[:n]
	// log.Println("Received:", byteHex)

	// bytehexStr := hex.EncodeToString(byteHex)
	// fmt.Println("HEX:", bytehexStr)
	// res := ParseByte(byteHex)
	// fmt.Println("解析结果：", res)
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
			fc.Type = buf[i+1]
			fc.Command = buf[i+2]
		}
		if buf[i] == 0xdd {

		}
		i++
	}
	return *fc
}
