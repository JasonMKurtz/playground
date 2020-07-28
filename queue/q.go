package main

import "fmt"

type QItem struct {
	value    int
	priority int // 0, 1, 2
}

type Q struct {
	items  []QItem
	hipri  int
	lowpri int
}

func (q *Q) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Q) Add(i QItem) {
	q.items = append(q.items, i)
}

func (q *Q) dequeueHead() (bool, QItem) {
	if len(q.items) == 0 {
		return false, QItem{}
	}

	first := q.items[0]
	q.items = q.items[1:]
	return true, first
}

func (q *Q) dequeueTail() (bool, QItem) {
	if len(q.items) == 0 {
		return false, QItem{}
	}

	latest := q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	return true, latest
}

func (q *Q) insert(index int, item QItem) {
	if len(q.items) == 0 {
		q.Add(item)
		return
	}

	q.items = append(q.items, QItem{})
	copy(q.items[index+1:], q.items[index:])
	q.items[index] = item
}

type LQ struct{ Q }
type FQ struct{ Q }
type PQ struct{ Q }

func (lq *LQ) Dequeue() (bool, QItem) {
	return lq.dequeueTail()
}

func (fq *FQ) Dequeue() (bool, QItem) {
	return fq.dequeueHead()
}

func (pq *PQ) Add(item QItem) {
	if len(pq.items) == 0 {
		pq.items = append(pq.items, item)
		return
	}

	if item.priority > pq.hipri {
		pq.insert(0, item)
		pq.hipri = item.priority
	} else if item.priority == pq.hipri {
		pq.insert(1, item)
	} else if item.priority > pq.lowpri {
		pq.insert(len(pq.items)-1, item)
	} else {
		pq.items = append(pq.items, item)
		pq.lowpri = item.priority
	}
}

func (pq *PQ) Dequeue() (bool, QItem) {
	return pq.dequeueHead()
}

type Queue interface {
	Dequeue() (bool, QItem)
	Add(QItem)
}

func MakeQItems(nums []int) []QItem {
	var items []QItem
	for _, i := range nums {
		pri := 0
		if i > 10 {
			pri = 1
		}

		if i > 100 {
			pri = 2
		}

		items = append(items, QItem{value: i, priority: pri})
	}

	return items
}

func MakeQ(q Queue, items []QItem) {
	for _, i := range items {
		q.Add(i)
	}
}

func TakeItem(q Queue) string {
	ok, item := q.Dequeue()
	if !ok {
		return ""
	}

	return fmt.Sprintf("Item (%d) [Priority: %d]", item.value, item.priority)
}

func main() {
	items := MakeQItems([]int{1, 500, 30, 2, 400})

	var fq FQ
	MakeQ(&fq, items)

	var lq LQ
	MakeQ(&lq, items)

	var pq PQ
	MakeQ(&pq, items)

	fmt.Println("Fifo first:")
	for i := 0; i < 5; i++ {
		fmt.Println(TakeItem(&fq))
	}

	fmt.Println("\nAnd then lifo:")
	for i := 0; i < 5; i++ {
		fmt.Println(TakeItem(&lq))
	}

	fmt.Println("\nAnd now priority:")
	for i := 0; i < 5; i++ {
		fmt.Println(TakeItem(&pq))
	}
}
