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
	items  []Tag
	queued []Tag

	sync.RWMutex
}

func (tg *TagRegistry) GetById(tagID ID) *Tag {
	tg.RLock()
	defer tg.RUnlock()

	for i := len(tg.items) - 1; i >= 0; i-- {
		v := tg.items[i]
		if v.ID == tagID {
			return &v
		}
	}
	return nil
}

func (tg *TagRegistry) GetByName(n string) *Tag {
	tg.RLock()
	defer tg.RUnlock()

	for i := len(tg.items) - 1; i >= 0; i-- {
		v := tg.items[i]
		if string(v.Name) == n {
			return &v
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

	tags := make(map[ID]Tag)
	for _, tag := range tg.items { // the last readed is the most actual version
		tags[tag.ID] = tag
	}
	return tags
}

func (tg *TagRegistry) Add(t Tag) int {
	tg.Lock()
	defer tg.Unlock()

	tg.items = append(tg.items, t)
	return 1
}

func (tg *TagRegistry) AddQueued(t Tag) {
	tg.Lock()
	defer tg.Unlock()

	tg.queued = append(tg.queued, t)
}

func (tg *TagRegistry) SyncQueued() []Tag {
	tg.RLock()
	defer tg.RUnlock()

	return tg.queued
}

func CreateTagRegistry() *TagRegistry      { return &TagRegistry{} }
func (tg *TagRegistry) Load() (int, error) { return Load(tg, TAGS_FILE) }
func (tg *TagRegistry) Save() (int, error) { return Save(tg, TAGS_FILE) }
