package pow

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ClientWorkerMock implements ClientWorker
type ClientWorkerMock struct {
	t minimock.Tester

	funcDoWork          func(ctx context.Context, data Data) (d1 Data, err error)
	inspectFuncDoWork   func(ctx context.Context, data Data)
	afterDoWorkCounter  uint64
	beforeDoWorkCounter uint64
	DoWorkMock          mClientWorkerMockDoWork
}

// NewClientWorkerMock returns a mock for ClientWorker
func NewClientWorkerMock(t minimock.Tester) *ClientWorkerMock {
	m := &ClientWorkerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DoWorkMock = mClientWorkerMockDoWork{mock: m}
	m.DoWorkMock.callArgs = []*ClientWorkerMockDoWorkParams{}

	return m
}

type mClientWorkerMockDoWork struct {
	mock               *ClientWorkerMock
	defaultExpectation *ClientWorkerMockDoWorkExpectation
	expectations       []*ClientWorkerMockDoWorkExpectation

	callArgs []*ClientWorkerMockDoWorkParams
	mutex    sync.RWMutex
}

// ClientWorkerMockDoWorkExpectation specifies expectation struct of the ClientWorker.DoWork
type ClientWorkerMockDoWorkExpectation struct {
	mock    *ClientWorkerMock
	params  *ClientWorkerMockDoWorkParams
	results *ClientWorkerMockDoWorkResults
	Counter uint64
}

// ClientWorkerMockDoWorkParams contains parameters of the ClientWorker.DoWork
type ClientWorkerMockDoWorkParams struct {
	ctx  context.Context
	data Data
}

// ClientWorkerMockDoWorkResults contains results of the ClientWorker.DoWork
type ClientWorkerMockDoWorkResults struct {
	d1  Data
	err error
}

// Expect sets up expected params for ClientWorker.DoWork
func (mmDoWork *mClientWorkerMockDoWork) Expect(ctx context.Context, data Data) *mClientWorkerMockDoWork {
	if mmDoWork.mock.funcDoWork != nil {
		mmDoWork.mock.t.Fatalf("ClientWorkerMock.DoWork mock is already set by Set")
	}

	if mmDoWork.defaultExpectation == nil {
		mmDoWork.defaultExpectation = &ClientWorkerMockDoWorkExpectation{}
	}

	mmDoWork.defaultExpectation.params = &ClientWorkerMockDoWorkParams{ctx, data}
	for _, e := range mmDoWork.expectations {
		if minimock.Equal(e.params, mmDoWork.defaultExpectation.params) {
			mmDoWork.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDoWork.defaultExpectation.params)
		}
	}

	return mmDoWork
}

// Inspect accepts an inspector function that has same arguments as the ClientWorker.DoWork
func (mmDoWork *mClientWorkerMockDoWork) Inspect(f func(ctx context.Context, data Data)) *mClientWorkerMockDoWork {
	if mmDoWork.mock.inspectFuncDoWork != nil {
		mmDoWork.mock.t.Fatalf("Inspect function is already set for ClientWorkerMock.DoWork")
	}

	mmDoWork.mock.inspectFuncDoWork = f

	return mmDoWork
}

// Return sets up results that will be returned by ClientWorker.DoWork
func (mmDoWork *mClientWorkerMockDoWork) Return(d1 Data, err error) *ClientWorkerMock {
	if mmDoWork.mock.funcDoWork != nil {
		mmDoWork.mock.t.Fatalf("ClientWorkerMock.DoWork mock is already set by Set")
	}

	if mmDoWork.defaultExpectation == nil {
		mmDoWork.defaultExpectation = &ClientWorkerMockDoWorkExpectation{mock: mmDoWork.mock}
	}
	mmDoWork.defaultExpectation.results = &ClientWorkerMockDoWorkResults{d1, err}
	return mmDoWork.mock
}

// Set uses given function f to mock the ClientWorker.DoWork method
func (mmDoWork *mClientWorkerMockDoWork) Set(f func(ctx context.Context, data Data) (d1 Data, err error)) *ClientWorkerMock {
	if mmDoWork.defaultExpectation != nil {
		mmDoWork.mock.t.Fatalf("Default expectation is already set for the ClientWorker.DoWork method")
	}

	if len(mmDoWork.expectations) > 0 {
		mmDoWork.mock.t.Fatalf("Some expectations are already set for the ClientWorker.DoWork method")
	}

	mmDoWork.mock.funcDoWork = f
	return mmDoWork.mock
}

// When sets expectation for the ClientWorker.DoWork which will trigger the result defined by the following
// Then helper
func (mmDoWork *mClientWorkerMockDoWork) When(ctx context.Context, data Data) *ClientWorkerMockDoWorkExpectation {
	if mmDoWork.mock.funcDoWork != nil {
		mmDoWork.mock.t.Fatalf("ClientWorkerMock.DoWork mock is already set by Set")
	}

	expectation := &ClientWorkerMockDoWorkExpectation{
		mock:   mmDoWork.mock,
		params: &ClientWorkerMockDoWorkParams{ctx, data},
	}
	mmDoWork.expectations = append(mmDoWork.expectations, expectation)
	return expectation
}

// Then sets up ClientWorker.DoWork return parameters for the expectation previously defined by the When method
func (e *ClientWorkerMockDoWorkExpectation) Then(d1 Data, err error) *ClientWorkerMock {
	e.results = &ClientWorkerMockDoWorkResults{d1, err}
	return e.mock
}

// DoWork implements ClientWorker
func (mmDoWork *ClientWorkerMock) DoWork(ctx context.Context, data Data) (d1 Data, err error) {
	mm_atomic.AddUint64(&mmDoWork.beforeDoWorkCounter, 1)
	defer mm_atomic.AddUint64(&mmDoWork.afterDoWorkCounter, 1)

	if mmDoWork.inspectFuncDoWork != nil {
		mmDoWork.inspectFuncDoWork(ctx, data)
	}

	mm_params := &ClientWorkerMockDoWorkParams{ctx, data}

	// Record call args
	mmDoWork.DoWorkMock.mutex.Lock()
	mmDoWork.DoWorkMock.callArgs = append(mmDoWork.DoWorkMock.callArgs, mm_params)
	mmDoWork.DoWorkMock.mutex.Unlock()

	for _, e := range mmDoWork.DoWorkMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.d1, e.results.err
		}
	}

	if mmDoWork.DoWorkMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDoWork.DoWorkMock.defaultExpectation.Counter, 1)
		mm_want := mmDoWork.DoWorkMock.defaultExpectation.params
		mm_got := ClientWorkerMockDoWorkParams{ctx, data}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDoWork.t.Errorf("ClientWorkerMock.DoWork got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmDoWork.DoWorkMock.defaultExpectation.results
		if mm_results == nil {
			mmDoWork.t.Fatal("No results are set for the ClientWorkerMock.DoWork")
		}
		return (*mm_results).d1, (*mm_results).err
	}
	if mmDoWork.funcDoWork != nil {
		return mmDoWork.funcDoWork(ctx, data)
	}
	mmDoWork.t.Fatalf("Unexpected call to ClientWorkerMock.DoWork. %v %v", ctx, data)
	return
}

// DoWorkAfterCounter returns a count of finished ClientWorkerMock.DoWork invocations
func (mmDoWork *ClientWorkerMock) DoWorkAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDoWork.afterDoWorkCounter)
}

// DoWorkBeforeCounter returns a count of ClientWorkerMock.DoWork invocations
func (mmDoWork *ClientWorkerMock) DoWorkBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDoWork.beforeDoWorkCounter)
}

// Calls returns a list of arguments used in each call to ClientWorkerMock.DoWork.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDoWork *mClientWorkerMockDoWork) Calls() []*ClientWorkerMockDoWorkParams {
	mmDoWork.mutex.RLock()

	argCopy := make([]*ClientWorkerMockDoWorkParams, len(mmDoWork.callArgs))
	copy(argCopy, mmDoWork.callArgs)

	mmDoWork.mutex.RUnlock()

	return argCopy
}

// MinimockDoWorkDone returns true if the count of the DoWork invocations corresponds
// the number of defined expectations
func (m *ClientWorkerMock) MinimockDoWorkDone() bool {
	for _, e := range m.DoWorkMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoWorkMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoWorkCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDoWork != nil && mm_atomic.LoadUint64(&m.afterDoWorkCounter) < 1 {
		return false
	}
	return true
}

// MinimockDoWorkInspect logs each unmet expectation
func (m *ClientWorkerMock) MinimockDoWorkInspect() {
	for _, e := range m.DoWorkMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ClientWorkerMock.DoWork with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DoWorkMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDoWorkCounter) < 1 {
		if m.DoWorkMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ClientWorkerMock.DoWork")
		} else {
			m.t.Errorf("Expected call to ClientWorkerMock.DoWork with params: %#v", *m.DoWorkMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDoWork != nil && mm_atomic.LoadUint64(&m.afterDoWorkCounter) < 1 {
		m.t.Error("Expected call to ClientWorkerMock.DoWork")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ClientWorkerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDoWorkInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ClientWorkerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ClientWorkerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDoWorkDone()
}