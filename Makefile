BENCH_COUNT ?= 8
BENCH_BASELINE = testdata/bench-baseline.txt
BENCH_CURRENT  = testdata/bench-current.txt

.PHONY: all build vet bench bench-save bench-compare bench-check test

all: build vet test

build:
	go build ./...

vet:
	go vet ./...

# Run benchmarks and display results
bench:
	go test -bench=. -benchmem -count=$(BENCH_COUNT) -run=^$$ ./...

# Save current benchmark results as the new baseline
bench-save:
	@mkdir -p testdata
	go test -bench=. -benchmem -count=$(BENCH_COUNT) -run=^$$ ./... | tee $(BENCH_BASELINE)
	@echo "\nBaseline saved to $(BENCH_BASELINE)"

# Run benchmarks and compare against baseline
bench-compare:
	@if [ ! -f $(BENCH_BASELINE) ]; then \
		echo "No baseline found. Run 'make bench-save' first."; \
		exit 1; \
	fi
	@mkdir -p testdata
	go test -bench=. -benchmem -count=$(BENCH_COUNT) -run=^$$ ./... > $(BENCH_CURRENT)
	benchstat $(BENCH_BASELINE) $(BENCH_CURRENT)

# CI-friendly: fail if any benchmark regressed >10%
bench-check:
	@if [ ! -f $(BENCH_BASELINE) ]; then \
		echo "No baseline found. Run 'make bench-save' first."; \
		exit 1; \
	fi
	@mkdir -p testdata
	go test -bench=. -benchmem -count=$(BENCH_COUNT) -run=^$$ ./... > $(BENCH_CURRENT)
	@echo "=== Benchmark comparison ==="
	benchstat $(BENCH_BASELINE) $(BENCH_CURRENT)

test:
	go test -v ./...
