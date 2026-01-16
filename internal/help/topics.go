package help

import (
	"embed"
)

//go:embed topics/*.txt
var topicsFS embed.FS

// Topic represents a conceptual documentation topic (not a command).
// Topics are static text embedded in the binary via go:embed.
type Topic struct {
	Name    string
	Summary string
	content string // loaded lazily from embedded files
}

// Content returns the topic's full documentation text.
func (t *Topic) Content() string {
	return t.content
}

// topics is the registry of all available help topics.
var topics = map[string]*Topic{
	"overview": {
		Name:    "overview",
		Summary: "What fp is and what it does",
	},
	"workflow": {
		Name:    "workflow",
		Summary: "Typical daily usage and mental model",
	},
	"hooks": {
		Name:    "hooks",
		Summary: "How fp integrates with git hooks",
	},
	"data": {
		Name:    "data",
		Summary: "What data fp records and how to interpret it",
	},
}

// TopicOrder defines the display order for topic listings.
var TopicOrder = []string{"overview", "workflow", "hooks", "data"}

func init() {
	// Load topic content from embedded files at startup.
	for name, topic := range topics {
		data, err := topicsFS.ReadFile("topics/" + name + ".txt")
		if err != nil {
			// Should never happen if files are properly embedded
			topic.content = "Error loading topic: " + err.Error()
			continue
		}
		topic.content = string(data)
	}
}

// LookupTopic returns a topic by name, or nil if not found.
func LookupTopic(name string) *Topic {
	return topics[name]
}

// AllTopics returns all topics in display order.
func AllTopics() []*Topic {
	result := make([]*Topic, 0, len(TopicOrder))
	for _, name := range TopicOrder {
		if t := topics[name]; t != nil {
			result = append(result, t)
		}
	}
	return result
}
