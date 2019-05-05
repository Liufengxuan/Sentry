package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	gomail "gopkg.in/gomail.v2"
)

type LocalInfo struct {
	ExtranetIp string
	LocalIp    []string
	HostName   string
	SystemName string
	Arch       string
}

var localInfo LocalInfo

func SendMail(subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "msfengxuan@sina.com",
		"pass": "--------~",
		"host": "smtp.sina.com",
		"port": "25",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", "Sentry<"+mailConn["user"]+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", "813949669@qq.com")               //发送给多个用户
	m.SetHeader("Subject", subject)                     //设置邮件主题
	m.SetBody("text/html", body)                        //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}
func main() {
	var a = 0
	for {
		a++
		err := CreateMail()
		if err == nil {
			return
		}
		if a > 6 {
			return
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}
func CreateMail() error {
	var err error
	localInfo.ExtranetIp, err = GetLocalPublicIpUseDnspod()
	if err != nil {
		localInfo.ExtranetIp = "获取外网IP地址失败！"
	}
	localInfo.HostName, err = os.Hostname()
	if err != nil {
		localInfo.HostName = "获取主机名失败！"
	}

	localInfo.LocalIp, err = net.LookupHost(localInfo.HostName)
	if err != nil {
		localInfo.LocalIp[0] = "获取内网IP地址失败！"
	}
	localInfo.SystemName = runtime.GOOS
	localInfo.Arch = runtime.GOARCH
	fmt.Println(localInfo)

	body := `外网地址：<font size="2" color="red">` + localInfo.ExtranetIp + "</font><br>" +
		`内网地址：<font size="2" color="red">` + fmt.Sprintln(localInfo.LocalIp) + "</font><br>" +
		"操作系统：" + localInfo.SystemName + `-` + localInfo.Arch + "<br>" +
		"<br><br><br>" +
		`<font size="1" color="black">` + time.Now().Format("2006-01-02 15:04:05") + "</font>" + " - " +
		`<font size="1" color="black">` + "  哨兵  </font>"

	//邮件主题为"Hello"
	subject := "<" + localInfo.HostName + ">已上线 "
	err2 := SendMail(subject, body)
	if err2 != nil {
		fmt.Print(err2)
		return err2
	} else {
		return nil
	}

}

func GetLocalPublicIpUseDnspod() (string, error) {
	timeout := time.Nanosecond * 30
	conn, err := net.DialTimeout("tcp", "ns1.dnspod.net:6666", timeout*time.Second)
	defer func() {
		if x := recover(); x != nil {
			log.Println("Can't get public ip", x)
		}
		if conn != nil {
			conn.Close()
		}
	}()
	if err == nil {
		var bytes []byte
		deadline := time.Now().Add(timeout * time.Second)
		err = conn.SetDeadline(deadline)
		if err != nil {
			return "", err
		}
		bytes, err = ioutil.ReadAll(conn)
		if err == nil {
			return string(bytes), nil
		}
	}
	return "", err
}
