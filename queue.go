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
	lock_Mutex  sync.Mutex
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

func (q *Queue) Size() int {
	return q.size
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	q.lock_Mutex.Lock()

	//if log1==nil { log1 = stdlog.GetFromFlags() }	// declarated in tools.audioplayback.go
	//fmt.Printf("Queue Push: q.head=%d q.tail=%d q.count=%d\n",q.head,q.tail,q.count)
	q.count++
	if q.count >= q.size {
		q.Pop()
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
	q.lock_Mutex.Unlock()
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) PopOldest() *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	if q.count == 0 {
		return nil
	}

	q.lock_Mutex.Lock()
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
	q.lock_Mutex.Unlock()
	return node
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	if q.count == 0 {
		return nil
	}

	q.lock_Mutex.Lock()
	q.count--

	q.tail = (q.tail - 1)
	if q.tail < 0 {
		q.tail = len(q.nodes)
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
	q.lock_Mutex.Unlock()
	return node
}

func (q *Queue) InQueue(search string) bool {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	//fmt.Printf("InQueue: "+search+" q.count=%d %d\n",q.count,q.head)
	if q.count == 0 {
		//fmt.Printf("InQueue: empty\n")
		return false
	}

	q.lock_Mutex.Lock()
	//fmt.Printf("InQueue: q.head=%d q.tail=%d\n",q.head,q.tail)
	for i := q.head; i != q.tail; {
		if search == q.nodes[i].Value {
			//log1.Debugf("InQueue: "+search+" / "+q.nodes[i].Value+" [%d] FOUND", i)
			q.lock_Mutex.Unlock()
			return true
		}
		//fmt.Printf("InQueue: "+search+" / "+q.nodes[i].Value+" [%d] not found\n",i)
		i++
		if i >= q.size {
			i = 0
		}
	}
	q.lock_Mutex.Unlock()
	return false
}
