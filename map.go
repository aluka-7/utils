package utils

import "sync"

type SyncMap struct {
	sync.Map
	sync.RWMutex
	len int64
}

func (m *SyncMap) MyStore(key, value interface{}) {
	m.Store(key, value)
	m.Lock()
	defer m.Unlock()
	m.len++
}
func (m *SyncMap) MyLoad(key interface{}) (value interface{}, ok bool) {
	return m.Load(key)
}
func (m *SyncMap) MyLoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	return m.LoadOrStore(key, value)
}
func (m *SyncMap) MyDelete(key interface{}) {
	m.Delete(key)
	m.Lock()
	defer m.Unlock()
	m.len--
}
func (m *SyncMap) MyRange(f func(key, value interface{}) bool) {
	m.Range(f)
}
func (m *SyncMap) MyLen() int64 {
	m.RLock()
	defer m.RUnlock()
	return m.len
}
func ClearSyncMap(m *SyncMap) {
	if m != nil {
		m.MyRange(func(key, value interface{}) bool {
			m.MyDelete(key)
			return true
		})
	}
}
