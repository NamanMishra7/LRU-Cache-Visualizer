package utils

import (
	"encoding/json"
	"sync"
	"time"
)

type Message struct {
	Content   string `json:"content,omitempty"`
}

type ListNode struct {
	key string
	value string
	exp int
	time time.Time
	prev *ListNode
	next *ListNode
}

type LRUCache struct {
	
	cache map[string]*ListNode
	head *ListNode
	tail *ListNode
	mutex sync.Mutex
}

func CreateLRUCache() *LRUCache {
	l := LRUCache{
			cache:    make(map[string]*ListNode),
			head:     &ListNode{},
			tail:     &ListNode{},
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return &l
}

func (lruCache *LRUCache) addToHead(node *ListNode) {
	node.next = lruCache.head.next
	node.prev = lruCache.head
	lruCache.head.next.prev = node
	lruCache.head.next = node
}

func (lruCache *LRUCache) removeNode(node *ListNode) {
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
}

func (lruCache *LRUCache) moveToHead(node *ListNode) {
	lruCache.removeNode(node)
	lruCache.addToHead(node)
}

func (lruCache *LRUCache) Get(key string) (string, bool) {
	lruCache.mutex.Lock()
	defer lruCache.mutex.Unlock()
	if node, ok := lruCache.cache[key]; ok {
			lruCache.moveToHead(node)
			return node.value, true
	}
	return "", false
}

func (lruCache *LRUCache) Put(key string, value string, expSec int, fn func(msg string)) {
	lruCache.mutex.Lock()
	defer lruCache.mutex.Unlock()
	if node, ok := lruCache.cache[key]; ok {

			node.value = value
			lruCache.moveToHead(node)

	} else {

			newNode := &ListNode{key: key, value: value, exp: expSec, time: time.Now()}
			lruCache.cache[key] = newNode
			lruCache.addToHead(newNode)

			// if len(lruCache.cache) > lruCache.capacity {
			// 		tail := lruCache.tail.prev
			// 		lruCache.removeNode(tail)
			// 		delete(lruCache.cache, tail.key)
			// }

			// go func () {
			// 	lruCache.mutex.Lock()
			// 	defer lruCache.mutex.Unlock()
			// 	if node, ok := lruCache.cache[key]; ok {
			// 		lruCache.removeNode(node)
			// 		delete(lruCache.cache, node.key)
			// 	}
			// }()

			// time.AfterFunc(time.Duration(expSec)*time.Second, func() {
			// 	lruCache.mutex.Lock()
			// 	defer lruCache.mutex.Unlock()
			// 	if node, ok := lruCache.cache[key]; ok && node == newNode {
			// 			lruCache.removeNode(newNode)
			// 			delete(lruCache.cache, newNode.key)
			// 			str, err := json.Marshal(lruCache.ToArray())
			// 			if err != nil {
			// 				return
			// 			}
			// 			fn(string(str))
			// 	}
			// })

	}
}

func (lruCache *LRUCache) CheckForExp(fn func(str string)) {
	for key, elem := range lruCache.cache {
		diff := time.Since(elem.time)
		// time.Since(elem.time)
		if (diff >= time.Duration(elem.exp) * time.Second) {
			lruCache.removeNode(elem)
			delete(lruCache.cache, key)
			// fmt.Printf("time: %v and elapsed: %v and key: %v", elem.time, diff, key)
			str, err := json.Marshal(lruCache.ToArray())
			if err != nil {
				continue
			}
			fn(string(str))
		}
	}
}

func (lruCache *LRUCache) Delete(key string) bool {
	lruCache.mutex.Lock()
	defer lruCache.mutex.Unlock()
	if node, ok := lruCache.cache[key]; ok {
			lruCache.removeNode(node)
			delete(lruCache.cache, node.key)
			return true
	}
	return false
}

func (lruCache *LRUCache) ToArray() []map[string]string {
	arr := []map[string]string{}
	curr := lruCache.head.next
	for curr != lruCache.tail {
		arr = append(arr, map[string]string{curr.key: curr.value})
		curr = curr.next
	}
	return arr
}