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

// NewGoWorker 初始化一个调度器, 并指定Worker和Job的最大数量
func NewGoWorker(maxWorkers, maxQueue int) (wf *GoWorker) {
	wf = new(GoWorker)
	wf.dispatch = NewDispatcher(maxWorkers, maxQueue)
	// wf.dispatch.Debugger = false
	wf.dispatch.Run()
	return wf
}

// NewSimpleJob  初始化一个调度器, 并指定Worker和Job的最大数量
func NewSimpleJob(fuc func(*SimpleJob), id string, args interface{}) *SimpleJob {
	return &SimpleJob{
		ID:   id,
		fuc:  fuc,
		Args: args,
	}
}

// GoWorker GoWorker
type GoWorker struct {
	dispatch *Dispatcher
}

// AddJob 向流水线中放入工作job
func (wf *GoWorker) AddJob(wJob Job) {
	wf.dispatch.GJobQueue <- wJob
}

// WaitGoWorkerClose 等待工作流结果
func (wf *GoWorker) WaitGoWorkerClose() {
	<-wf.dispatch.Closed
}

// CloseGoWorker 关闭
func (wf *GoWorker) CloseGoWorker() {
	wf.dispatch.EndDispatch <- 0
}

// SimpleJob 简单工作负载
type SimpleJob struct {
	ID   string
	Args interface{}
	fuc  func(*SimpleJob)
}

// AddPayload 添加负载
func (job *SimpleJob) AddPayload(fuc func(*SimpleJob)) {
	job.fuc = fuc
}

// Play Play
func (job *SimpleJob) Play() {
	job.fuc(job)
}
