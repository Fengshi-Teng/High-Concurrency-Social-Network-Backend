import matplotlib.pyplot as plt
import numpy as np

def read_seq_times(file_path, sizes, rounds=5):
    with open(file_path, "r") as f:
        lines = [float(x.strip()) for x in f if x.strip()]
    expected = len(sizes) * rounds
    if len(lines) != expected:
        print(f"[Warning] seq_results length {len(lines)} != expected {expected}")

    seq_times = []
    for i in range(len(sizes)):
        start = i * rounds
        end = start + rounds
        avg = np.mean(lines[start:end])
        seq_times.append(avg)
    return seq_times


def read_par_times(file_path, sizes, threads, rounds=5):
    with open(file_path, "r") as f:
        lines = [float(x.strip()) for x in f if x.strip()]
    expected = len(sizes) * len(threads) * rounds
    if len(lines) != expected:
        print(f"[Warning] par_results length {len(lines)} != expected {expected}")

    par_times = []
    idx = 0
    for _ in sizes:
        size_times = []
        for _ in threads:
            chunk = lines[idx: idx + rounds]
            avg = np.mean(chunk) if chunk else None
            size_times.append(avg)
            idx += rounds
        par_times.append(size_times)
    return par_times


def generate_speedup_graph(sizes, seq_times, par_times, threads, outfile="speedup2.png"):
    plt.figure(figsize=(8, 5))
    for i, size in enumerate(sizes):
        seq_time = seq_times[i]
        speedups = [seq_time / t if t else 0 for t in par_times[i]]
        plt.plot(threads, speedups, marker="o", label=size)

    plt.xlabel("Number of Threads")
    plt.ylabel("Speedup (T_seq / T_par)")
    plt.title("Speedup vs Threads for Different Input Sizes")
    plt.legend()
    plt.grid(True)
    plt.ylim(0, None)
    plt.savefig(outfile, dpi=200)
    print(f"✅ Saved plot to {outfile}")


def main():
    threads = [2, 4, 6, 8, 12]
    sizes = ["xsmall", "small", "medium", "large", "xlarge"]
    rounds = 5

    seq_times = read_seq_times("seq_results.txt", sizes, rounds)
    par_times = read_par_times("par_results.txt", sizes, threads, rounds)

    generate_speedup_graph(sizes, seq_times, par_times, threads, outfile="speedup.png")


if __name__ == "__main__":
    main()
