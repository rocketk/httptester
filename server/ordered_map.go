package server

import "sync"

type OrderedUserMap struct {
	ids   []string
	m     map[string]User
	mutex *sync.Mutex
}

func NewOrderedUserMap(cap int) OrderedUserMap {
	oum := OrderedUserMap{
		ids:   make([]string, 0, cap),
		m:     make(map[string]User, cap),
		mutex: &sync.Mutex{},
	}
	return oum
}

func (m *OrderedUserMap) Put(id string, user User) {
	m.mutex.Lock()
	if _, ok := m.m[id]; !ok {
		m.ids = append(m.ids, id)
	}
	m.m[id] = user
	m.mutex.Unlock()
}

func (m *OrderedUserMap) Get(id string) (User, bool) {
	u, ok := m.m[id]
	return u, ok
}

func (m *OrderedUserMap) GetAllOrdered() []User {
	list := make([]User, len(m.ids))
	for idx, id := range m.ids {
		list[idx] = m.m[id]
	}
	return list
}

func (m *OrderedUserMap) Delete(id string) {
	m.mutex.Lock()
	if _, ok := m.m[id]; !ok {
		return
	}
	idx := indexOf(m.ids, id)
	m.ids = append(m.ids[:idx], m.ids[idx+1:]...)
	delete(m.m, id)
	m.mutex.Unlock()
}

func (m *OrderedUserMap) Size() int {
	return len(m.ids)
}

func indexOf(slice []string, query string) int {
	for idx, item := range slice {
		if item == query {
			return idx
		}
	}
	return -1
}
