package main

import (
	"fmt"

	"github.com/KrysPow/go_blog_aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	err = config.SetUser("Chris", conf)
	if err != nil {
		fmt.Println(err)
	}
	conf, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", conf)
	return
}
