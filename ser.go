/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-05-27 14:25:33
 * @LastEditTime: 2023-05-28 11:57:24
 */
package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
	sl "go.bug.st/serial"
)

func main2() {
	ports, err := sl.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	// 配置串口参数
	config := &serial.Config{
		Name:        "COM3",      // 串口设备名称
		Baud:        57600,       // 波特率
		ReadTimeout: time.Second, // 读超时时间
	}

	// 打开串口
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal("失败", err)
	}
	hexStr := "02FFFF210000DF" // 十六进制字符串
	bytes, err := hex.DecodeString(hexStr)
	// 写入数据
	data := bytes
	_, err = port.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	// 读取数据
	buf := make([]byte, 128)
	n, err := port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// 处理读取的数据
	//receivedData := buf[:n]
	fmt.Println(n)
	log.Println("Received:", buf)

	// 关闭串口
	err = port.Close()
	if err != nil {
		log.Fatal(err)
	}
}
