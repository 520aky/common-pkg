package snowflake

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// 经典位宽
const (
	workerBits   = 10 // 节点ID 0..1023
	sequenceBits = 12 // 每毫秒序列 0..4095
	maxWorkerID  = (1 << workerBits) - 1
	seqMask      = (1 << sequenceBits) - 1

	workerShift = sequenceBits
	timeShift   = workerBits + sequenceBits
)

// 自定义纪元(毫秒)：2024-01-01T00:00:00Z
const epochMs = int64(1704067200000)

type Generator struct {
	mu       sync.Mutex
	workerID int64
	seq      int64
	lastMs   int64
}

// New 创建生成器，workerID 取值 0..1023
func New(workerID int64) (*Generator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("workerID out of range: %d (0..%d)", workerID, maxWorkerID)
	}
	return &Generator{workerID: workerID}, nil
}

// NewFromEnv 支持从环境变量读取：SNOWFLAKE_WORKER_ID
func NewFromEnv() *Generator {
	if s := os.Getenv("SNOWFLAKE_WORKER_ID"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil && n >= 0 && n <= maxWorkerID {
			g, _ := New(n)
			return g
		}
	}
	// 简单兜底：基于主机名做一个稳定hash到 0..1023（避免碰撞请使用注册中心/Redis分配）
	hn, _ := os.Hostname()
	var h uint32 = 2166136261
	for i := 0; i < len(hn); i++ {
		h ^= uint32(hn[i])
		h *= 16777619
	}
	g, _ := New(int64(h % uint32(maxWorkerID+1)))
	return g
}

// NextID 返回 19 位十进制可排序的 int64 ID（全局唯一）
func (g *Generator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := nowMs()
	if now < g.lastMs {
		// 时钟回拨：等待至上次时间戳（保守策略）
		now = waitUntil(g.lastMs)
	}
	if now == g.lastMs {
		g.seq = (g.seq + 1) & seqMask
		if g.seq == 0 {
			// 本毫秒序列耗尽，等待下一毫秒
			now = waitUntil(g.lastMs + 1)
		}
	} else {
		g.seq = 0
	}
	g.lastMs = now

	ts := now - epochMs
	id := (ts << timeShift) | (g.workerID << workerShift) | g.seq
	return id
}

// NextString 同 NextID，但返回字符串
func (g *Generator) NextString() string { return strconv.FormatInt(g.NextID(), 10) }

// Parse 解析ID，返回(时间戳ms, workerID, sequence)
func Parse(id int64) (tsMs int64, workerID int64, seq int64, err error) {
	if id < 0 {
		return 0, 0, 0, errors.New("invalid id")
	}
	ts := (id >> timeShift) + epochMs
	worker := (id >> workerShift) & int64(maxWorkerID)
	sequence := id & int64(seqMask)
	return ts, worker, sequence, nil
}

func nowMs() int64 { return time.Now().UnixNano() / 1e6 }

func waitUntil(targetMs int64) int64 {
	for {
		n := nowMs()
		if n >= targetMs {
			return n
		}
		time.Sleep(time.Millisecond)
	}
}
