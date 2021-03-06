package db

type SubtractNode struct {
	left  Node
	right Node
}

func (s *SubtractNode) String() string {
	return "SubtractNode{}"
}

func (s *SubtractNode) Children() []Node {
	return []Node{s.left, s.right}
}

func (s *SubtractNode) ResultIter(tx *txWithContext, qr *QueryResult) (qrIterator, error) {
	leftIter, err := s.left.ResultIter(tx, qr)
	if err != nil {
		return nil, err
	}
	rightIter, err := s.right.ResultIter(tx, qr)
	if err != nil {
		return nil, err
	}
	return &subtractIterator{
		tx:        tx,
		qr:        qr,
		left:      leftIter,
		right:     rightIter,
		rightRows: map[ID]*QueryResult{},
	}, nil
}

type subtractIterator struct {
	tx        *txWithContext
	qr        *QueryResult
	left      qrIterator
	right     qrIterator
	rightRows map[ID]*QueryResult
	value     *QueryResult
	err       error
}

func (i *subtractIterator) loadLeft() (qr *QueryResult, err error) {
	ok := i.left.Next()
	if ok {
		qr = i.left.Value()
	} else {
		err = i.left.Err()
	}
	return
}

func (i *subtractIterator) attemptMatch(leftQR *QueryResult) (found bool, err error) {
	leftID := leftQR.Record.ID()
	if _, ok := i.rightRows[leftQR.Record.ID()]; ok {
		return true, nil
	}
	for i.right.Next() {
		rightQR := i.right.Value()
		rightID := rightQR.Record.ID()
		i.rightRows[rightID] = rightQR
		if rightID == leftID {
			return true, nil
		}
	}
	if i.right.Err() == Done {
		return false, nil
	}
	return false, i.right.Err()
}

func (i *subtractIterator) Next() bool {
	for {
		leftQR, err := i.loadLeft()
		if err != nil {
			i.err = err
			return false
		}
		match, err := i.attemptMatch(leftQR)
		if err != nil {
			i.err = err
			return false
		}
		if !match {
			i.value = leftQR
			return true
		}
	}
}

func (i *subtractIterator) Value() *QueryResult {
	return i.value
}

func (i *subtractIterator) Err() error {
	return i.err
}
