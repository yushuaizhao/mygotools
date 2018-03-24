package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	loginfo, err := http.Get("http://net.tsinghua.edu.cn/do_login.php?action=logout")
	if err != nil {
		fmt.Println("can't logout")
		panic(err)
	}
	logoutbyte, err1 := ioutil.ReadAll(loginfo.Body)
	defer loginfo.Body.Close()
	if err1 != nil {
		fmt.Println("can't logout ")
		panic(err1)
	}
	fmt.Println(string(logoutbyte))
	fmt.Println("按回车键退出")
	fmt.Scanln()
}
