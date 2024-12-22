package strutil

import (
	"sync"
	"testing"
	"time"
)

func TestUUIDBasic(t *testing.T) {
	uuid := GetUUID()
	if len(uuid) != 32 {
		t.Errorf("UUID length should be 32, got %d", len(uuid))
	}
}

func TestUUIDUniqueness(t *testing.T) {
	count := 10000
	uuids := make(map[string]bool, count)
	for i := 0; i < count; i++ {
		uuid := GetUUID()
		if uuids[uuid] {
			t.Errorf("Duplicate UUID found: %s", uuid)
		}
		uuids[uuid] = true
	}
}

func TestUUIDConcurrent(t *testing.T) {
	count := 100
	concurrent := 10
	uuids := sync.Map{}
	wg := sync.WaitGroup{}

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				uuid := GetUUID()
				if _, loaded := uuids.LoadOrStore(uuid, true); loaded {
					t.Errorf("Duplicate UUID found in concurrent test: %s", uuid)
				}
			}
		}()
	}
	wg.Wait()
}

func TestUUIDClockBackwards(t *testing.T) {
	origLastNano := lastNano
	defer func() { lastNano = origLastNano }()

	lastNano = uint64(time.Now().Add(time.Second).UnixNano())
	uuids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		uuid := GetUUID()
		if uuids[uuid] {
			t.Errorf("Duplicate UUID found during clock backwards: %s", uuid)
		}
		uuids[uuid] = true
	}
}

func TestUUIDPerformance(t *testing.T) {
	start := time.Now()
	count := 100000
	for i := 0; i < count; i++ {
		GetUUID()
	}
	duration := time.Since(start)

	avgDuration := duration.Nanoseconds() / int64(count)
	if avgDuration > 1000 {
		t.Errorf("UUID generation too slow: %d ns/op", avgDuration)
	}
}

// 基准测试
func BenchmarkGetUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUUID()
	}
}

// 并发基准测试
func BenchmarkGetUUIDConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GetUUID()
		}
	})
}

// 示例
func ExampleGetUUID() {
	uuid := GetUUID()
	// 输出UUID长度
	println(len(uuid)) // 32
}
