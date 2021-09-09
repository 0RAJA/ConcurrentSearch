package Tire

import (
	"sync"
)

// Trie 前缀树
type Trie struct {
	prefix byte           //此结构的首字符
	result []string       //结果数组
	suffix map[byte]*Trie //后驱
}

const (
	MaxWorker = 32 //最大工人数
)

var (
	mutex      sync.Mutex            //互斥锁
	resultChan = make(chan string)   //结果通道
	downChan   = make(chan struct{}) //工人完成工作
	worksChan  = make(chan *Trie)    //任务
	nowWorker  = 0                   //目前工人
)

// Constructor /** Initialize your data structure here. */
func Constructor() Trie {
	return Trie{suffix: map[byte]*Trie{}}
}

// Insert /** Inserts a word into the trie. */
func (this *Trie) Insert(word, value string) {
	mutex.Lock()
	defer mutex.Unlock()
	root := this
	for i := 0; i < len(word); i++ {
		if root.suffix[word[i]] == nil {
			root.suffix[word[i]] = &Trie{
				prefix: word[i],
				suffix: map[byte]*Trie{},
			}
		}
		root = root.suffix[word[i]]
	}
	root.result = append(root.result, value)
}

//搜索控制中心
func searchCenter() (ret []string) {
	for {
		select {
		case result := <-resultChan:
			ret = append(ret, result)
		case <-downChan:
			nowWorker--
			if nowWorker == 0 {
				return
			}
		case works := <-worksChan:
			nowWorker++
			go search(works, true)
		}
	}
}

func search(tire *Trie, master bool) {
	if master {
		defer func() {
			downChan <- struct{}{}
		}()
	}
	if tire == nil {
		return
	}
	for _, v := range tire.suffix {
		if nowWorker < MaxWorker {
			worksChan <- v
		} else {
			search(v, false)
		}
	}
	for _, v := range tire.result {
		resultChan <- v
	}
}

// Search /** Returns if the word is in the trie. */
func (this *Trie) Search(word string) (ret []string) {
	nowWorker = 1
	go func() {
		root := this
		for i := 0; i < len(word); i++ {
			if root.suffix[word[i]] != nil {
				root = root.suffix[word[i]]
			} else {
				downChan <- struct{}{} //结束任务
				return
			}
		}
		search(root, false)
		downChan <- struct{}{} //结束任务
	}()
	return searchCenter()
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */
