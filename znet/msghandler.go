package znet

import (
	"fmt"
	"myzinx/config"
	"myzinx/skeapool"
	"myzinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint64]ziface.IRouter // 存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint64                    // 业务工作Worker池的数量
	pool           *skeapool.Pool            //worker池
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint64]ziface.IRouter),
		WorkerPoolSize: config.GlobalObject.WorkerPoolSize,
	}
}

// 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint64, router ziface.IRouter) {
	// 1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgID = " + fmt.Sprintf("%d", msgID))
	}
	// 2 添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}

// 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	var err error
	mh.pool, err = skeapool.NewPool(int(mh.WorkerPoolSize))
	if err != nil {
		panic(err)
	}
}

// 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	mh.pool.Submit(
		func() {
			mh.DoMsgHandler(request)
		},
	)
}
