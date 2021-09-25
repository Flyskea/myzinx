package znet

type Message struct {
	ID      uint64 //消息的ID
	DataLen uint64 //消息的长度
	Data    []byte //消息的内容
}

//创建一个Message消息包
func NewMsgPackage(ID uint64, data []byte) *Message {
	return &Message{
		ID:      ID,
		DataLen: uint64(len(data)),
		Data:    data,
	}
}

//获取消息数据段长度
func (msg *Message) GetDataLen() uint64 {
	return msg.DataLen
}

//获取消息ID
func (msg *Message) GetMsgID() uint64 {
	return msg.ID
}

//获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//设置消息数据段长度
func (msg *Message) SetDataLen(len uint64) {
	msg.DataLen = len
}

//设计消息ID
func (msg *Message) SetMsgID(msgID uint64) {
	msg.ID = msgID
}

//设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
