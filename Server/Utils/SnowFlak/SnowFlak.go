package SnowFlak

import (
	"AITranslatio/Global/Consts"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type SnowFlakManager interface {
	GetID() int64
	GetIDString() string
}

type SnowFlake struct {
	sync.Mutex
	timestamp           int64 //上次生成ID是时间戳(毫秒)
	machineId           int64
	sequence            int64
	getTime             func() int64
	RollbackThresholdMs int64
	logger              map[string]*zap.Logger
}

// 创建一个雪花算法生成器(生成工厂)
func CreateSnowflakeFactory(RollbackThresholdMs int64, logger map[string]*zap.Logger) *SnowFlake {
	return &SnowFlake{
		timestamp:           0,
		machineId:           1,
		sequence:            0,
		getTime:             func() int64 { return time.Now().UnixNano() / 1e6 }, //依赖注入，默认获取当前时间，测试时注入测试时间
		RollbackThresholdMs: RollbackThresholdMs,
		logger:              logger,
	}
}

// 生成ID
func (s *SnowFlake) GetID() int64 {

	s.Lock()
	defer s.Unlock()

	now := s.getTime()

	//同毫秒内发号
	if now == s.timestamp {
		// 同毫秒内序列号自增1，&SequenceMask保证不会溢出，在0-4095内循环
		s.sequence = (s.sequence + 1) & Consts.SequenceMask
		//如果溢出了，则阻塞到1下一个1毫秒
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
		threshold := s.RollbackThresholdMs
		if threshold <= 0 {
			threshold = 5
		}
		diff := s.timestamp - now
		if diff <= threshold {
			// 小回拨：沿用 lastTs 作为逻辑时间，并推进序列
			s.logger["Business"].Warn("snowflake 时钟回退，使用逻辑时间", zap.Int64("diff_ms", diff), zap.Int64("last_ts", s.timestamp))
			now = s.timestamp
			s.sequence = (s.sequence + 1) & Consts.SequenceMask
			if s.sequence == 0 {
				now = waitNextMillis(s.timestamp)
			}
		} else {
			// 大回拨：阻塞到 lastTs
			s.logger["Business"].Error("snowflake clock rollback (major), blocking until last timestamp", zap.Int64("diff_ms", diff), zap.Int64("last_ts", s.timestamp))
			now = waitNextMillis(s.timestamp)
			s.sequence = 0
		}
	}

	s.timestamp = now

	r := (now-Consts.StartTimeStamp)<<Consts.TimestampShift | (s.machineId << Consts.MachineIdShift) | (s.sequence)

	return r
}

func (s *SnowFlake) GetIDString() string {
	return strconv.FormatInt(s.GetID(), 10)
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
