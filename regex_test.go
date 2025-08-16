package main

import "testing"

func TestForwardToFormatting(t *testing.T) {

	testStrings := map[string]string{
		"http://namespace.net/queue.1": "queue.1",
		"https://namespace.net/queue.2": "queue.2",
		"https://www.namespace.net/queue.3": "queue.3",
		"https://www.namespace.domain.net/queue.4": "queue.4",
		"https://www.name-space.co.uk/topic.1/subscription.1": "topic.1/subscription.1",
		"topic.2/subscription.2": "topic.2/subscription.2",
		"queue.4": "queue.4",
	}

	for k, v := range testStrings {
		val := ensureForwardToFormatting(&k)
		if *val != v {
			t.Errorf("Expected %s but was %s", v, *val)
		}
	}
}
