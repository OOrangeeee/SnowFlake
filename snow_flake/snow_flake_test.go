package snowflake

import (
	"sync"
	"testing"
)

func TestNewSnowFlakeCreatorForClusterWithDataCenter(t *testing.T) {
	t.Run("正常参数", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("不应该panic")
			}
		}()
		NewSnowFlakeCreatorForClusterWithDataCenter(1, 5, 1, 5)
	})

	t.Run("参数超过位数限制", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("应该触发panic")
			}
		}()
		NewSnowFlakeCreatorForClusterWithDataCenter(1, 20, 1, 20)
	})
}

func TestSingleInstanceIDGeneration(t *testing.T) {
	sf := NewSnowFlakeCreatorForSingle()
	prev := sf.GetId()
	for i := 0; i < 1000; i++ {
		newID := sf.GetId()
		if newID <= prev {
			t.Fatalf("ID不是单调递增: %d -> %d", prev, newID)
		}
		prev = newID
	}
}

func TestIDStructure(t *testing.T) {
	dcID := int64(3)
	workerID := int64(7)
	sf := NewSnowFlakeCreatorForClusterWithDataCenter(dcID, 5, workerID, 5)

	id := sf.GetId()
	// 验证各组成部分
	timePart := id >> (5 + 5 + 12)
	workerPart := (id >> (5 + 12)) & 0x1F
	dcPart := (id >> 12) & 0x1F
	seqPart := id & 0xFFF

	if workerPart != workerID || dcPart != dcID || seqPart != 0 {
		t.Errorf("ID结构错误: time=%d dc=%d worker=%d seq=%d", timePart, dcPart, workerPart, seqPart)
	}
}

func TestConcurrentGeneration(t *testing.T) {
	sf := NewSnowFlakeCreatorForSingle()
	const goroutines = 100
	const perRoutine = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	ids := make(chan int64, goroutines*perRoutine)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perRoutine; j++ {
				ids <- sf.GetId()
			}
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[int64]bool)
	for id := range ids {
		if seen[id] {
			t.Fatal("发现重复ID")
		}
		seen[id] = true
	}
}
