package msgrouter

import "sync"

type Msgrouter[T any] struct {
	lock     sync.RWMutex
	revivers map[string][]chan (T)
}

func New[T any]() *Msgrouter[T] {
	return &Msgrouter[T]{
		revivers: make(map[string][]chan (T)),
	}
}

func (m *Msgrouter[T]) Sub(key string, toadd chan (T)) {

	m.lock.Lock()
	defer m.lock.Unlock()
	if m.revivers == nil {
		m.revivers = make(map[string][]chan (T))
	}

	if reviverchanlist, ok := m.revivers[key]; ok {
		reviverchanlist = append(reviverchanlist, toadd)
		m.revivers[key] = reviverchanlist

	} else {
		reviverchanlist = make([]chan (T), 1)
		reviverchanlist[0] = toadd
		m.revivers[key] = reviverchanlist

	}
}

func (m *Msgrouter[T]) Pub(key string, msg T) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if reviverchanlist, ok := m.revivers[key]; ok {
		for _, reviverchan := range reviverchanlist {

			reviverchan <- msg

		}
	}
}

func (m *Msgrouter[T]) Unsub(key string, currentchan chan (T)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if reviverchanlist, ok := m.revivers[key]; ok {
		for index, reviverchan := range reviverchanlist {
			if reviverchan == currentchan {
				reviverchanlist = append(reviverchanlist[:index], reviverchanlist[index+1:]...)
				if len(reviverchan) == 0 {
					delete(m.revivers, key)
				} else {
					m.revivers[key] = reviverchanlist
				}
			}
		}
	}
}
