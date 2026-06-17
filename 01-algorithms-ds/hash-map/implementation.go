package hashmap

type linkedList struct {
	next  *linkedList
	key   string
	value any
}

// return new pointer of appended linkedList and a bool for increase list or just replace
func (l *linkedList) Append(key string, value any) (*linkedList, bool) {
	if l == nil {
		return &linkedList{key: key, value: value}, true
	}

	var isIncrease bool
	if l.key == key {
		l.value = value
	} else {
		l.next, isIncrease = l.next.Append(key, value)
	}
	return l, isIncrease
}

// Delete return new linkedList pointer and decrease item or not
func (l *linkedList) Delete(key string) (*linkedList, bool) {
	if l == nil {
		return l, false
	}
	if l.key == key {
		return l.next, true
	} else {
		var isDecrease bool
		l.next, isDecrease = l.next.Delete(key)
		return l, isDecrease
	}
}

type Hash struct {
	sli []*linkedList
	cap int
	len int
}

func NewHash(cap int) Hash {
	return Hash{
		sli: make([]*linkedList, cap),
		cap: cap,
		len: 0,
	}
}

func (h *Hash) hashFunc(key string) int {
	const prime = 10000007
	sum := 0
	for _, char := range key {
		sum = (sum + int(char)*prime) % h.cap
	}
	return sum % h.cap
}

func (h *Hash) Len() int {
	return h.len
}

func (h *Hash) Put(key string, value any) {
	hashedKey := h.hashFunc(key)
	var isIncrease bool
	h.sli[hashedKey], isIncrease = h.sli[hashedKey].Append(key, value)
	if isIncrease {
		h.len++
	}

	h.loadFactorAndResize()
}

func (h *Hash) Get(key string) (any, bool) {
	hashedKey := h.hashFunc(key)
	current := h.sli[hashedKey]
	for current != nil {
		if current.key == key {
			return current.value, true
		}
		current = current.next
	}
	return nil, false
}

func (h *Hash) Delete(key string) {
	hashedKey := h.hashFunc(key)
	var isDecrease bool
	h.sli[hashedKey], isDecrease = h.sli[hashedKey].Delete(key)
	if isDecrease {
		h.len--
	}
}

func (h *Hash) loadFactorAndResize() {
	loadFactorPct := 100 * h.len / h.cap

	if loadFactorPct >= 75 {
		h.resize(h.cap * 2)
	}
}

func (h *Hash) resize(newCap int) {
	currentSli := h.sli
	h.sli = make([]*linkedList, newCap)
	h.cap = newCap
	h.len = 0

	for _, head := range currentSli {
		for ; head != nil; head = head.next {
			h.Put(head.key, head.value)
		}
	}
}
