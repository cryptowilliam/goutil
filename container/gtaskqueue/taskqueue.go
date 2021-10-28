package gtaskqueue

// Task sender queue in memory with statistic function.

import (
	"github.com/cryptowilliam/goutil/sys/gtime"
	"sync"
	"time"
)

type Priority int

const (
	PriorityHighest Priority = iota + 0
	PriorityHigh
	PriorityMedium
	PriorityLow
	PriorityLowest
)

const (
	PriorityCount int = 5
)

var (
	ErrorExecTimeout error = gerrors.Errorf("Exec timeout")
)

type Task struct {
	Id            string
	PriorityValue Priority
	ExecTimeout   time.Duration
	DataPtr       interface{}

	// Statistics.
	CreateTime    time.Time
	ExecBeginTime time.Time
	ExecDoneTime  time.Time
	LastError     error
}

type Statistic struct {
	// 此刻排队任务个数
	Now2doSize int64
	// 此刻正在执行的任务个数
	NowExecSize int64

	// 近期的开始时间
	LatelyStatisticBeginTime time.Time
	// 近期任务成功次数
	LatelySuccessSize int64
	// 近期等待执行超时次数
	//LatelyWaitTimeoutSize int64
	// 近期执行超时次数
	LatelyExecTimeoutSize int64
	// 近期执行出错次数
	LatelyExecErrorSize int64
	// 近期执行平均耗时
	LatelyAvgExecDuration time.Duration
}

type TaskQueue struct {
	to2Do    [PriorityCount]chan *Task
	exec     map[string]*Task
	execLock sync.RWMutex
	done     chan *Task

	tempDataLock sync.RWMutex
	// 近期的开始时间
	tempLatelyStatisticBeginTime time.Time
	// 近期任务成功次数
	tempLatelySuccessSize int64
	// 近期执行超时次数
	tempLatelyExecTimeoutSize int64
	// 近期执行出错次数
	tempLatelyExecErrorSize int64
	// 近期执行的任务的耗时总和，包括执行超时的在内
	tempLatelyExecDurationSum time.Duration
	// 近期执行的任务的次数，包括执行超时的在内
	tempLatelyExecCount int64
}

func New(autoResetStatistic time.Duration, size int) *TaskQueue {
	q := TaskQueue{}
	q.exec = make(map[string]*Task)
	for i := 0; i < PriorityCount; i++ {
		q.to2Do[Priority(i)] = make(chan *Task, size)
	}

	// Auto reset lately statistic data.
	go func() {
		time.Sleep(autoResetStatistic)
		q.tempDataLock.Lock()
		q.tempLatelyStatisticBeginTime = time.Now()
		q.tempLatelySuccessSize = 0
		q.tempLatelyExecTimeoutSize = 0
		q.tempLatelyExecErrorSize = 0
		q.tempDataLock.Unlock()
	}()

	// Every 2 seconds, auto check whether task timeout when executing.
	go func() {
		time.Sleep(time.Second * 2)

		q.execLock.RLock()
		execClone := q.exec
		q.execLock.RUnlock()
		for k := range execClone {
			item := execClone[k]
			if time.Now().Sub(item.ExecBeginTime) > item.ExecTimeout {
				q.PushDone(item.Id, ErrorExecTimeout)
			}
		}
	}()

	return &q
}

func (q *TaskQueue) GetStatistic() *Statistic {
	res := Statistic{}

	q.tempDataLock.RLock()
	res.LatelyStatisticBeginTime = q.tempLatelyStatisticBeginTime
	res.LatelySuccessSize = q.tempLatelySuccessSize
	res.LatelyExecTimeoutSize = q.tempLatelyExecTimeoutSize
	res.LatelyExecErrorSize = q.tempLatelyExecErrorSize
	res.LatelyAvgExecDuration = gtime.NsecToDuration(q.tempLatelyExecDurationSum.Nanoseconds() / q.tempLatelyExecCount)
	q.tempDataLock.RUnlock()

	for i := 0; i < PriorityCount; i++ {
		res.Now2doSize += int64(len(q.to2Do[i]))
	}
	q.execLock.RLock()
	res.NowExecSize += int64(len(q.exec))
	q.execLock.RUnlock()

	return &res
}

// Wait until pushed.
func (q *TaskQueue) Push2doWait(priority Priority /*waitTimeout, */, execTimeout time.Duration, dataPtr interface{}) {
	item := Task{}
	item.PriorityValue = priority
	item.ExecTimeout = execTimeout
	item.DataPtr = dataPtr
	item.CreateTime = time.Now()
	q.to2Do[priority] <- &item
}

// Peek a task to exec.
func (q *TaskQueue) Pop2do() (*Task, error) {
	var item *Task

	for i := 0; i < PriorityCount; i++ {
		for {
			// Pop one item.
			if len(q.to2Do[i]) == 0 { // Empty channel.
				break
			}
			select {
			case <-time.After(time.Millisecond * 50): // Empty channel.
				continue
			case item = <-q.to2Do[i]:
			}

			// Push into exec map.
			q.execLock.Lock()
			q.exec[item.Id] = item
			q.execLock.Unlock()
			return item, nil
		}
	}
	return nil, gerrors.Errorf("No to2Do item")
}

// Tell taskQueue some task is done.
func (q *TaskQueue) PushDone(id string, err error) {
	// Remove from exec map.
	q.execLock.Lock()
	item := q.exec[id]
	delete(q.exec, id)
	q.execLock.Unlock()
	if item == nil {
		return
	}

	// Write exec done time.
	item.ExecDoneTime = time.Now()
	item.LastError = err

	// Write statistic data.
	q.tempDataLock.Lock()
	q.tempLatelyExecDurationSum += item.ExecDoneTime.Sub(item.ExecBeginTime)
	q.tempLatelyExecCount++
	if err == nil {
		q.tempLatelySuccessSize++
	} else if err == ErrorExecTimeout {
		q.tempLatelyExecTimeoutSize++
	} else {
		q.tempLatelyExecErrorSize++
	}
	q.tempDataLock.Unlock()

	// Push into done channel.
	q.done <- item
}

// Pop Task from done channel.
func (q *TaskQueue) PopDone() (*Task, error) {
	if len(q.done) == 0 {
		return nil, gerrors.Errorf("No done item")
	}
	return <-q.done, nil
}

// Reset taskQueue.
func (q *TaskQueue) Close() {
	for i := 0; i < PriorityCount; i++ {
		close(q.to2Do[i])
	}

	q.execLock.Lock()
	for k := range q.exec {
		delete(q.exec, k)
	}
	q.execLock.Unlock()

	close(q.done)

	q.tempDataLock.Lock()
	q.tempLatelyStatisticBeginTime = time.Time{}
	q.tempLatelySuccessSize = 0
	q.tempLatelyExecTimeoutSize = 0
	q.tempLatelyExecErrorSize = 0
	q.tempLatelyExecDurationSum = time.Duration(0)
	q.tempLatelyExecCount = 0
	q.tempDataLock.Unlock()
}
