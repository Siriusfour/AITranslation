package SnowFlak

import (
	"AITranslatio/Global"
	"go.uber.org/zap"
	"sync"
	"time"
)

// 创建一个雪花算法生成器(生成工厂)
func CreateSnowflakeFactory() *SnowFlake {
	return &SnowFlake{
		timestamp: 0,
		machineId: Global.Config.GetInt64("SnowFlake.SnowFlakeMachineId"),
		sequence:  0,
	}
}

type SnowFlake struct {
	sync.Mutex
	timestamp int64
	machineId int64
	sequence  int64
}

// 生成分布式ID
func (s *SnowFlake) GetId() int64 {

	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1e6

	//同毫秒内发号
	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & Global.SequenceMask
		if s.sequence == 0 {
			now = waitNextMillis(s.timestamp)
		}
	}

	//新毫秒发号
	if now > s.timestamp {
		s.sequence = 0
	}

	//时钟回拨
	if now < s.timestamp {
		threshold := Global.Config.GetInt64("SnowFlake.RollbackThresholdMs")
		if threshold <= 0 {
			threshold = 5
		}
		diff := s.timestamp - now
		if diff <= threshold {
			// 小回拨：沿用 lastTs 作为逻辑时间，并推进序列
			Global.Logger.Warn("snowflake clock rollback (minor), using logical time", zap.Int64("diff_ms", diff), zap.Int64("last_ts", s.timestamp))
			now = s.timestamp
			s.sequence = (s.sequence + 1) & Global.SequenceMask
			if s.sequence == 0 {
				now = waitNextMillis(s.timestamp)
			}
		} else {
			// 大回拨：阻塞到 lastTs
			Global.Logger.Error("snowflake clock rollback (major), blocking until last timestamp", zap.Int64("diff_ms", diff), zap.Int64("last_ts", s.timestamp))
			now = waitNextMillis(s.timestamp)
			s.sequence = 0
		}
	}

	s.timestamp = now

	r := (now-Global.StartTimeStamp)<<Global.TimestampShift | (s.machineId << Global.MachineIdShift) | (s.sequence)

	return r
}

// 等待到下一毫秒，直到当前毫秒时间 strictly 大于 lastTs
func waitNextMillis(lastTs int64) int64 {
	now := time.Now().UnixNano() / 1e6
	for now <= lastTs {
		time.Sleep(time.Millisecond)
		now = time.Now().UnixNano() / 1e6
	}
	return now
}
