package miser

import (
	"sync"
)

// Special system tag names:
const (
	Initial    = "Initial"
	Unexpected = "Unexpected"

	// todo add analysis of transactions during the load and mark some transactions as:
	OverAverage = "OverAverage"
	Periodic    = "Periodic"
)

// Tag ties name of a tag and its id, nothing more.
// For tagging object use TagMap.
type Tag struct {
	ID      ID
	Name    EncryptedString
	Deleted bool
}

type TagRegistry struct {
	items  map[ID]Tag
	queued map[ID]Tag

	sync.RWMutex
}

func (tg *TagRegistry) GetById(tagID ID) *Tag {
	tg.RLock()
	defer tg.RUnlock()

	t, ok := tg.items[tagID]
	if ok {
		return &t
	}
	return nil
}

func (tg *TagRegistry) GetByName(n string) *Tag {
	tg.RLock()
	defer tg.RUnlock()

	for _, t := range tg.items {
		if string(t.Name) == n {
			return &t
		}
	}
	return nil
}

func (tg *TagRegistry) Create(n string) *Tag {
	t := Tag{ID: CreateID(), Name: EncryptedString(n)}
	tg.Add(t)
	tg.AddQueued(t)
	return &t
}

// List all tags.
func (tg *TagRegistry) List() map[ID]Tag {
	tg.RLock()
	defer tg.RUnlock()

	return tg.items
}

func (tg *TagRegistry) Add(t Tag) int {
	tg.Lock()
	defer tg.Unlock()

	tg.items[t.ID] = t
	return 1
}

func (tg *TagRegistry) AddQueued(t Tag) {
	tg.Lock()
	defer tg.Unlock()

	tg.queued[t.ID] = t
}

func (tg *TagRegistry) SyncQueued() (changes []Tag) {
	tg.RLock()
	defer tg.RUnlock()

	for _, t := range tg.queued {
		changes = append(changes, t)
	}
	return
}

func CreateTagRegistry() *TagRegistry {
	return &TagRegistry{
		items:  make(map[ID]Tag),
		queued: make(map[ID]Tag),
	}
}

func (tg *TagRegistry) Load() (int, error) { return Load(tg, TAGS_FILE) }
func (tg *TagRegistry) Save() (int, error) { return Save(tg, TAGS_FILE) }
