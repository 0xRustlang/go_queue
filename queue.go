// (c) 2018 timur mobi

package go_queue

import (
	"sync"
)

type Node struct {
	Value string
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes       []*Node
	size        int
	head        int
	tail        int
	count       int
	lock_Mutex  sync.RWMutex
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

func (q *Queue) Size() int {
	q.lock_Mutex.Lock()
	defer q.lock_Mutex.Unlock()
	return q.size
}

func (q *Queue) Count() int {
	q.lock_Mutex.RLock()
	defer q.lock_Mutex.RUnlock()
	return q.count
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	q.lock_Mutex.Lock()
	defer q.lock_Mutex.Unlock()

	//if log1==nil { log1 = stdlog.GetFromFlags() }	// declarated in tools.audioplayback.go
	//fmt.Printf("Queue Push: q.head=%d q.tail=%d q.count=%d\n",q.head,q.tail,q.count)
	q.count++
	if q.count >= q.size {
		q.PopOldest(true)
		//log1.Debugf("Queue Push: pop one: q.head=%d q.tail=%d q.count=%d", q.head, q.tail, q.count)
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes) // position of the next push object
	//fmt.Printf("Queue Push: q.head=%d q.tail=%d q.count=%d\n",q.head,q.tail,q.count)
	/* print queue
	for i:=q.head; i!=q.tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.nodes[i].Value)
		i++
		if i>=q.size { i=0 }
	}
	*/
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) PopOldest(skiplock bool) *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	if !skiplock {
		q.lock_Mutex.Lock()
		defer q.lock_Mutex.Unlock()
	}
	if q.count == 0 {
		return nil
	}

	q.count--

	node := q.nodes[q.head] // head points to the oldest entry
	q.head = (q.head + 1) % len(q.nodes)

	/*
	fmt.Printf("Queue Pop: q.tail=%d q.count=%d\n",q.tail,q.count)
	for i:=q.head; i!=q.tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.nodes[i].Value)
		i++
		if i>=q.size { i=0 }
	}
	*/
	return node
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	q.lock_Mutex.Lock()
	defer q.lock_Mutex.Unlock()
	if q.count == 0 {
		return nil
	}

	q.count--

	q.tail = (q.tail - 1)
	if q.tail < 0 {
		q.tail = len(q.nodes)-1
	}
	node := q.nodes[q.tail]

	/*
	fmt.Printf("Queue Pop: q.tail=%d q.count=%d\n",q.tail,q.count)
	for i:=q.head; i!=q.tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.nodes[i].Value)
		i++
		if i>=q.size { i=0 }
	}
	*/
	return node
}

func (q *Queue) InQueue(search string) bool {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	q.lock_Mutex.RLock()
	defer q.lock_Mutex.RUnlock()
	//fmt.Printf("InQueue: "+search+" q.count=%d %d\n",q.count,q.head)
	if q.count == 0 {
		//fmt.Printf("InQueue: empty\n")
		return false
	}

	//fmt.Printf("InQueue: q.head=%d q.tail=%d\n",q.head,q.tail)
	for i := q.head; i != q.tail; {
		if search == q.nodes[i].Value {
			//log1.Debugf("InQueue: "+search+" / "+q.nodes[i].Value+" [%d] FOUND", i)
			return true
		}
		//fmt.Printf("InQueue: "+search+" / "+q.nodes[i].Value+" [%d] not found\n",i)
		i++
		if i >= q.size {
			i = 0
		}
	}
	return false
}

