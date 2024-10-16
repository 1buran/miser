package miser

import (
	"sync"
)

const (
	AccountTag = iota
	BalanceTag
	TransactionTag
)

// Special system tag names:
const (
	Initial    = "Initial"
	Unexpected = "Unexpected"

	// todo add analysis of transactions during the load and mark some transactions as:
	OverAverage = "OverAverage"
	Periodic    = "Periodic"
)

type Tag struct {
	ID      ID
	Name    EncryptedString
	Deleted bool
}

type TagRegistry struct {
	Items  []Tag
	Queued []Tag

	sync.RWMutex
}

func (tr *TagRegistry) GetById(tagID ID) *Tag {
	tr.RLock()
	defer tr.RUnlock()

	for _, v := range tr.Items {
		if v.ID == tagID {
			return &v
		}
	}
	return nil
}

func (tr *TagRegistry) GetByName(n string) *Tag {
	tr.RLock()
	defer tr.RUnlock()

	for _, v := range tr.Items {
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
func (tr *TagRegistry) List() []Tag {
	tr.RLock()
	defer tr.RUnlock()

	return tr.Items
}

func (tr *TagRegistry) Add(t Tag) int {
	tr.Lock()
	defer tr.Unlock()

	tr.Items = append(tr.Items, t)
	return 1
}

func (tr *TagRegistry) AddQueued(t Tag) {
	tr.Lock()
	defer tr.Unlock()

	tr.Queued = append(tr.Queued, t)
}

func (tr *TagRegistry) SyncQueued() []Tag {
	tr.RLock()
	defer tr.RUnlock()

	return tr.Queued
}

var Tags = TagRegistry{}

func LoadTags() (int, error) { return Load(&Tags, TAGS_FILE) }
func SaveTags() (int, error) { return Save(&Tags, TAGS_FILE) }

type TagMap struct {
	Tag, Item ID
	Type      int
}

type TagMapRegistry struct {
	Items  []TagMap
	Queued []TagMap

	sync.RWMutex
}

func (tmr *TagMapRegistry) Create(tagID, itemID ID, t int) {
	tm := TagMap{Tag: tagID, Item: itemID, Type: t}
	tmr.Add(tm)
	tmr.AddQueued(tm)
}

func (tmr *TagMapRegistry) Add(tm TagMap) int {
	tmr.Lock()
	defer tmr.Unlock()

	tmr.Items = append(tmr.Items, tm)
	return 1
}

func (tmr *TagMapRegistry) AddQueued(tm TagMap) {
	tmr.Lock()
	defer tmr.Unlock()

	tmr.Queued = append(tmr.Queued, tm)
}

func (tmr *TagMapRegistry) SyncQueued() []TagMap {
	tmr.RLock()
	defer tmr.RUnlock()

	return tmr.Queued
}

var TagsMap = TagMapRegistry{}

func LoadTagsMap() (int, error) { return Load(&TagsMap, TAGS_MAPPING_FILE) }
func SaveTagsMap() (int, error) { return Save(&TagsMap, TAGS_MAPPING_FILE) }

func (tmr *TagMapRegistry) GetByTagId(tagID ID) (tmList []TagMap) {
	tmr.RLock()
	defer tmr.RUnlock()

	for _, v := range tmr.Items {
		if v.Tag == tagID {
			tmList = append(tmList, v)
		}
	}
	return
}

func (tmr *TagMapRegistry) GetByItemId(itemID ID) (tmList []TagMap) {
	tmr.RLock()
	defer tmr.RUnlock()

	for _, v := range tmr.Items {
		if v.Item == itemID {
			tmList = append(tmList, v)
		}
	}
	return
}

func (tmr *TagMapRegistry) listItems(id ID, t int) (ids []ID) {
	tmr.RLock()
	defer tmr.RUnlock()

	for _, v := range tmr.Items {
		if v.Tag == id && v.Type == t {
			ids = append(ids, v.Item)
		}
	}
	return
}

func (tmr *TagMapRegistry) Accounts(tagID ID) []ID {
	return tmr.listItems(tagID, AccountTag)
}
func (tmr *TagMapRegistry) Transactions(tagID ID) []ID {
	return tmr.listItems(tagID, TransactionTag)
}
