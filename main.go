//网上找的代码，修改了可以接收参数。然后可以用于zabbix报警。
//编译前需要修改部分代码，看代码注释，一般修改三处即可，邮件发送人邮箱，授权码，smtp服务器地址。
//golang编译过程
//在centos上按装golang即可，没有什么第三方包。
//在root家目录创建go目录，go目录下创建pkg，src，bin三个目录
//在src目录下创建main目录，main目录创建man.go文件，将代码粘贴进去即可
//编译命令  cd /root/go/src/main;go build main.go
//生成可执行文件，加上三个参数即可运行。

//modify by 87733838@qq.com

package main

import (
	"fmt"
	"net/smtp"
	"crypto/tls"
	"net"
	"log"
	"os"
)

var Usage = func() {
	fmt.Println("Usage: COMMAND args1 args2 args3")
	fmt.Println("args1 is email address")
	fmt.Println("args2 is the mesages's title")
	fmt.Println("args3 is messages's content")
}

func main(){

	args := os.Args

	if args == nil || len(args) < 2 {
		Usage() //如果用户没有输入,或参数个数不够,则调用该函数提示用户
		return
	}

	host :="smtp.163.com"  //编译前需要修改。可以是其他邮件服务提供商的例如smtp.qq.com
	port :=465             //端口一般不需要修改
	email  := "yang@163.com" //编译前需要修改，发送邮件的地址，就是使用哪个邮箱地址给被人发送邮件
	pwd := "yang"  // 编译前修改按需修改，这里填你的授权码 ，邮件服务提供商处申请授权码
	toEmail := &args[1]  // 目标地址 ，这个是程序运行是的参数。

	header   :=  make(map[string]string)

	header["From"] = "test"+"<" +email+">"
	header["To"] = *toEmail
	header["Subject"] = args[2]
	header["Content-Type"] = "text/html;chartset=UTF-8"

	body  := args[3]

	message := ""

	for k,v :=range header{
		message  += fmt.Sprintf("%s:%s\r\n",k,v)
	}

	message +="\r\n"+body


	auth :=smtp.PlainAuth(
		"",
		email,
		pwd,
		host,
	)

	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{*toEmail},
		[]byte(message),
	)

	if err  !=nil{
		panic(err)
	}

}


//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Panicln("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}