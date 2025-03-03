package snowflake

import (
	"errors"
	"math"
	"strconv"
	"sync"
	"time"
)

type SnowFlakeCreator struct {
	mu                sync.Mutex
	lastTimestamp     int64
	datacenterId      int64
	datacenterIdShift int64
	workerId          int64
	workerIdShift     int64
	sequence          int64
	sequenceShift     int64
	maxSequence       int64
}

const (
	defaultSequenceShift     = 12
	defaultWorkerIdShift     = 5
	defaultDataCenterIdShift = 5
	allShift                 = 22
	maxTime                  = 1<<41 - 1
)

var (
	startTime int64
)

func init() {
	startTimeTmp := time.Date(2025, 3, 3, 0, 0, 0, 0, time.UTC).UnixMilli()
	startTime = startTimeTmp
}

// NewSnowFlakeCreatorForClusterWithDataCenter
// 生成一个分布式、区分数据中心和worker的雪花算法生成器。
// 如果数据中心id和workerid的长度总和大小超过allShift则触发panic。
// 如果数值大于其位数最大值则panic。
func NewSnowFlakeCreatorForClusterWithDataCenter(datacenterId int64, datacenterIdShift int64, workerId int64, workerIdShift int64) *SnowFlakeCreator {
	if datacenterIdShift+workerIdShift >= allShift {
		panic(errors.New("参数错误，数据中心ID+工作机器ID长度超过" + strconv.Itoa(allShift)))
	} else if workerId >= int64(math.Pow(2.0, float64(workerIdShift))) || datacenterId >= int64(math.Pow(2.0, float64(datacenterIdShift))) {
		panic(errors.New("参数错误，ID值大于其位数最大值"))
	}
	return &SnowFlakeCreator{
		mu:                sync.Mutex{},
		datacenterId:      datacenterId,
		datacenterIdShift: datacenterIdShift,
		workerId:          workerId,
		workerIdShift:     workerIdShift,
		sequenceShift:     allShift - datacenterIdShift - workerIdShift,
		maxSequence:       int64(math.Pow(2.0, float64(allShift-datacenterIdShift-workerIdShift))) - 1,
	}
}

// NewSnowFlakeCreatorForClusterWithoutDataCenter
// 生成一个分布式、不区分数据中心和worker的雪花算法生成器。
// 如果workerid长度不为workerIdDataCenterIdShift则panic。
// 如果数值大于其位数最大值则panic。
func NewSnowFlakeCreatorForClusterWithoutDataCenter(workerId int64, workerIdShift int64) *SnowFlakeCreator {
	if workerIdShift >= allShift {
		panic(errors.New("参数错误，机器ID不等于" + strconv.Itoa(defaultSequenceShift)))
	} else if workerId >= int64(math.Pow(2.0, float64(workerIdShift))) {
		panic(errors.New("参数错误，机器ID大于其位数最大值"))
	}
	return &SnowFlakeCreator{
		mu:                sync.Mutex{},
		datacenterId:      0,
		datacenterIdShift: 0,
		workerId:          workerId,
		workerIdShift:     workerIdShift,
		sequenceShift:     allShift - workerIdShift,
		maxSequence:       int64(math.Pow(2.0, float64(allShift-workerIdShift))) - 1,
	}
}

// NewSnowFlakeCreatorForSingle
// 生成一个单机下的雪花算法生成器
func NewSnowFlakeCreatorForSingle() *SnowFlakeCreator {
	return &SnowFlakeCreator{
		mu:                sync.Mutex{},
		datacenterId:      0,
		datacenterIdShift: 0,
		workerId:          0,
		workerIdShift:     0,
		sequenceShift:     allShift,
		maxSequence:       int64(math.Pow(2.0, float64(allShift))) - 1,
	}
}

func (sf *SnowFlakeCreator) getNowTime() {
	now := time.Now().UnixMilli()
	now = now - startTime
	if now > maxTime {
		panic("时间错大于41位")
	}
	// 处理时钟回调
	if now <= sf.lastTimestamp {
		return
	}
	sf.lastTimestamp = now
	sf.sequence = 0
}

func (sf *SnowFlakeCreator) createId() int64 {
	return sf.lastTimestamp<<(sf.workerIdShift+sf.datacenterIdShift+sf.sequenceShift) |
		sf.workerId<<(sf.datacenterIdShift+sf.sequenceShift) |
		sf.datacenterId<<(sf.sequenceShift) |
		sf.sequence
}

// GetId
// 获取一个雪花ID
func (sf *SnowFlakeCreator) GetId() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	sf.getNowTime()
	// 这一毫秒的id数量用完了
	if sf.sequence > sf.maxSequence {
		for sf.sequence != 0 {
			sf.getNowTime()
		}
	}
	ans := sf.createId()
	sf.sequence++
	return ans
}
