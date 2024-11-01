// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	types "github.com/rollkit/go-execution/types"
)

// MockExecute is an autogenerated mock type for the Execute type
type MockExecute struct {
	mock.Mock
}

type MockExecute_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExecute) EXPECT() *MockExecute_Expecter {
	return &MockExecute_Expecter{mock: &_m.Mock}
}

// ExecuteTxs provides a mock function with given fields: txs, blockHeight, timestamp, prevStateRoot
func (_m *MockExecute) ExecuteTxs(txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash) (types.Hash, uint64, error) {
	ret := _m.Called(txs, blockHeight, timestamp, prevStateRoot)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteTxs")
	}

	var r0 types.Hash
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func([]types.Tx, uint64, time.Time, types.Hash) (types.Hash, uint64, error)); ok {
		return rf(txs, blockHeight, timestamp, prevStateRoot)
	}
	if rf, ok := ret.Get(0).(func([]types.Tx, uint64, time.Time, types.Hash) types.Hash); ok {
		r0 = rf(txs, blockHeight, timestamp, prevStateRoot)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func([]types.Tx, uint64, time.Time, types.Hash) uint64); ok {
		r1 = rf(txs, blockHeight, timestamp, prevStateRoot)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func([]types.Tx, uint64, time.Time, types.Hash) error); ok {
		r2 = rf(txs, blockHeight, timestamp, prevStateRoot)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockExecute_ExecuteTxs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteTxs'
type MockExecute_ExecuteTxs_Call struct {
	*mock.Call
}

// ExecuteTxs is a helper method to define mock.On call
//   - txs []types.Tx
//   - blockHeight uint64
//   - timestamp time.Time
//   - prevStateRoot types.Hash
func (_e *MockExecute_Expecter) ExecuteTxs(txs interface{}, blockHeight interface{}, timestamp interface{}, prevStateRoot interface{}) *MockExecute_ExecuteTxs_Call {
	return &MockExecute_ExecuteTxs_Call{Call: _e.mock.On("ExecuteTxs", txs, blockHeight, timestamp, prevStateRoot)}
}

func (_c *MockExecute_ExecuteTxs_Call) Run(run func(txs []types.Tx, blockHeight uint64, timestamp time.Time, prevStateRoot types.Hash)) *MockExecute_ExecuteTxs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]types.Tx), args[1].(uint64), args[2].(time.Time), args[3].(types.Hash))
	})
	return _c
}

func (_c *MockExecute_ExecuteTxs_Call) Return(updatedStateRoot types.Hash, maxBytes uint64, err error) *MockExecute_ExecuteTxs_Call {
	_c.Call.Return(updatedStateRoot, maxBytes, err)
	return _c
}

func (_c *MockExecute_ExecuteTxs_Call) RunAndReturn(run func([]types.Tx, uint64, time.Time, types.Hash) (types.Hash, uint64, error)) *MockExecute_ExecuteTxs_Call {
	_c.Call.Return(run)
	return _c
}

// GetTxs provides a mock function with given fields:
func (_m *MockExecute) GetTxs() ([]types.Tx, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetTxs")
	}

	var r0 []types.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]types.Tx, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []types.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExecute_GetTxs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTxs'
type MockExecute_GetTxs_Call struct {
	*mock.Call
}

// GetTxs is a helper method to define mock.On call
func (_e *MockExecute_Expecter) GetTxs() *MockExecute_GetTxs_Call {
	return &MockExecute_GetTxs_Call{Call: _e.mock.On("GetTxs")}
}

func (_c *MockExecute_GetTxs_Call) Run(run func()) *MockExecute_GetTxs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExecute_GetTxs_Call) Return(_a0 []types.Tx, _a1 error) *MockExecute_GetTxs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExecute_GetTxs_Call) RunAndReturn(run func() ([]types.Tx, error)) *MockExecute_GetTxs_Call {
	_c.Call.Return(run)
	return _c
}

// InitChain provides a mock function with given fields: genesisTime, initialHeight, chainID
func (_m *MockExecute) InitChain(genesisTime time.Time, initialHeight uint64, chainID string) (types.Hash, uint64, error) {
	ret := _m.Called(genesisTime, initialHeight, chainID)

	if len(ret) == 0 {
		panic("no return value specified for InitChain")
	}

	var r0 types.Hash
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func(time.Time, uint64, string) (types.Hash, uint64, error)); ok {
		return rf(genesisTime, initialHeight, chainID)
	}
	if rf, ok := ret.Get(0).(func(time.Time, uint64, string) types.Hash); ok {
		r0 = rf(genesisTime, initialHeight, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func(time.Time, uint64, string) uint64); ok {
		r1 = rf(genesisTime, initialHeight, chainID)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(time.Time, uint64, string) error); ok {
		r2 = rf(genesisTime, initialHeight, chainID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockExecute_InitChain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitChain'
type MockExecute_InitChain_Call struct {
	*mock.Call
}

// InitChain is a helper method to define mock.On call
//   - genesisTime time.Time
//   - initialHeight uint64
//   - chainID string
func (_e *MockExecute_Expecter) InitChain(genesisTime interface{}, initialHeight interface{}, chainID interface{}) *MockExecute_InitChain_Call {
	return &MockExecute_InitChain_Call{Call: _e.mock.On("InitChain", genesisTime, initialHeight, chainID)}
}

func (_c *MockExecute_InitChain_Call) Run(run func(genesisTime time.Time, initialHeight uint64, chainID string)) *MockExecute_InitChain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time), args[1].(uint64), args[2].(string))
	})
	return _c
}

func (_c *MockExecute_InitChain_Call) Return(stateRoot types.Hash, maxBytes uint64, err error) *MockExecute_InitChain_Call {
	_c.Call.Return(stateRoot, maxBytes, err)
	return _c
}

func (_c *MockExecute_InitChain_Call) RunAndReturn(run func(time.Time, uint64, string) (types.Hash, uint64, error)) *MockExecute_InitChain_Call {
	_c.Call.Return(run)
	return _c
}

// SetFinal provides a mock function with given fields: blockHeight
func (_m *MockExecute) SetFinal(blockHeight uint64) error {
	ret := _m.Called(blockHeight)

	if len(ret) == 0 {
		panic("no return value specified for SetFinal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(blockHeight)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExecute_SetFinal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetFinal'
type MockExecute_SetFinal_Call struct {
	*mock.Call
}

// SetFinal is a helper method to define mock.On call
//   - blockHeight uint64
func (_e *MockExecute_Expecter) SetFinal(blockHeight interface{}) *MockExecute_SetFinal_Call {
	return &MockExecute_SetFinal_Call{Call: _e.mock.On("SetFinal", blockHeight)}
}

func (_c *MockExecute_SetFinal_Call) Run(run func(blockHeight uint64)) *MockExecute_SetFinal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint64))
	})
	return _c
}

func (_c *MockExecute_SetFinal_Call) Return(_a0 error) *MockExecute_SetFinal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExecute_SetFinal_Call) RunAndReturn(run func(uint64) error) *MockExecute_SetFinal_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExecute creates a new instance of MockExecute. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExecute(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExecute {
	mock := &MockExecute{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
