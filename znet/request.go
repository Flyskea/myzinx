package znet

import (
	"fmt"
	"myzinx/ziface"
)

type Request struct {
	conn  ziface.IConnection // 已经和客户端建立好的 链接
	msg   ziface.IMessage    // 客户端请求的数据
	codec ziface.ICodec      //编解码器
}

// 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求的消息的ID
func (r *Request) GetMsgID() uint64 {
	return r.msg.GetMsgID()
}

func (r *Request) Bind(v interface{}) error {
	if r.codec == nil {
		return fmt.Errorf("message codec is nil")
	}
	return r.codec.Decode(r.GetData(), v)
}
