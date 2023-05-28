package main

import (
	"jian/serialbit"
)

func main() {
	coms := serialbit.FindCom()
	oport := serialbit.OpenCom(coms[0])
	if oport != nil {
		serialbit.ReadReaderInfo(oport)
		serialbit.ReadFunc(oport)
	}

}
