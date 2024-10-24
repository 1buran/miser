package miser

import (
	"fmt"
	"sync"
)

// TagMap ties together a tag and some object: account, transaction etc.
type TagMap struct {
	Tag, Item ID
	Deleted   bool
}

func (tgm TagMap) Key() string { return fmt.Sprintf("%s-%s", tgm.Tag, tgm.Item) }

type TagMapRegistry struct {
	items  map[string]TagMap
	queued map[string]TagMap

	sync.RWMutex
}

func (tm *TagMapRegistry) Create(tagID, itemID ID) {
	t := TagMap{Tag: tagID, Item: itemID}
	tm.Add(t)
	tm.AddQueued(t)
}

func (tm *TagMapRegistry) Add(t TagMap) int {
	tm.Lock()
	defer tm.Unlock()

	tm.items[t.Key()] = t
	return 1
}

func (tm *TagMapRegistry) AddQueued(t TagMap) {
	tm.Lock()
	defer tm.Unlock()

	tm.queued[t.Key()] = t
}

func (tm *TagMapRegistry) SyncQueued() (changes []TagMap) {
	tm.RLock()
	defer tm.RUnlock()

	for _, t := range tm.queued {
		changes = append(changes, t)
	}
	return
}

func CreateTagsMapRegistry() *TagMapRegistry {
	return &TagMapRegistry{
		items:  make(map[string]TagMap),
		queued: make(map[string]TagMap),
	}
}

func (tm *TagMapRegistry) Load() (int, error) { return Load(tm, TAGS_MAPPING_FILE) }
func (tm *TagMapRegistry) Save() (int, error) { return Save(tm, TAGS_MAPPING_FILE) }

func (tm *TagMapRegistry) Items(tagID ID) (items []ID) {
	tm.RLock()
	defer tm.RUnlock()

	for _, v := range tm.items {
		if v.Tag == tagID {
			items = append(items, v.Item)
		}
	}
	return
}

func (tm *TagMapRegistry) Tags(itemID ID) (tags []ID) {
	tm.RLock()
	defer tm.RUnlock()

	for _, v := range tm.items {
		if v.Item == itemID {
			tags = append(tags, v.Tag)
		}
	}
	return
}
