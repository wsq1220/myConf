package main

import (
	"fmt"

	"github.com/wsq1220/myConf/conf"
)

func main() {
	var c = &conf.Config{}
	err := conf.ParseConf("xxx.conf", c)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", c)

}
