package main

import (
	"fmt"
	"os"
)

func createmain() {
	var maingo string = `package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/idun886/GoBot/Context"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

// 定义接口函数  让所有注册的插件都实现这个方法 框架去自动调用这个方法 形成插件化开发
type Plugin interface {
	PluginFunc()
}

// 创建一个类  里面有属性 端口号
type WebSocketServer struct {
	SOCKET_PORT string
	Conn        *websocket.Conn
}

// 升级websocket所用的up
var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	//升级成了websocket连接   返回两个参数 一个连接 一个错误
	conn, err := UP.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		//如果有错误直接退出函数  如升级失败
		return
	}
	for {
		_, context, err := conn.ReadMessage()
		if err != nil {
			//如果出现错误 则跳出循环
			fmt.Printf("go-cqhttp客户端已退出\n")
			break
		} else {
			meta := Context.Meta{Conn: conn, Context: string(context)}
			meta.Login()
			if meta.MetaType == "lifecycle" {
				fmt.Printf("go-cqhttp已链接,登录账号为%d\n", meta.SelfID)
			}
			//fmt.Println(string(Context))
			RegisteredPlugins(conn, string(context))
		}
	}
}

func (w *WebSocketServer) ReadConfig() {

}

// http升级websocket用的方法
func Receiving() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	port := gjson.Get(string(bytes), "port").String()
	fmt.Println("程序运行在" + port + "端口")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)

}

func main() {

	Receiving()

}`
	// 创建文件
	//
	file, err := os.Create("main.go")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// WriteString写数据，返回值参数1是写的字符长度，参数2是异常信息
	len, err := file.WriteString(maingo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("main文件创建完成,写入字符%d个。\n", len)
}

func createreg() {
	var regfile string = `package main

import (
	"github.com/gorilla/websocket"
)

func RegisteredPlugins(conn *websocket.Conn, context string) {
	
}`
	file, err := os.Create("Registered.go")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// WriteString写数据，返回值参数1是写的字符长度，参数2是异常信息
	len, err := file.WriteString(regfile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("插件注册文件创建完成,写入字符%d个。\n", len)
}

func createconfig() {
	var config string = `{
  "port":8086
}
`
	file, err := os.Create("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	// WriteString写数据，返回值参数1是写的字符长度，参数2是异常信息
	len, err := file.WriteString(config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("配置文件创建完毕,写入字符%d个。\n", len)

}

func main() {
	createmain()
	createreg()
	createconfig()
}
