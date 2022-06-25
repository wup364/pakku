// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

//

package upool

import (
	"fmt"
	"strconv"

	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/strutil"
)

// Job 负载
type Job interface {
	Play()
}

// Worker 流水线
type Worker struct {
	Debugger        bool
	WorkerId        string
	Workbench       chan Job
	GWorkbenchQueue chan chan Job
	Finished        chan bool
}

// NewWorker 新建一条流水线
func NewWorker(WorkbenchQueue chan chan Job, Id string, debugger bool) *Worker {
	return &Worker{
		Debugger:        debugger,
		WorkerId:        Id,
		Workbench:       make(chan Job),
		GWorkbenchQueue: WorkbenchQueue,
		Finished:        make(chan bool),
	}
}

// Start 开始工作
func (w *Worker) Start() {
	go func() {
		for {
			// 将当前未分配待加工产品的工作台添加到工作台队列中
			w.GWorkbenchQueue <- w.Workbench
			if w.Debugger {
				logs.Infof("Add worker[%s] to workbench\r\n", w.WorkerId)
			}
			select {
			// 接收到了新的WorkerJob
			case wJob := <-w.Workbench:
				wJob.Play()
			case bFinished := <-w.Finished:
				if bFinished {
					if w.Debugger {
						logs.Infof("Worker [%s] closed\r\n", w.WorkerId)
					}
					return
				}
			}
		}
	}()
}

func (w *Worker) Stop() {
	// w.QuitChannel <- true
	go func() {
		w.Finished <- true
	}()
}

// Dispatcher 调度器
type Dispatcher struct {
	Debugger        bool
	DispatcherId    string        // 流水线ID
	MaxWorkers      int           // 流水线上的Worker最大数量
	Workers         []*Worker     // 流水线上所有Worker对象集合
	Closed          chan bool     // 流水线工作状态通道
	EndDispatch     chan int      // 流水线停止工作信号
	GJobQueue       chan Job      // 流水线上的所有待加工Job队列通道
	GWorkbenchQueue chan chan Job // 流水线上的所有操作台队列通道
}

// NewDispatcher 初始化调度者
func NewDispatcher(maxWorkers, maxQueue int) *Dispatcher {
	Closed := make(chan bool)
	EndDispatch := make(chan int)
	JobQueue := make(chan Job, maxQueue)
	WorkbenchQueue := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		DispatcherId:    strutil.GetUUID(),
		MaxWorkers:      maxWorkers,
		Closed:          Closed,
		EndDispatch:     EndDispatch,
		GJobQueue:       JobQueue,
		GWorkbenchQueue: WorkbenchQueue,
	}
}

// Run 开始运行
func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.GWorkbenchQueue, fmt.Sprintf("work-%s-%s", strconv.Itoa(i), strutil.GetUUID()), d.Debugger)
		if d.Debugger {
			logs.Infof("Created worker: %s\r\n", worker.WorkerId)
		}
		d.Workers = append(d.Workers, worker)
		worker.Start()
	}
	go d.Dispatch()
}

// Dispatch 监控
func (d *Dispatcher) Dispatch() {
FLAG:
	for {
		select {
		case endDispatch := <-d.EndDispatch:
			if d.Debugger {
				logs.Infof("Dispatcher close signal [%v]", endDispatch)
			}
			close(d.GJobQueue)
		case wJob, Ok := <-d.GJobQueue:
			if Ok {
				go func(wJob Job) {
					// 获取未使用的的工作台 & 将Job放入工作台
					Workbench := <-d.GWorkbenchQueue
					Workbench <- wJob
				}(wJob)
			} else {
				for _, w := range d.Workers {
					w.Stop()
				}
				d.Closed <- true
				break FLAG
			}
		}
	}
}
