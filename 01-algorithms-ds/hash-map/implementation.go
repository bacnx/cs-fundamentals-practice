package hashmap

type linkedList struct {
	next  *linkedList
	key   string
	value any
}

func newLinkedList(key string, value any) *linkedList {
	return &linkedList{
		next:  nil,
		key:   key,
		value: value,
	}
}

func (l *linkedList) Value() (string, any, bool) {
	if l != nil {
		return l.key, l.value, true
	}
	return "", nil, false
}

func (l *linkedList) Next() (*linkedList, bool) {
	if l != nil {
		return l.next, true
	}
	return nil, false
}

func (l *linkedList) Append(key string, value any) *linkedList {
	if l != nil {
		l.next = l.next.Append(key, value)
	} else {
		l = &linkedList{key: key, value: value}
	}
	return l
}

func (l *linkedList) Delete(key string) *linkedList {
	if l == nil {
		return l
	}
	if l.key == key {
		return l.next
	}
	if next, ok := l.Next(); ok && next.key == key {
		l.next = next.next
	}
	return l
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
	for char := range key {
		sum = (sum + char*prime) % h.cap
	}
	return sum % h.cap
}

func (h *Hash) Put(key string, value any) {
	hashedKey := h.hashFunc(key)
	if h.sli[hashedKey] == nil {
		h.sli[hashedKey] = newLinkedList(key, value)
	} else {
		if h.sli[hashedKey].key == key {
			h.sli[hashedKey].value = value
		} else {
			h.sli[hashedKey].Append(key, value)
		}
	}
	h.len++
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
	h.sli[hashedKey] = h.sli[hashedKey].Delete(key)
	h.len--
}
