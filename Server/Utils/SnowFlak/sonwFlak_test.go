package SnowFlak

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

type mockConfig struct{}

func (logger mockConfig) GetInt64(keyName string) int64 { return 0 }

func initSnowflake() *SnowFlakeGenerator {
	c := CreateSnowflakeFactory()
	c.logger = map[string]*zap.Logger{
		"Business": zap.NewNop(),
	}
	c.config = mockConfig{}
	return c
}

func TestGetID(t *testing.T) {

	//成功情况，同毫秒内生成的1000个ID不重合且递增
	success := initSnowflake()

	var lastID int64
	for count := 0; count < 100; count++ {
		ID := success.GetID()
		if ID <= lastID {
			t.Errorf("验证失败，ID重复")
		}
		fmt.Println(ID)
		lastID = ID
	}

	//测试同毫秒发号
	a := initSnowflake()

	a.getTime = func() int64 {
		return a.timestamp
	}
	IDA := a.GetID()
	IDB := a.GetID()
	if IDA == IDB {
		t.Errorf("同毫秒发情况下两个ID相同：A:%d，B:%d", IDA, IDB)
	}

	//测试同毫秒发号，但序列号溢出情况
	b := initSnowflake()
	b.getTime = func() int64 {
		return b.timestamp
	}
	b.sequence = -1
	IDC := b.GetID()
	IDD := b.GetID()

	if IDC == IDD {
		t.Errorf("同毫秒发号,且超过了4095的情况下两个ID相同：C:%d，D:%d", IDC, IDD)
	}

	//测试小回拨，判断新的ID相对于当前时间生成的ID是不是递增的
	c := initSnowflake()

	c.timestamp = 10
	c.getTime = func() int64 {
		return c.timestamp
	}
	IDE := c.GetID()

	c.getTime = func() int64 {
		return c.timestamp - 10
	}
	IDF := c.GetID()

	if IDE > IDF {
		t.Errorf("小回拨情况下ID倒退")
	}

	//测试大回拨，判断新的ID相对于当前时间生成的ID是不是递增的
	d := initSnowflake()

	d.config = mockConfig{}
	d.timestamp = 50
	d.getTime = func() int64 {
		return d.timestamp
	}
	IDG := d.GetID()

	d.getTime = func() int64 {
		return d.timestamp - 50
	}
	IDH := d.GetID()

	if IDG > IDH {
		t.Errorf("大回拨情况下ID倒退")
	}
}
