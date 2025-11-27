package idgen

import (
	"sync"
	"testing"
	"time"
)

func TestNewShortIDGenerator(t *testing.T) {
	tests := []struct {
		name      string
		workerID  int64
		wantError bool
	}{
		{"valid worker ID 0", 0, false},
		{"valid worker ID 5", 5, false},
		{"valid worker ID 9", 9, false},
		{"invalid worker ID -1", -1, true},
		{"invalid worker ID 10", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewShortIDGenerator(tt.workerID)
			if (err != nil) != tt.wantError {
				t.Errorf("NewShortIDGenerator() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestShortIDGenerator_Generate(t *testing.T) {
	gen, err := NewShortIDGenerator(1)
	if err != nil {
		t.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	// 生成100个ID
	ids := make([]int64, 100)
	for i := 0; i < 100; i++ {
		id, err := gen.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		ids[i] = id

		// 检查ID长度（应该是13位）
		if id < 1000000000000 || id > 9999999999999 {
			t.Errorf("ID %d is not 13 digits", id)
		}
	}

	// 检查唯一性
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("Duplicate ID found: %d", id)
		}
		seen[id] = true
	}

	// 检查递增性
	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("IDs not increasing: %d <= %d", ids[i], ids[i-1])
		}
	}
}

func TestShortIDGenerator_ConcurrentGeneration(t *testing.T) {
	gen, err := NewShortIDGenerator(1)
	if err != nil {
		t.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	const goroutines = 10
	const idsPerGoroutine = 50
	const totalIDs = goroutines * idsPerGoroutine

	ids := make(chan int64, totalIDs)
	var wg sync.WaitGroup

	// 并发生成ID
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := gen.Generate()
				if err != nil {
					t.Errorf("Generate failed: %v", err)
					return
				}
				ids <- id
			}
		}()
	}

	wg.Wait()
	close(ids)

	// 检查唯一性
	seen := make(map[int64]bool)
	count := 0
	for id := range ids {
		if seen[id] {
			t.Fatalf("Duplicate ID in concurrent test: %d", id)
		}
		seen[id] = true
		count++
	}

	if count != totalIDs {
		t.Errorf("Expected %d IDs, got %d", totalIDs, count)
	}
}

func TestParseShortID(t *testing.T) {
	gen, err := NewShortIDGenerator(7)
	if err != nil {
		t.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	id, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	timestamp, workerID, sequence := ParseShortID(id)

	// 检查worker ID
	if workerID != 7 {
		t.Errorf("Expected worker ID 7, got %d", workerID)
	}

	// 检查时间戳（应该接近当前时间）
	now := time.Now().Unix()
	if timestamp < now-2 || timestamp > now+2 {
		t.Errorf("Timestamp out of range: %d, expected around %d", timestamp, now)
	}

	// 序列号应该在有效范围内
	if sequence < 0 || sequence > maxShortSequence {
		t.Errorf("Sequence out of range: %d", sequence)
	}

	t.Logf("ID: %d = timestamp:%d + workerID:%d + sequence:%02d",
		id, timestamp, workerID, sequence)
}

func TestGetShortIDTimestamp(t *testing.T) {
	gen, err := NewShortIDGenerator(1)
	if err != nil {
		t.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	beforeGen := time.Now()
	id, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	afterGen := time.Now()

	extractedTime := GetShortIDTimestamp(id)

	// 提取的时间应该在生成前后之间（允许2秒误差）
	if extractedTime.Before(beforeGen.Add(-2*time.Second)) ||
		extractedTime.After(afterGen.Add(2*time.Second)) {
		t.Errorf("Extracted time %v not between %v and %v",
			extractedTime, beforeGen, afterGen)
	}
}

func TestInitShort(t *testing.T) {
	err := InitShort(5)
	if err != nil {
		t.Fatalf("InitShort failed: %v", err)
	}

	// 测试GenerateShortID
	id, err := GenerateShortID()
	if err != nil {
		t.Fatalf("GenerateShortID failed: %v", err)
	}

	if id <= 0 || id > 9999999999999 {
		t.Errorf("Generated ID should be 13 digits, got %d", id)
	}

	t.Logf("Generated short ID: %d", id)
}

func TestDifferentWorkerIDs_ShortID(t *testing.T) {
	gen1, err := NewShortIDGenerator(1)
	if err != nil {
		t.Fatalf("NewShortIDGenerator(1) failed: %v", err)
	}

	gen2, err := NewShortIDGenerator(2)
	if err != nil {
		t.Fatalf("NewShortIDGenerator(2) failed: %v", err)
	}

	// 同时生成ID
	id1, err := gen1.Generate()
	if err != nil {
		t.Fatalf("Generate from gen1 failed: %v", err)
	}

	id2, err := gen2.Generate()
	if err != nil {
		t.Fatalf("Generate from gen2 failed: %v", err)
	}

	// ID应该不同
	if id1 == id2 {
		t.Errorf("IDs from different workers should be different: %d == %d", id1, id2)
	}

	// 验证worker ID
	_, worker1, _ := ParseShortID(id1)
	_, worker2, _ := ParseShortID(id2)

	if worker1 != 1 {
		t.Errorf("Expected worker ID 1, got %d", worker1)
	}
	if worker2 != 2 {
		t.Errorf("Expected worker ID 2, got %d", worker2)
	}

	t.Logf("ID1: %d (worker %d), ID2: %d (worker %d)",
		id1, worker1, id2, worker2)
}

func TestShortID_Format(t *testing.T) {
	gen, err := NewShortIDGenerator(9)
	if err != nil {
		t.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	for i := 0; i < 10; i++ {
		id, err := gen.Generate()
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		// 检查ID格式
		ts, wid, seq := ParseShortID(id)
		reconstructed := ts*1000 + wid*100 + seq

		if reconstructed != id {
			t.Errorf("ID reconstruction failed: %d != %d", reconstructed, id)
		}

		t.Logf("ID: %13d = %10d (ts) + %d (worker) + %02d (seq)",
			id, ts, wid, seq)
	}
}

func BenchmarkShortIDGenerator_Generate(b *testing.B) {
	gen, err := NewShortIDGenerator(1)
	if err != nil {
		b.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Generate()
		if err != nil {
			b.Fatalf("Generate failed: %v", err)
		}
	}
}

func BenchmarkShortIDGenerator_GenerateParallel(b *testing.B) {
	gen, err := NewShortIDGenerator(1)
	if err != nil {
		b.Fatalf("NewShortIDGenerator failed: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := gen.Generate()
			if err != nil {
				b.Fatalf("Generate failed: %v", err)
			}
		}
	})
}
