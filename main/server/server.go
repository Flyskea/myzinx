package main

import (
	"fmt"
	"myzinx/ziface"
	"myzinx/znet"
	"myzinx/ztimer"
	"time"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (hzr *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello Zinx Router V0.10"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	// =============设置两个链接属性，在连接创建之后===========
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Flyskea")
	conn.SetProperty("Home", "https:// github.com/Flyskea")
	// ===================================================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	// ============在连接销毁之前，查询conn的Name，Home属性=====
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	// ===================================================

	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	// 创建一个server句柄
	s := znet.NewServer(&znet.ServerOption{
		IPVersion:  "tcp4",
		IP:         "127.0.0.1",
		Packer:     znet.NewDataPack(),
		Port:       7777,
		ConnMgr:    znet.NewConnManager(),
		MsgHandler: znet.NewMsgHandle(),
		Name:       "fuck",
	})

	// 注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	timer := ztimer.NewAutoExecTimerScheduler()
	timer.CreateTimerAfter(ztimer.NewDelayFunc(func(v ...interface{}) {
		fmt.Println("我是延迟函数")
	}, nil), 1*time.Microsecond)
	timer.CreateTimerAfter(ztimer.NewDelayFunc(func(v ...interface{}) {
		fmt.Println("我是延迟函数")
	}, nil), 10*time.Second)
	// 开启服务
	s.Serve()
}
