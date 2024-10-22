package miser

import (
	"sync"
)

// TagMap ties together a tag and some object: account, transaction etc.
type TagMap struct {
	Tag, Item ID
}

type TagMapRegistry struct {
	items  []TagMap
	queued []TagMap

	sync.RWMutex
}

func (tm *TagMapRegistry) Create(tagID, itemID ID) {
	tm := TagMap{Tag: tagID, Item: itemID}
	tm.Add(tm)
	tm.AddQueued(tm)
}

func (tm *TagMapRegistry) Add(tm TagMap) int {
	tm.Lock()
	defer tm.Unlock()

	tm.items = append(tm.items, tm)
	return 1
}

func (tm *TagMapRegistry) AddQueued(tm TagMap) {
	tm.Lock()
	defer tm.Unlock()

	tm.queued = append(tm.queued, tm)
}

func (tm *TagMapRegistry) SyncQueued() []TagMap {
	tm.RLock()
	defer tm.RUnlock()

	return tm.queued
}

func CreateTagsMapRegistry() *TagMapRegistry  { return &TagMapRegistry{} }
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
