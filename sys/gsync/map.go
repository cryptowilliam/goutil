package gsync

import "sync"

func CleanupMap(m *sync.Map)  {
	m.Range(func(key, value interface{}) bool {
		m.Delete(key)
		return true
	})
}
