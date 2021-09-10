package server

import (
	mTire "ConcurrentSearch/utils/Tire"
	"os"
	"sync"
)

/*不要通过共享内存来通信，要通过通信来共享内存*/

const (
	MaxWorker = 32 //最大工人数
)

var (
	name          = "你的"              //文件名
	nowWorker     = 0                 //目前工作工人数
	searchRequest = make(chan string) //发现需要工人执行的任务的通道
	doneWork      = make(chan bool)   //通知工人完成工作的通道
	tire          = mTire.Constructor()
	paths         = "CDEFGHIJKLMNOPQRSTUVWXYZ" //搜索路径
	mutex         = sync.Mutex{}               //保证工人数不会超
)

func Constructor() {
	tire = mTire.Constructor()
}

//路径格式化
func pathFormat(p string) string {
	return p + ":/"
}

// InputMessage 搜索信息
func InputMessage(n string) []string {
	name = n
	return returnResults()
}

// ToSearch 建立索引
func ToSearch() {
	nowWorker = 0 //目前占用一个工人执行启动工作
	for i := 0; i < len(paths); i++ {
		path := pathFormat(string(paths[i]))
		if _, err := os.Stat(path); err == nil {
			nowWorker++
			go search(path, true)
		}
	}
	waitingCenter()
}

// returnResults 返回结果
func returnResults() []string {
	return tire.Search(name)
}

func waitingCenter() { //控制中心
	for { //一直监听各路消息
		select {
		case path := <-searchRequest: //接收到需求,分配任务
			go search(path, true)
		case <-doneWork: //工人完成工作,增加空闲工人数
			nowWorker--
			if nowWorker == 0 { //所有工人都完成了工作就over
				return
			} //所有工人都完成了任务,结束程序
		}
	}
}

//master说明此次执行是否是由工人(go)完成的,true由go完成,false递归实现
func search(path string, master bool) {
	fileInfoList, err := os.ReadDir(path)
	if err == nil {
		for _, fileInfo := range fileInfoList {
			if path[len(path)-1] != '/' {
				path += "/"
			}
			newPath := path + fileInfo.Name()
			if fileInfo.IsDir() { //文件夹--发现任务--通知控制中心
				mutex.Lock()
				if nowWorker < MaxWorker { //有可用工人就用其他工人干
					nowWorker++
					mutex.Unlock()
					searchRequest <- newPath
				} else { //没有多余工人就自己干
					mutex.Unlock()
					search(newPath, false)
				}
			} else {
				tire.Insert(fileInfo.Name(), newPath)
			}
		}
	}
	if master { //一个任务完成了,如果是工人完成的就通知控制中心
		doneWork <- true
	}
}
