package server

import (
	"encoding/json"
	"fmt"
	"proj2/feed"
	"proj2/queue"
	"sync"
	"sync/atomic"
)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
}

type SharedContext struct {
	config *Config
	mu     *sync.Mutex
	outMu  *sync.Mutex
	cd     *sync.Cond
	q      *queue.LockFreeQueue
	wg     *sync.WaitGroup
	fd     feed.Feed
	closed atomic.Bool
}

// Run starts up the twitter server based on the configuration
// information provided and only returns when the server is fully
// shutdown.
func Run(config Config) {
	mu := &sync.Mutex{}
	outMu := &sync.Mutex{}
	cd := sync.NewCond(mu)
	ctx := &SharedContext{
		config: &config,
		mu:     mu,
		outMu:  outMu,
		cd:     cd,
		q:      queue.NewLockFreeQueue(),
		wg:     &sync.WaitGroup{},
		fd:     feed.NewFeed(),
	}

	if config.Mode == "s" {
		config.ConsumersCount = 1
	}

	for i := 1; i <= config.ConsumersCount; i++ {
		ctx.wg.Add(1)
		go consumer(ctx)
	}
	producer(ctx)
	ctx.wg.Wait()
}

// Producer&Consumer Model
func consumer(ctx *SharedContext) {
	defer ctx.wg.Done()
	// fmt.Fprintln(os.Stderr, ">>> Starting consumer...")
	config := ctx.config
	mu := ctx.mu
	cd := ctx.cd
	q := ctx.q
	fd := ctx.fd

	for {
		if ctx.closed.Load() && q.IsEmpty() {
			return
		}

		var req *queue.Request
		// fmt.Fprintln(os.Stderr, req)
		for {
			req = q.Dequeue()
			if req != nil {
				break
			}
			mu.Lock()
			cd.Wait()
			mu.Unlock()

			if ctx.closed.Load() && q.IsEmpty() {
				return
			}
		}

		switch req.Command {
		case "ADD":
			Add(config, *req, fd, ctx.outMu)
		case "REMOVE":
			Remove(config, *req, fd, ctx.outMu)
		case "CONTAINS":
			Contains(config, *req, fd, ctx.outMu)
		case "FEED":
			Feed(config, *req, fd, ctx.outMu)
		}
	}
}

func producer(ctx *SharedContext) {
	// fmt.Fprintln(os.Stderr, ">>> Starting producer...")
	for {
		var req map[string]interface{}
		err := ctx.config.Decoder.Decode(&req)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Decode error:", err)
			}
			return
		}

		newReq := &queue.Request{}

		cmd, _ := req["command"].(string)
		newReq.Command = cmd

		switch cmd {
		case "ADD":
			if v, ok := req["id"].(float64); ok {
				newReq.Id = int(v)
			}
			if v, ok := req["body"].(string); ok {
				newReq.Body = v
			}
			if v, ok := req["timestamp"].(float64); ok {
				newReq.TimeStamp = v
			}

		case "REMOVE":
			if v, ok := req["id"].(float64); ok {
				newReq.Id = int(v)
			}
			if v, ok := req["timestamp"].(float64); ok {
				newReq.TimeStamp = v
			}

		case "CONTAINS":
			if v, ok := req["id"].(float64); ok {
				newReq.Id = int(v)
			}
			if v, ok := req["timestamp"].(float64); ok {
				newReq.TimeStamp = v
			}

		case "FEED":
			if v, ok := req["id"].(float64); ok {
				newReq.Id = int(v)
			}

		case "DONE":
			ctx.closed.Store(true)
			ctx.mu.Lock()
			ctx.cd.Broadcast()
			ctx.mu.Unlock()
			return
		}

		ctx.q.Enqueue(newReq)
		ctx.mu.Lock()
		ctx.cd.Signal()
		ctx.mu.Unlock()
	}
}

// Woker Funcitons
func Add(config *Config, req queue.Request, fd feed.Feed, outMu *sync.Mutex) {
	id := req.Id
	body := req.Body
	timestamp := req.TimeStamp

	fd.Add(body, timestamp)

	resp := map[string]interface{}{
		"success": true,
		"id":      id,
	}
	outMu.Lock()
	config.Encoder.Encode(resp)
	outMu.Unlock()
}

func Remove(config *Config, req queue.Request, fd feed.Feed, outMu *sync.Mutex) {
	id := req.Id
	timestamp := req.TimeStamp

	ok := fd.Remove(timestamp)

	resp := map[string]interface{}{
		"success": ok,
		"id":      id,
	}
	outMu.Lock()
	config.Encoder.Encode(resp)
	outMu.Unlock()
}

func Contains(config *Config, req queue.Request, fd feed.Feed, outMu *sync.Mutex) {
	id := req.Id
	timestamp := req.TimeStamp

	ok := fd.Contains(timestamp)

	resp := map[string]interface{}{
		"success": ok,
		"id":      id,
	}

	outMu.Lock()
	config.Encoder.Encode(resp)
	outMu.Unlock()
}

func Feed(config *Config, req queue.Request, fd feed.Feed, outMu *sync.Mutex) {
	id := req.Id

	allPost := fd.AllPost()

	resp := map[string]interface{}{
		"id":   id,
		"feed": allPost,
	}
	outMu.Lock()
	config.Encoder.Encode(resp)
	outMu.Unlock()
}
