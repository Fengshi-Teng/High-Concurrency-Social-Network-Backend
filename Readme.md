# Concurrent Twitter-Style Server (Go)

A high-performance, concurrent Twitter-like server implemented in Go. This system processes client requests such as `ADD`, `REMOVE`, and `CONTAINS` using a multi-threaded architecture. It features a custom **Producer-Consumer model**, **Lock-free Task Queues**, and **Read-Write synchronization** to handle high-concurrency workloads.

## Key Features

* **Multi-Mode Execution**: Supports both sequential and parallel execution modes for comparative performance benchmarking.


* **Thread-Safe Feed**: Implements a user news feed as a linked list with thread-safe operations.


* **Custom Synchronization**: Utilizes a custom Read-Write Lock implemented using `sync.Cond` and `sync.Mutex` with a limit of 32 concurrent readers.
* **Lock-Free Task Queue**: Employs an unbounded, non-blocking queue using atomic CAS (Compare-And-Swap) operations for task distribution.
* **JSON Protocol**: Communicates via a streaming JSON-encoded request/response protocol over `os.Stdin` and `os.Stdout`.

## Architecture

The system mimics a real-life client-server model where requests are streamed to a server that handles tasks asynchronously.

### Component Breakdown

* **`twitter.go`**: The entry point that configures the server, manages the feed, and initiates the run.


* **`server.go`**: Contains the main logic for spawning consumer goroutines (workers) and the producer function that populates the task queue.


* **`feed.go`**: The core data structure. A linked list of posts protected by synchronization primitives to ensure correctness under concurrent access.



## Performance Analysis

The project includes a benchmarking suite to evaluate speedup across various input sizes: `xsmall`, `small`, `medium`, `large`, and `xlarge`.

### Analysis Insights

* **Scaling**: Speedup improves significantly for larger inputs as the overhead of synchronization is amortized over more computation.


* **Peak Efficiency**: Performance typically peaks at **6-8 threads** before saturation or contention causes a drop in speedup.


* **Contention**: At high thread counts (e.g., 12 threads), performance may degrade due to high contention in CAS retry loops within the lock-free queue.


* **Hardware Impact**: Performance is heavily influenced by CPU frequency and cache hierarchy, as linked-list operations are sensitive to pointer-chasing and cache misses.



## Getting Started

### Prerequisites

* Go (Golang) 1.18 or higher.
* A Linux/Unix-based environment (required for Slurm benchmarking scripts).



### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd proj2

```


2. Run the correctness tests:
```bash
cd grader
go run proj2/grader proj2

```



### Running the Server

To run the server in parallel mode with 4 consumer threads:

```bash
go run twitter/twitter.go 4 < tasks.txt

```

## Benchmarking

To replicate the performance analysis and generate the speedup graphs:

1. Navigate to the benchmark directory:
```bash
cd proj2/benchmark

```


2. Execute the Slurm script:
```bash
sbatch benchmark-generate_speedup_graph.sh

```



The results will produce a `speedup.png` file for analysis.

## 🔍 Future Optimizations

* **Work-Stealing**: Implementing per-thread queues to reduce cache invalidation and pointer contention.


* **Fine-Grained Locking**: Moving away from a coarse-grained RWLock to a lock-free linked list or skip-list for the feed.


* **Hybrid Spinning**: Implementing hybrid spin-locks with backoff to reduce latency in the producer/consumer signaling.



---

**Would you like me to generate a specific `docker-compose.yml` file so you can run this entire benchmarking environment in a container?**