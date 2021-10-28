package main

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gbindata"
	"os"
)

func main() {
	fmt.Println("Example:\nenc souce-binary-filename package-name var-name")
	if len(os.Args) != 4 {
		fmt.Println("Arguments number should be 4")
		return
	}

	binfile := os.Args[1]
	pkgname := os.Args[2]
	varname := os.Args[3]
	fmt.Println("Start to encode", binfile)
	err := gbindata.Enc(binfile, binfile+".go", pkgname, varname)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(binfile, "encode success")
}
