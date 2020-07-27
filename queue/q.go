package main

import "fmt"

type Q struct {
	items []int
}

func MakeQ(q Queue, nums []int) {
	for _, i := range nums {
		q.Add(i)
	}
}

func (q *Q) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Q) Add(v int) {
	q.items = append(q.items, v)
}

func (q *Q) Items() []int {
	return q.items
}

type LQ struct{ Q }
type FQ struct{ Q }

func (lq *LQ) Dequeue() int {
	if len(lq.items) == 0 {
		return -1
	}

	latest := lq.items[len(lq.items)-1]
	lq.items = lq.items[:len(lq.items)-1]
	return latest
}

func (fq *FQ) Dequeue() int {
	if len(fq.items) == 0 {
		return -1
	}

	first := fq.items[0]
	fq.items = fq.items[1:]
	return first
}

func TakeItem(q Queue) int {
	return q.Dequeue()
}

type Queue interface {
	Dequeue() int
	Add(int)
}

func main() {
	nums := []int{1, 5, 3, 2, 4}

	var fq FQ
	MakeQ(&fq, nums)

	var lq LQ
	MakeQ(&lq, nums)

	fmt.Println("Fifo first:")
	for i := 0; i < 5; i++ {
		fmt.Println(TakeItem(&fq))
	}

	fmt.Println("And then lifo:")
	for i := 0; i < 5; i++ {
		fmt.Println(TakeItem(&lq))
	}
}
