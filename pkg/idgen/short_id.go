package idgen

import (
	"errors"
	"sync"
	"time"
)

// ShortID生成器 - 生成13位数字ID
// ID结构：时间戳(秒,10位) + 机器ID(1位) + 序列号(2位)
// 示例：1732701603 + 1 + 23 = 1732701603123
//
// 特点：
// - 13位数字，易读易记
// - 支持10台机器（0-9）
// - 每秒每台机器可生成100个ID（00-99）
// - 趋势递增，利于数据库索引

const (
	// 机器ID位数（支持0-9，共10台机器）
	shortWorkerIDBits = 1
	// 序列号位数（支持00-99，每秒100个ID）
	shortSequenceBits = 2

	// 最大机器ID（0-9）
	maxShortWorkerID = 9
	// 最大序列号（0-99）
	maxShortSequence = 99
)

// ShortIDGenerator 短ID生成器
type ShortIDGenerator struct {
	mu            sync.Mutex
	lastTimestamp int64  // 上次生成ID的时间戳（秒）
	workerID      int64  // 机器ID（0-9）
	sequence      int64  // 序列号（0-99）
}

// NewShortIDGenerator 创建短ID生成器
func NewShortIDGenerator(workerID int64) (*ShortIDGenerator, error) {
	if workerID < 0 || workerID > maxShortWorkerID {
		return nil, errors.New("worker ID must be between 0 and 9")
	}

	return &ShortIDGenerator{
		workerID: workerID,
		sequence: 0,
	}, nil
}

// Generate 生成下一个短ID
func (g *ShortIDGenerator) Generate() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().Unix() // 秒级时间戳

	// 时钟回退检测
	if timestamp < g.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	// 同一秒内
	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) % (maxShortSequence + 1)
		// 序列号溢出，等待下一秒
		if g.sequence == 0 {
			timestamp = g.waitNextSecond(g.lastTimestamp)
		}
	} else {
		// 新的秒，序列号重置
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	// 组装13位ID：时间戳(10位) + 机器ID(1位) + 序列号(2位)
	id := timestamp*1000 + g.workerID*100 + g.sequence

	return id, nil
}

// waitNextSecond 等待下一秒
func (g *ShortIDGenerator) waitNextSecond(lastTimestamp int64) int64 {
	timestamp := time.Now().Unix()
	for timestamp <= lastTimestamp {
		time.Sleep(time.Millisecond * 10) // 休眠10毫秒后重试
		timestamp = time.Now().Unix()
	}
	return timestamp
}

// ParseShortID 解析短ID
func ParseShortID(id int64) (timestamp int64, workerID int64, sequence int64) {
	timestamp = id / 1000                    // 前10位
	workerID = (id % 1000) / 100             // 第11位
	sequence = id % 100                       // 后2位
	return
}

// GetShortIDTimestamp 从短ID中提取时间戳
func GetShortIDTimestamp(id int64) time.Time {
	timestamp, _, _ := ParseShortID(id)
	return time.Unix(timestamp, 0)
}

// 全局短ID生成器实例
var globalShortGenerator *ShortIDGenerator
var shortOnce sync.Once

// InitShort 初始化全局短ID生成器
func InitShort(workerID int64) error {
	var err error
	shortOnce.Do(func() {
		globalShortGenerator, err = NewShortIDGenerator(workerID)
	})
	return err
}

// GenerateShortID 生成全局唯一短ID
func GenerateShortID() (int64, error) {
	if globalShortGenerator == nil {
		return 0, errors.New("short ID generator not initialized, call InitShort() first")
	}
	return globalShortGenerator.Generate()
}

// Init 初始化ID生成器（兼容接口）
func Init(workerID int64) error {
	return InitShort(workerID)
}

// GenerateID 生成ID（兼容接口）
func GenerateID() (int64, error) {
	return GenerateShortID()
}
