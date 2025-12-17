package feed

import "proj2/lock"

// Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	Add(body string, timestamp float64)
	Remove(timestamp float64) bool
	Contains(timestamp float64) bool
	AllPost() []postContent
}

// feed is the internal representation of a user's twitter feed (hidden from outside packages)
type feed struct {
	start *post // a pointer to the beginning post
	mu    *lock.RWlock
}

// post is the internal representation of a post on a user's twitter feed (hidden from outside packages)
type post struct {
	body      string  // the text of the post
	timestamp float64 // Unix timestamp of the post
	next      *post   // the next post in the feed
}

type postContent struct {
	Body      string  `json:"body"`
	Timestamp float64 `json:"timestamp"`
}

// NewPost creates and returns a new post value given its body and timestamp
func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

// NewFeed creates a empy user feed
func NewFeed() Feed {
	return &feed{
		start: nil,
		mu:    lock.NewRWLock(),
	}
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *feed) Add(body string, timestamp float64) {
	newP := newPost(body, timestamp, nil)

	f.mu.Lock()
	defer f.mu.Unlock()
	// Case 1: Empty feed
	if f.start == nil {
		f.start = newP
		return
	}

	// Case 2: New post is the most recent (insert at front)
	if timestamp > f.start.timestamp {
		newP.next = f.start
		f.start = newP
		return
	}

	// Traverse to find insertion point
	prev := f.start
	cur := f.start.next
	for cur != nil && cur.timestamp > timestamp {
		prev = cur
		cur = cur.next
	}

	// Insert between prev and cur
	prev.next = newP
	newP.next = cur
}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp float64) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.start == nil || timestamp > f.start.timestamp {
		return false
	}

	if timestamp == f.start.timestamp {
		f.start = f.start.next
		return true
	}

	// Traverse to find insertion point
	prev := f.start
	cur := f.start.next
	for cur != nil && cur.timestamp > timestamp {
		prev = cur
		cur = cur.next
	}

	if cur != nil && (cur.timestamp == timestamp) {
		// Remove
		prev.next = cur.next
		return true
	} else {
		return false
	}
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp float64) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for cur := f.start; cur != nil; cur = cur.next {
		if cur.timestamp == timestamp {
			return true
		}
		if cur.timestamp < timestamp {
			return false
		}
	}
	return false
}

func (f *feed) AllPost() []postContent {
	allPost := []postContent{}

	f.mu.RLock()
	defer f.mu.RUnlock()

	for cur := f.start; cur != nil; cur = cur.next {
		allPost = append(allPost, postContent{
			Body:      cur.body,
			Timestamp: cur.timestamp,
		})
	}

	return allPost
}
