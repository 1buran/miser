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

func (tr *TagRegistry) GetById(tagID ID) *Tag {
	tr.RLock()
	defer tr.RUnlock()

	for i := len(tr.items) - 1; i >= 0; i-- {
		v := tr.items[i]
		if v.ID == tagID {
			return &v
		}
	}
	return nil
}

func (tr *TagRegistry) GetByName(n string) *Tag {
	tr.RLock()
	defer tr.RUnlock()

	for i := len(tr.items) - 1; i >= 0; i-- {
		v := tr.items[i]
		if string(v.Name) == n {
			return &v
		}
	}
	return nil
}

func (tr *TagRegistry) Create(n string) *Tag {
	t := Tag{ID: CreateID(), Name: EncryptedString(n)}
	tr.Add(t)
	tr.AddQueued(t)
	return &t
}

// List all tags.
func (tr *TagRegistry) List() map[ID]Tag {
	tr.RLock()
	defer tr.RUnlock()

	tags := make(map[ID]Tag)
	for _, tag := range tr.items { // the last readed is the most actual version
		tags[tag.ID] = tag
	}
	return tags
}

func (tr *TagRegistry) Add(t Tag) int {
	tr.Lock()
	defer tr.Unlock()

	tr.items = append(tr.items, t)
	return 1
}

func (tr *TagRegistry) AddQueued(t Tag) {
	tr.Lock()
	defer tr.Unlock()

	tr.queued = append(tr.queued, t)
}

func (tr *TagRegistry) SyncQueued() []Tag {
	tr.RLock()
	defer tr.RUnlock()

	return tr.queued
}

func CreateTagRegistry() *TagRegistry      { return &TagRegistry{} }
func (tr *TagRegistry) Load() (int, error) { return Load(tr, TAGS_FILE) }
func (tr *TagRegistry) Save() (int, error) { return Save(tr, TAGS_FILE) }
