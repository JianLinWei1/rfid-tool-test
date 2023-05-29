/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-05-28 11:56:27
 * @LastEditTime: 2023-05-29 16:35:34
 */
package main

import (
	"encoding/hex"
	"fmt"
	"jian/serialbit"
	"strconv"
	"strings"
)

// 硬件信息
const HARDWARE = "AA000300010004DD"

// 单次读取
const SINGLE_EPC = "AA0022000022DD"

// 多次读取
const MULTi_EPC = ""

// 读取TID 区
const EPC_TID = "AA003900090000000002000000064ADD"

// AA 00 49 00 19 00 00 00 00 01 00 00 00 08 1C 6F 30 00 E2 80 69 95 00 00 50 03 76 85 B1 AB 30 DD
// AA 00 49 00 19 00 00 00 00 01 00 00 00 08 0C 4E 30 00 E2 80 69 95 00 00 50 03 76 85 B1 AC 00 DD
// 0x01 EPC  0x02 TID 0x03 user
const WIRTE_EPC_START = "00490019000000000100000008"

// Selct EPC
// AA 00 12 00 01 01 14 DD
const SELECT_EPC = "AA001200010014DD"

func main() {
	hexb, _ := hex.DecodeString("004900190000000001000000080C4E3000E2806995000050037685B1AC")
	fmt.Println(serialbit.CheckSumCalc(hexb))

	coms := serialbit.FindCom()
	oport := serialbit.OpenCom(coms[0])
	if oport != nil {
		//serialbit.ReadReaderInfo(oport)
		//serialbit.ReadFunc(oport)
		for i := 0; i < 1; i++ {
			fmt.Println("第" + strconv.Itoa(i) + "次")
			fmt.Println("*********读取EPC+TID********")
			obj1 := serialbit.WriteCmd(oport, SINGLE_EPC)
			obj2 := serialbit.WriteCmd(oport, EPC_TID)
			if strings.EqualFold(obj1.EPC, obj2.EPC) {
				obj1.TID = obj2.TID
			}
			fmt.Println("合并后：", obj1)

			if (obj1 != serialbit.EpcObj{}) {
				//设置select
				serialbit.WriteCmd(oport, SELECT_EPC)
				fmt.Println("*********写EPC********")
				//写命令
				writeCmd := WIRTE_EPC_START
				writeCmd += obj1.CRC + obj1.PC + "E2806995000050037685B112"
				byteHex, _ := hex.DecodeString(writeCmd)
				checkSum := serialbit.CheckSumCalc(byteHex)
				fmt.Println(checkSum)
				fmt.Println("写命令：", writeCmd)
				writeCmd = "AA" + writeCmd + checkSum + "DD"
				fmt.Println("写命令：", writeCmd)
				serialbit.WriteCmd(oport, writeCmd)
			}

		}

	}

}
