package tag

import (
	"github.com/df-mc/atomic"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
	"sync"
)

// Tag represents a tag that can be applied to a player. It is used to format the chat messages
type Tag interface {
	Name() string
	Format() string
}

var (
	// tags is a list of all tags that are available.
	tags []Tag
	// tagsByName is a map of all tags that are available indexed by their name.
	tagsByName = map[string]Tag{}
)

// All returns all registered tags.
func All() []Tag {
	return tags
}

// Register registers a tag to the tags list.
func Register(tag Tag) {
	tags = append(tags, tag)
	tagsByName[tag.Name()] = tag
}

// ByName returns the tag with the given name. If no tag with the given name is registered, the second return
// value is false.
func ByName(name string) (Tag, bool) {
	tag, ok := tagsByName[name]
	return tag, ok
}

// Tags represents a list of tags that can be applied to a player.
type Tags struct {
	tagMu sync.Mutex
	tags  []Tag

	active atomic.Value[Tag]
}

// NewTags creates a new Tags instance.
func NewTags(tags []Tag) *Tags {
	return &Tags{
		tags:   tags,
		active: *atomic.NewValue[Tag](nil),
	}
}

// Active returns the active tag of the list of tags.
func (t *Tags) Active() (Tag, bool) {
	tag := t.active.Load()
	return tag, tag != nil
}

// UpdateActive updates the active tag of the list of tags.
func (t *Tags) UpdateActive(tag Tag) {
	t.active.Store(tag)
}

// Add adds a tag to the list of tags.
func (t *Tags) Add(tag Tag) {
	t.tagMu.Lock()
	defer t.tagMu.Unlock()
	t.tags = append(t.tags, tag)
}

// Remove removes a tag from the list of tags.
func (t *Tags) Remove(tag Tag) {
	t.tagMu.Lock()
	defer t.tagMu.Unlock()
	i := slices.IndexFunc(t.tags, func(other Tag) bool {
		return tag == other
	})
	t.tags = slices.Delete(t.tags, i, i+1)
}

// Contains returns true if the list of tags contains the tag provided.
func (t *Tags) Contains(tag Tag) bool {
	t.tagMu.Lock()
	defer t.tagMu.Unlock()
	return slices.Contains(t.tags, tag)
}

// All returns all tags that are currently applied to the list of tags.
func (t *Tags) All() []Tag {
	t.tagMu.Lock()
	defer t.tagMu.Unlock()
	return t.tags
}

// tagsData is a struct that is used to store the tags of a player in a database.
type tagsData struct {
	Tags   []string `json:"tags"`
	Active string   `json:"active"`
}

// MarshalBSON ...
func (t *Tags) MarshalBSON() ([]byte, error) {
	var d tagsData
	t.tagMu.Lock()
	defer t.tagMu.Unlock()

	for _, tag := range t.tags {
		d.Tags = append(d.Tags, tag.Name())
	}

	if tg, active := t.Active(); active {
		d.Active = tg.Name()
	}
	return bson.Marshal(d)
}

// UnmarshalBSON ...
func (t *Tags) UnmarshalBSON(data []byte) error {
	var d tagsData
	if err := bson.Unmarshal(data, &d); err != nil {
		return err
	}

	for _, name := range d.Tags {
		if tag, ok := ByName(name); ok {
			t.Add(tag)
		}
	}

	if tag, ok := ByName(d.Active); ok && t.Contains(tag) {
		t.active = *atomic.NewValue(tag)
	}
	return nil
}
