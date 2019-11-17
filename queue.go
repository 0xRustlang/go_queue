// (c) 2018,2019 timur mobi

package go_queue

import (
	"sync"
)

type Node struct {
	Value string
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	Nodes       []*Node
	SizeIntern  int		// room
	Head        int
	Tail        int
	CountInterm int		// active elements
	Lock_Mutex  sync.RWMutex
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		Nodes: make([]*Node, size),
		SizeIntern:  size,
	}
}

func (q *Queue) Size() int {
	q.Lock_Mutex.Lock()
	defer q.Lock_Mutex.Unlock()
	return q.SizeIntern
}

func (q *Queue) Count() int {
	q.Lock_Mutex.RLock()
	defer q.Lock_Mutex.RUnlock()
	return q.CountInterm
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	q.PushIntern(n, false)
}
func (q *Queue) PushIntern(n *Node, skiplock bool) {
	if !skiplock {
		q.Lock_Mutex.Lock()
		defer q.Lock_Mutex.Unlock()
	}

	//if log1==nil { log1 = stdlog.GetFromFlags() }	// declarated in tools.audioplayback.go
	//fmt.Printf("Queue Push: q.Head=%d q.Tail=%d q.CountInterm=%d\n",q.Head,q.Tail,q.CountInterm)
	q.CountInterm++
	if q.CountInterm >= q.SizeIntern {
		q.PopOldest(true)
		//log1.Debugf("Queue Push: pop one: q.Head=%d q.Tail=%d q.CountInterm=%d", q.Head, q.Tail, q.CountInterm)
	}
	q.Nodes[q.Tail] = n
	q.Tail = (q.Tail + 1) % len(q.Nodes) // position of the next push object
	//fmt.Printf("Queue Push: q.Head=%d q.Tail=%d q.CountInterm=%d\n",q.Head,q.Tail,q.CountInterm)
	/* print queue
	for i:=q.Head; i!=q.Tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.Nodes[i].Value)
		i++
		if i>=q.SizeIntern { i=0 }
	}
	*/
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) PopOldest(skiplock bool) *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	if !skiplock {
		q.Lock_Mutex.Lock()
		defer q.Lock_Mutex.Unlock()
	}
	if q.CountInterm == 0 {
		return nil
	}

	q.CountInterm--

	node := q.Nodes[q.Head] // Head points to the oldest entry
	q.Head = (q.Head + 1) % len(q.Nodes)

	/*
	fmt.Printf("Queue Pop: q.Tail=%d q.CountInterm=%d\n",q.Tail,q.CountInterm)
	for i:=q.Head; i!=q.Tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.Nodes[i].Value)
		i++
		if i>=q.SizeIntern { i=0 }
	}
	*/
	return node
}

func (q *Queue) PeekOldest(skiplock bool) *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	if !skiplock {
		q.Lock_Mutex.Lock()
		defer q.Lock_Mutex.Unlock()
	}
	if q.CountInterm == 0 {
		return nil
	}

	node := q.Nodes[q.Head] // Head points to the oldest entry

	/*
	fmt.Printf("Queue Pop: q.Tail=%d q.CountInterm=%d\n",q.Tail,q.CountInterm)
	for i:=q.Head; i!=q.Tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.Nodes[i].Value)
		i++
		if i>=q.SizeIntern { i=0 }
	}
	*/
	return node
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	q.Lock_Mutex.Lock()
	defer q.Lock_Mutex.Unlock()
	if q.CountInterm == 0 {
		return nil
	}

	q.CountInterm--

	q.Tail = (q.Tail - 1)
	if q.Tail < 0 {
		q.Tail = len(q.Nodes)-1
	}
	node := q.Nodes[q.Tail]

	/*
	fmt.Printf("Queue Pop: q.Tail=%d q.CountInterm=%d\n",q.Tail,q.CountInterm)
	for i:=q.Head; i!=q.Tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.Nodes[i].Value)
		i++
		if i>=q.SizeIntern { i=0 }
	}
	*/
	return node
}

func (q *Queue) Peek() *Node {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	q.Lock_Mutex.Lock()
	defer q.Lock_Mutex.Unlock()
	if q.CountInterm == 0 {
		return nil
	}

	peek := (q.Tail - 1)
	if peek < 0 {
		peek = len(q.Nodes)-1
	}
	node := q.Nodes[peek]

	/*
	fmt.Printf("Queue Pop: q.Tail=%d q.CountInterm=%d\n",q.Tail,q.CountInterm)
	for i:=q.Head; i!=q.Tail; {
		fmt.Printf("Queue: i=%d element:%s\n",i,q.Nodes[i].Value)
		i++
		if i>=q.SizeIntern { i=0 }
	}
	*/
	return node
}

func (q *Queue) InQueue(search string) bool {
	//if log1==nil { log1 = stdlog.GetFromFlags() }
	q.Lock_Mutex.RLock()
	defer q.Lock_Mutex.RUnlock()
	//fmt.Printf("InQueue: "+search+" q.CountInterm=%d %d\n",q.CountInterm,q.Head)
	if q.CountInterm == 0 {
		//fmt.Printf("InQueue: empty\n")
		return false
	}

	//fmt.Printf("InQueue: q.Head=%d q.Tail=%d\n",q.Head,q.Tail)
	for i := q.Head; i != q.Tail; {
		if search == q.Nodes[i].Value {
			//log1.Debugf("InQueue: "+search+" / "+q.Nodes[i].Value+" [%d] FOUND", i)
			return true
		}
		//fmt.Printf("InQueue: "+search+" / "+q.Nodes[i].Value+" [%d] not found\n",i)
		i++
		if i >= q.SizeIntern {
			i = 0
		}
	}
	return false
}

func (q *Queue) Remove(search string) *Node {
	q.Lock_Mutex.Lock()
	defer q.Lock_Mutex.Unlock()
	//fmt.Printf("InQueue: "+search+" q.CountInterm=%d %d\n",q.CountInterm,q.Head)
	if q.CountInterm == 0 {
		//fmt.Printf("InQueue: empty\n")
		return nil
	}

	//fmt.Printf("InQueue: q.Head=%d q.Tail=%d\n",q.Head,q.Tail)
	for i := q.Head; i != q.Tail; {
		if search == q.Nodes[i].Value {
			//log1.Debugf("InQueue: "+search+" / "+q.Nodes[i].Value+" [%d] FOUND", i)
/*
			copy(q.Nodes[i:], q.Nodes[i+1:])    // Shift a[i+1:] left one index.
			q.Nodes[len(q.Nodes)-1] = Node{}    // Erase last element (write zero value).
			q.Nodes = q.Nodes[:len(q.Nodes)-1]  // Truncate slice.
			q.SizeIntern--
			if q.Head>i {
				q.Head--
			}
			if q.Tail>i {
				q.Tail--
			}
			q.CountInterm--
*/
			newq := NewQueue(q.SizeIntern)
			var returnNode* Node = nil
			oldest := q.PopOldest(true)
			for oldest!=nil {
				if oldest.Value != search {
					newq.PushIntern(&Node{oldest.Value},true)
				} else {
					returnNode = oldest
				}
				oldest = q.PopOldest(true)
			}
			oldest = newq.PopOldest(true)
			for oldest!=nil {
				q.PushIntern(&Node{oldest.Value},true)
				oldest = newq.PopOldest(true)
			}
			return returnNode
		}
		//fmt.Printf("InQueue: "+search+" / "+q.Nodes[i].Value+" [%d] not found\n",i)
		i++
		if i >= q.SizeIntern {
			i = 0
		}
	}
	return nil
}

