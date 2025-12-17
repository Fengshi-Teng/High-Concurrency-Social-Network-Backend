#!/bin/bash
#
#SBATCH --mail-user=tengfengshi@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=benchmark-generate_speedup_graph
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/tengfengshi/project-2-Fengshi-Teng/proj2/benchmark
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=50:00


# ======== CONFIG ========
THREADS=(2 4 6 8 12)
SIZES=("xsmall" "small" "medium" "large" "xlarge")
ROUNDS=5
GO_CMD="go run benchmark.go"

# =======================
# sequential baseline
for size in "${SIZES[@]}"; do
    for ((i=1; i<=ROUNDS; i++)); do
        $GO_CMD s $size >> seq_results.txt
    done
done

# parallel versions
for size in "${SIZES[@]}"; do
    for t in "${THREADS[@]}"; do
        for ((i=1; i<=ROUNDS; i++)); do
            $GO_CMD p $size $t >> par_results.txt
        done
    done
done

python generate_speedup_graph.py
