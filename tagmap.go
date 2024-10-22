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

func (tmr *TagMapRegistry) Create(tagID, itemID ID) {
	tm := TagMap{Tag: tagID, Item: itemID}
	tmr.Add(tm)
	tmr.AddQueued(tm)
}

func (tmr *TagMapRegistry) Add(tm TagMap) int {
	tmr.Lock()
	defer tmr.Unlock()

	tmr.items = append(tmr.items, tm)
	return 1
}

func (tmr *TagMapRegistry) AddQueued(tm TagMap) {
	tmr.Lock()
	defer tmr.Unlock()

	tmr.queued = append(tmr.queued, tm)
}

func (tmr *TagMapRegistry) SyncQueued() []TagMap {
	tmr.RLock()
	defer tmr.RUnlock()

	return tmr.queued
}

func CreateTagsMapRegistry() *TagMapRegistry   { return &TagMapRegistry{} }
func (tmr *TagMapRegistry) Load() (int, error) { return Load(tmr, TAGS_MAPPING_FILE) }
func (tmr *TagMapRegistry) Save() (int, error) { return Save(tmr, TAGS_MAPPING_FILE) }

func (tmr *TagMapRegistry) Items(tagID ID) (items []ID) {
	tmr.RLock()
	defer tmr.RUnlock()

	for _, v := range tmr.items {
		if v.Tag == tagID {
			items = append(items, v.Item)
		}
	}
	return
}

func (tmr *TagMapRegistry) Tags(itemID ID) (tags []ID) {
	tmr.RLock()
	defer tmr.RUnlock()

	for _, v := range tmr.items {
		if v.Item == itemID {
			tags = append(tags, v.Tag)
		}
	}
	return
}
