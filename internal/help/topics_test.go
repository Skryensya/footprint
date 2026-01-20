package help

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupTopic_Existing(t *testing.T) {
	topic := LookupTopic("overview")

	require.NotNil(t, topic)
	require.Equal(t, "overview", topic.Name)
	require.NotEmpty(t, topic.Summary)
}

func TestLookupTopic_NotFound(t *testing.T) {
	topic := LookupTopic("nonexistent")

	require.Nil(t, topic)
}

func TestAllTopics(t *testing.T) {
	topics := AllTopics()

	require.NotEmpty(t, topics)
	require.GreaterOrEqual(t, len(topics), 4) // At least overview, workflow, hooks, data

	// Verify first topic is overview (based on TopicOrder)
	require.Equal(t, "overview", topics[0].Name)
}

func TestTopic_Content(t *testing.T) {
	topic := LookupTopic("overview")
	require.NotNil(t, topic)

	content := topic.Content()

	require.NotEmpty(t, content)
}

func TestTopicOrder(t *testing.T) {
	require.NotEmpty(t, TopicOrder)
	require.Contains(t, TopicOrder, "overview")
	require.Contains(t, TopicOrder, "workflow")
	require.Contains(t, TopicOrder, "hooks")
	require.Contains(t, TopicOrder, "data")
}

func TestAllTopicsMatchesTopicOrder(t *testing.T) {
	topics := AllTopics()

	// Verify that topics are returned in the defined order
	for i, topic := range topics {
		if i < len(TopicOrder) {
			require.Equal(t, TopicOrder[i], topic.Name,
				"Topic at position %d should be %s", i, TopicOrder[i])
		}
	}
}

func TestLookupTopic_AllDefinedTopics(t *testing.T) {
	// Test that all topics in TopicOrder can be looked up
	for _, name := range TopicOrder {
		topic := LookupTopic(name)
		require.NotNil(t, topic, "Topic %q should exist", name)
		require.Equal(t, name, topic.Name)
		require.NotEmpty(t, topic.Summary, "Topic %q should have a summary", name)
	}
}
