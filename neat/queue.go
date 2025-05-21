package neat

import "errors"


type Queue struct {
    Elements []int
}

func (q *Queue) Enqueue(elem int) {
    q.Elements = append(q.Elements, elem)
}

func (q *Queue) Dequeue() int {
    if q.IsEmpty() {
        return 0
    }
    element := q.Elements[0]
    if q.GetLength() == 1 {
        q.Elements = nil
        return element
    }
    q.Elements = q.Elements[1:]
    return element // Slice off the element once it is dequeued.
}

func (q *Queue) GetLength() int {
    return len(q.Elements)
}

func (q *Queue) IsEmpty() bool {
    return len(q.Elements) == 0
}

func (q *Queue) Peek() (int, error) {
    if q.IsEmpty() {
        return 0, errors.New("empty queue")
    }
    return q.Elements[0], nil
}
