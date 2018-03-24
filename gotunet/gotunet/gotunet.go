package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	var logstring string
	try := 1
	for checkCode := checkOnline(); checkCode != 2; {
		if checkCode == 0 {
			return
		}
		if checkCode == 1 {
			fmt.Println("waiting..........")
			logstring = login()
			if logstring == "E2553: Password is error." {
				return
			}
			checkCode = checkOnline()
		}
		if try != 1 {
			time.Sleep(4 * time.Second)
		}
		try = try + 1
		if try == 10 {
			fmt.Println("尝试10次失败，网络可能有问题")
			return
		}
	}
	usersplit := userInfo()
	transformBtoG := 1000.0 * 1000.0 * 1000.0
	loginSeconds := timeChange(usersplit[1])
	nowSeconds := timeChange(usersplit[2])
	download := haveUsed(usersplit[3]) / transformBtoG
	upload := haveUsed(usersplit[4]) / transformBtoG
	usage := haveUsed(usersplit[6]) / transformBtoG

	fmt.Println("用户名：", usersplit[0])
	fmt.Println("本次登录前已使用: ", usage, "G")
	fmt.Println("本次登录前还剩余: ", 25.0-usage, "G")
	fmt.Println("本次使用: ", download, "G")
	fmt.Println("本次上传: ", upload, "G")
	fmt.Println("登录时间：", time.Unix(loginSeconds, 0))
	fmt.Println("现在时间：", time.Unix(nowSeconds, 0))
	fmt.Println("本机ip: ", usersplit[8])
	fmt.Println("共使用: ", usage+download, "G")
	fmt.Println("共剩余: ", 25.0-usage-download, "G")
	fmt.Println("按回车键退出")
	fmt.Scanln()
}

//查看是否已登录
func checkOnline() int {
	var code int
	onlineinfo, err := http.Get("http://net.tsinghua.edu.cn/do_login.php?action=check_online")
	if err != nil {
		fmt.Println("can't get the check online url")
		return 0
	}
	onlinebyte, err1 := ioutil.ReadAll(onlineinfo.Body)
	defer onlineinfo.Body.Close()
	if err1 != nil {
		fmt.Println("can't check online ")
		return 0
	}
	if string(onlinebyte) == "not_online" {
		fmt.Println("not_online")
		code = 1
	} else {
		fmt.Println("online")
		code = 2
	}
	return code
}

//查看使用信息
func userInfo() []string {
	userinfo, err := http.Get("http://net.tsinghua.edu.cn/rad_user_info.php")
	if err != nil {
		fmt.Println("can't get the succeed url")
		panic(err)
	}
	userbyte, err1 := ioutil.ReadAll(userinfo.Body)
	defer userinfo.Body.Close()
	if err1 != nil {
		fmt.Println("无法获取用户信息")
		panic(err)
	}
	userstring := string(userbyte)
	usersplit := strings.Split(userstring, ",")
	return usersplit
}

//查看使用了多少流量
func haveUsed(used string) float64 {
	floatused, err := strconv.ParseFloat(used, 32)
	if err != nil {
		fmt.Println("获取流量额失败")
		panic(err)
	}
	return floatused
}

//时间转换，时间从1970年1月1日0时开始，东八区
func timeChange(sectime string) int64 {
	thistime, err := strconv.ParseInt(sectime, 10, 64)
	if err != nil {
		fmt.Println("时间有误")
		panic(err)
	}
	return thistime
}

func login() string {
	var err error
	passbyte, err := ioutil.ReadFile("./passwd.txt")
	if err != nil {
		fmt.Println("File is wrong")
		panic(err)
	}
	passwd := strings.Split(string(passbyte), ",")
	h := md5.New()
	io.WriteString(h, passwd[1])
	MD5 := hex.EncodeToString(h.Sum(nil))
	username := passwd[0]
	password := "{MD5_HEX}" + MD5
	data := map[string][]string{
		"action":   {"login"},
		"username": {username},
		"password": {password},
		"ac_id":    {"1"},
	}

	loginfo, err := http.PostForm("http://net.tsinghua.edu.cn/do_login.php", data)
	if err != nil {
		fmt.Println("can't login")
		panic(err)
	}
	loginbyte, err := ioutil.ReadAll(loginfo.Body)
	defer loginfo.Body.Close()
	if err != nil {
		fmt.Println("can't login ")
		panic(err)
	}
	loginInfo := string(loginbyte)
	fmt.Println(loginInfo)
	return loginInfo
}
