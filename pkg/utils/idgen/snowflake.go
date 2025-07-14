package idgen

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

// Init 初始化 sonyflake，应该在程序启动时调用一次
func Init() error {
	var initErr error
	once.Do(func() {
		settings := sonyflake.Settings{
			StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		sf = sonyflake.NewSonyflake(settings)
		if sf == nil {
			initErr = fmt.Errorf("failed to initialize sonyflake")
		}
	})
	return initErr
}

// NewID 生成一个新的唯一 ID（int64）
func NewID() int64 {
	if sf == nil {
		panic("idgen not initialized: call idgen.Init() first")
	}
	id, err := sf.NextID()
	if err != nil {
		panic(fmt.Sprintf("failed to generate ID: %v", err))
	}
	return int64(id)
}
