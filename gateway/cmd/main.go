package main

import (
	"fmt"

	_ "github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/Zhan9Yunhua/blog-svr/gateway/logger"
	lg "github.com/Zhan9Yunhua/logger"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		lg.Fatalln(err)
	}
	fmt.Println(logger)
}
