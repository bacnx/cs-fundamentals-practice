package hashmap

type hashItemFlag int

const (
	Default hashItemFlag = iota
	Using
	Deleted
)

type hashItem struct {
	key   string
	value any
	flag  hashItemFlag
}

type HashOpenAddressing struct {
	sli []hashItem
	cap int
	len int
}

func NewHashOA(cap int) HashOpenAddressing {
	return HashOpenAddressing{
		sli: make([]hashItem, cap),
		cap: cap,
		len: 0,
	}
}

func (h *HashOpenAddressing) Len() int {
	return h.len
}

func (h *HashOpenAddressing) Put(key string, value any) {
	if (h.len+1)*2 >= h.cap {
		h.doubleCap()
	}

	hashedKey := h.hashFunc(key)
	slotIndex := hashedKey
	for {
		item := h.sli[slotIndex]
		if item.flag != Using || item.key == key {
			break
		}
		slotIndex = (slotIndex + 1) % h.cap
	}

	item := h.sli[slotIndex]
	if item.flag != Using {
		h.sli[slotIndex].key = key
		h.sli[slotIndex].value = value
		h.sli[slotIndex].flag = Using
		h.len++
	}

	if item.key == key {
		h.sli[slotIndex].value = value
	}
}

func (h *HashOpenAddressing) Delete(key string) {
	hashedKey := h.hashFunc(key)

	for {
		item := h.sli[hashedKey]
		if item.flag == Default {
			break
		}
		if item.key == key {
			h.sli[hashedKey].flag = Deleted
			h.len--
			break
		}
		hashedKey = (hashedKey + 1) % h.cap
	}
}

func (h *HashOpenAddressing) Get(key string) (any, bool) {
	hashedKey := h.hashFunc(key)

	for {
		item := h.sli[hashedKey]
		if item.flag == Default {
			return nil, false
		}
		if item.key == key && item.flag == Using {
			return item.value, true
		}
		hashedKey = (hashedKey + 1) % h.cap
	}
}

func (h *HashOpenAddressing) hashFunc(key string) int {
	primeNum := 10000007
	sum := 0
	for char := range key {
		sum = (sum + char*primeNum%h.cap) % h.cap
	}
	return sum % h.cap
}

func (h *HashOpenAddressing) doubleCap() {
	oldSli := h.sli
	h.cap = h.cap * 2
	h.len = 0
	h.sli = make([]hashItem, h.cap)

	for _, item := range oldSli {
		if item.flag == Using {
			h.Put(item.key, item.value)
		}
	}
}
