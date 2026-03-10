package main

import (
	"math/big"
	"math/rand"
	"sort"
	"testing"
)

// seededRand returns a deterministic *rand.Rand for reproducible tests.
func seededRand() *rand.Rand {
	return rand.New(rand.NewSource(42))
}

func bi(n int64) *big.Int { return big.NewInt(n) }

func bigFromString(s string) *big.Int {
	v, _ := new(big.Int).SetString(s, 10)
	return v
}

// --- fisherYatesSample ---

func TestFisherYatesReturnsCorrectCount(t *testing.T) {
	result := fisherYatesSample(bi(1), 50, 10, seededRand())
	if len(result) != 10 {
		t.Errorf("got %d results, want 10", len(result))
	}
}

func TestFisherYatesAllUnique(t *testing.T) {
	result := fisherYatesSample(bi(0), 100, 50, seededRand())
	seen := make(map[string]struct{}, len(result))
	for _, v := range result {
		key := v.String()
		if _, exists := seen[key]; exists {
			t.Errorf("duplicate value: %s", key)
		}
		seen[key] = struct{}{}
	}
}

func TestFisherYatesValuesInRange(t *testing.T) {
	lower := bi(10)
	upper := bi(20) // inclusive
	count := 5
	result := fisherYatesSample(lower, 11, count, seededRand())
	for _, v := range result {
		if v.Cmp(lower) < 0 || v.Cmp(upper) > 0 {
			t.Errorf("value %s out of range [10, 20]", v)
		}
	}
}

func TestFisherYatesFullRange(t *testing.T) {
	// count == rangeSize: must return every value in [lower, lower+rangeSize)
	lower := bi(5)
	const rangeSize int64 = 20
	result := fisherYatesSample(lower, rangeSize, int(rangeSize), seededRand())
	if len(result) != int(rangeSize) {
		t.Fatalf("got %d results, want %d", len(result), rangeSize)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Cmp(result[j]) < 0 })
	for i, v := range result {
		want := new(big.Int).Add(lower, bi(int64(i)))
		if v.Cmp(want) != 0 {
			t.Errorf("sorted[%d] = %s, want %s", i, v, want)
		}
	}
}

func TestFisherYatesCountOne(t *testing.T) {
	result := fisherYatesSample(bi(7), 100, 1, seededRand())
	if len(result) != 1 {
		t.Fatalf("got %d results, want 1", len(result))
	}
	if result[0].Cmp(bi(7)) < 0 || result[0].Cmp(bi(106)) > 0 {
		t.Errorf("value %s out of range [7, 106]", result[0])
	}
}

func TestFisherYatesNegativeLowerBound(t *testing.T) {
	lower := bi(-50)
	upper := bi(-1) // inclusive; rangeSize = 50
	result := fisherYatesSample(lower, 50, 10, seededRand())
	for _, v := range result {
		if v.Cmp(lower) < 0 || v.Cmp(upper) > 0 {
			t.Errorf("value %s out of range [-50, -1]", v)
		}
	}
}

func TestFisherYatesDeterministic(t *testing.T) {
	r1 := rand.New(rand.NewSource(99))
	r2 := rand.New(rand.NewSource(99))
	a := fisherYatesSample(bi(0), 100, 10, r1)
	b := fisherYatesSample(bi(0), 100, 10, r2)
	for i := range a {
		if a[i].Cmp(b[i]) != 0 {
			t.Errorf("non-deterministic at index %d: %s vs %s", i, a[i], b[i])
		}
	}
}

// --- rejectionSample ---

func TestRejectionSampleReturnsCorrectCount(t *testing.T) {
	result := rejectionSample(bi(1), bi(1_000_000), 20, seededRand())
	if len(result) != 20 {
		t.Errorf("got %d results, want 20", len(result))
	}
}

func TestRejectionSampleAllUnique(t *testing.T) {
	result := rejectionSample(bi(0), bi(5_000_000), 1000, seededRand())
	seen := make(map[string]struct{}, len(result))
	for _, v := range result {
		key := v.String()
		if _, exists := seen[key]; exists {
			t.Errorf("duplicate value: %s", key)
		}
		seen[key] = struct{}{}
	}
}

func TestRejectionSampleValuesInRange(t *testing.T) {
	lower := bi(100)
	rangeSize := bi(500) // values in [100, 600)
	upper := new(big.Int).Add(lower, rangeSize) // exclusive upper = 600
	result := rejectionSample(lower, rangeSize, 50, seededRand())
	for _, v := range result {
		if v.Cmp(lower) < 0 || v.Cmp(upper) >= 0 {
			t.Errorf("value %s out of range [100, 600)", v)
		}
	}
}

func TestRejectionSampleLargeBigInt(t *testing.T) {
	// Range spanning numbers larger than int64
	lower := bigFromString("100000000000000000000")
	rangeSize := bigFromString("999999999999999999999")
	result := rejectionSample(lower, rangeSize, 5, seededRand())
	if len(result) != 5 {
		t.Fatalf("got %d results, want 5", len(result))
	}
	upper := new(big.Int).Add(lower, rangeSize)
	for _, v := range result {
		if v.Cmp(lower) < 0 || v.Cmp(upper) >= 0 {
			t.Errorf("value %s out of range", v)
		}
	}
}

func TestRejectionSampleCountOne(t *testing.T) {
	result := rejectionSample(bi(42), bi(1_000_000), 1, seededRand())
	if len(result) != 1 {
		t.Fatalf("got %d results, want 1", len(result))
	}
}

func TestRejectionSampleAllUniqueStrings(t *testing.T) {
	// Regression: key collisions between distinct big.Ints must not occur
	lower := bigFromString("9999999999999999999999999999990")
	rangeSize := bi(100)
	result := rejectionSample(lower, rangeSize, 50, seededRand())
	seen := make(map[string]bool, len(result))
	for _, v := range result {
		k := v.String()
		if seen[k] {
			t.Errorf("duplicate: %s", k)
		}
		seen[k] = true
	}
}

// --- maxFisherYates threshold ---

func TestMaxFisherYatesConstant(t *testing.T) {
	if maxFisherYates <= 0 {
		t.Errorf("maxFisherYates = %d, must be positive", maxFisherYates)
	}
}

// Verify that a range of exactly maxFisherYates uses Fisher-Yates (smoke test:
// just ensure it returns the right count without panicking).
func TestFisherYatesAtThreshold(t *testing.T) {
	result := fisherYatesSample(bi(0), maxFisherYates, 10, seededRand())
	if len(result) != 10 {
		t.Errorf("got %d results at threshold, want 10", len(result))
	}
}
