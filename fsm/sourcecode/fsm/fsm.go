// Copyright (c) 2013 - Max Persson <max@looplab.se>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fsm implements a finite state machine.
//
// It is heavily based on two FSM implementations:
//
// Javascript Finite State Machine
// https://github.com/jakesgordon/javascript-state-machine
//
// Fysom for Python
// https://github.com/oxplot/fysom (forked at https://github.com/mriehl/fysom)
package fsm

import (
	"context"
	"strings"
	"sync"
)

// transitioner 是 FSM 的状态转换函数接口。
type transitioner interface {
	transition(*FSM) error
}

// FSM 是持有「当前状态」的状态机。
// 必须使用 NewFSM 创建才能正常工作。
type FSM struct {
	// FSM 当前状态
	current string

	// transitions 将「事件和原状态」映射到「目标状态」。
	// key: event + src
	// val: dst
	transitions map[eKey]string

	// callbacks 将「回调类型和目标」映射到「回调函数」。
	// key: callbackType + target
	// val: callback（事件触发时调用的回调函数）
	callbacks map[cKey]Callback

	// transition 是内部状态转换函数，可以直接使用，也可以在异步状态转换时调用。
	transition func()
	// transitionerObj 用于调用 FSM 的 transition() 函数。
	transitionerObj transitioner

	// stateMu 保护对当前状态的访问。
	stateMu sync.RWMutex
	// eventMu 保护对 Event() 和 Transition() 两个函数的调用。
	eventMu sync.Mutex

	// metadata 可以用来存储和加载可能跨事件使用的数据
	// 使用 SetMetadata() 和 Metadata() 方法来存储和加载数据。
	metadata map[string]interface{}
	// metadataMu 保护对元数据的访问。
	metadataMu sync.RWMutex
}

// EventDesc represents an event when initializing the FSM.
//
// The event can have one or more source states that is valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
// EventDesc 表示初始化 FSM 时的一个事件。
type EventDesc struct {
	// Name is the event name used when calling for a transition.
	Name string

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	Src []string

	// Dst is the destination state that the FSM will be in if the transition
	// succeeds.
	Dst string
}

// Callback is a function type that callbacks should use. Event is the current
// event info as the callback happens.
type Callback func(context.Context, *Event)

// Events is a shorthand for defining the transition map in NewFSM.
type Events []EventDesc

// Callbacks is a shorthand for defining the callbacks in NewFSM.
type Callbacks map[string]Callback

// NewFSM 通过事件和回调函数构造一个有限状态机
//
// 事件和状态转换规则通过 Events 切片（slice）定义，每个 Event 对应一个或多个
// 从 Event.Src 到 Event.Dst 的内部转换规则
//
// 回调函数通过 Callbacks 映射表添加，键名按以下规则解析，调用顺序如下：
// 1. before_<EVENT>    - 在特定事件前执行（如 before_open）
// 2. before_event      - 全局事件前置钩子
// 3. leave_<OLD_STATE> - 离开旧状态前执行（如 leave_closed）
// 4. leave_state       - 全局状态离开钩子
// 5. enter_<NEW_STATE> - 进入新状态后执行（如 enter_open）
// 6. enter_state       - 全局状态进入钩子
// 7. after_<EVENT>     - 事件完成后执行（如 after_open）
// 8. after_event       - 全局事件后置钩子
//
// 对于最常用的回调函数也有两个简短版本的实现。
// 它们只是事件或状态的名称：
// 1. <NEW_STATE>       - 等效 5. enter_<NEW_STATE>（如 closed => enter_closed）
// 2. <EVENT>           - 等效 7. after_<EVENT>（如 close => after_close）
//
// 若同时定义短格式与完整格式回调，由于 Go map 的无序性，
// 最终生效的回调函数版本将不确定。当前实现不做重复键检查。
//
// 注册事件示例：
//
//	fsm.Events{
//		{Name: "open", Src: []string{"closed"}, Dst: "open"},
//		{Name: "close", Src: []string{"open"}, Dst: "closed"},
//	},
func NewFSM(initial string, events []EventDesc, callbacks map[string]Callback) *FSM {
	// 构造有限状态机 FSM
	f := &FSM{
		transitionerObj: &transitionerStruct{},        // 状态转换器，使用默认实现
		current:         initial,                      // 当前状态
		transitions:     make(map[eKey]string),        // 存储「事件和原状态」到「目标状态」的转换规则映射
		callbacks:       make(map[cKey]Callback),      // 回调函数映射表
		metadata:        make(map[string]interface{}), // 元信息
	}

	// 构建 f.transitions map，并且存储所有的「事件」和「状态」集合
	allEvents := make(map[string]bool) // 存储所有事件的集合
	allStates := make(map[string]bool) // 存储所有状态的集合
	for _, e := range events {         // 遍历事件列表，提取并存储所有事件和状态
		for _, src := range e.Src {
			f.transitions[eKey{e.Name, src}] = e.Dst
			allStates[src] = true
			allStates[e.Dst] = true
		}
		allEvents[e.Name] = true
	}

	// Map all callbacks to events/states.
	// 提取「回调函数」到「事件和原状态」的映射关系，并注册到 callbacks
	// 示例：
	// fsm.Callbacks{
	//     "enter_state": func(_ context.Context, e *fsm.Event) { d.enterState(e) },
	// }
	for name, fn := range callbacks {
		var target string    // 目标：状态/事件
		var callbackType int // 回调类型（决定了调用顺序）

		// 根据回调函数名称前缀分类
		switch {
		// 事件触发前执行
		case strings.HasPrefix(name, "before_"):
			target = strings.TrimPrefix(name, "before_")
			if target == "event" { // 全局事件前置钩子（任何事件触发都会调用，如用于日志记录场景）
				target = "" // 将 target 置空
				callbackType = callbackBeforeEvent
			} else if _, ok := allEvents[target]; ok { // 在特定事件前执行
				callbackType = callbackBeforeEvent
			}
		// 离开当前状态前执行
		case strings.HasPrefix(name, "leave_"):
			target = strings.TrimPrefix(name, "leave_")
			if target == "state" { // 全局状态离开钩子
				target = ""
				callbackType = callbackLeaveState
			} else if _, ok := allStates[target]; ok { // 离开旧状态前执行
				callbackType = callbackLeaveState
			}
		// 进入新状态后执行
		case strings.HasPrefix(name, "enter_"):
			target = strings.TrimPrefix(name, "enter_")
			if target == "state" { // 全局状态进入钩子
				target = ""
				callbackType = callbackEnterState
			} else if _, ok := allStates[target]; ok { // 进入新状态后执行
				callbackType = callbackEnterState
			}
		// 事件完成后执行
		case strings.HasPrefix(name, "after_"):
			target = strings.TrimPrefix(name, "after_")
			if target == "event" { // 全局事件后置钩子
				target = ""
				callbackType = callbackAfterEvent
			} else if _, ok := allEvents[target]; ok { // 事件完成后执行
				callbackType = callbackAfterEvent
			}
		// 处理未加前缀的回调（简短版本）
		default:
			target = name                       // 状态/事件
			if _, ok := allStates[target]; ok { // 如果 target 为某个状态，则 callbackType 会置为与 enter_[target] 相同，即二者等价
				callbackType = callbackEnterState
			} else if _, ok := allEvents[target]; ok { // 如果 target 为某个事件，则 callbackType 会置为与 after_[target] 相同，即二者等价
				callbackType = callbackAfterEvent
			}
		}

		// 记录 callbacks map
		if callbackType != callbackNone {
			// key: callbackType（用于决定执行顺序） + target（如果是全局钩子，则 target 为空，否则，target 为状态/事件）
			// val: 事件触发时需要执行的回调函数
			f.callbacks[cKey{target, callbackType}] = fn
		}
	}

	return f
}

// Current 返回 FSM 的当前状态。
func (f *FSM) Current() string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.current
}

// Is 判断 FSM 当前状态是否为指定状态。
func (f *FSM) Is(state string) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return state == f.current
}

// SetState 将 FSM 从当前状态转移到指定状态。
// 此调用不触发任何回调函数（如果定义）。
func (f *FSM) SetState(state string) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = state
}

// Can 判断 FSM 在当前状态下，是否可以触发指定事件，如果可以，则返回 true。
func (f *FSM) Can(event string) bool {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	_, ok := f.transitions[eKey{event, f.current}]
	return ok && (f.transition == nil)
}

// AvailableTransitions 返回当前状态下可用的转换列表。
func (f *FSM) AvailableTransitions() []string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	var transitions []string
	for key := range f.transitions {
		if key.src == f.current {
			transitions = append(transitions, key.event)
		}
	}
	return transitions
}

// Cannot returns true if event can not occur in the current state.
// It is a convenience method to help code read nicely.
func (f *FSM) Cannot(event string) bool {
	return !f.Can(event)
}

// Metadata 返回存储在元信息中的值
func (f *FSM) Metadata(key string) (interface{}, bool) {
	f.metadataMu.RLock()
	defer f.metadataMu.RUnlock()
	dataElement, ok := f.metadata[key]
	return dataElement, ok
}

// SetMetadata 存储 key、val 到元信息中
func (f *FSM) SetMetadata(key string, dataValue interface{}) {
	f.metadataMu.Lock()
	defer f.metadataMu.Unlock()
	f.metadata[key] = dataValue
}

// DeleteMetadata 从元信息中删除指定 key 对应的数据
func (f *FSM) DeleteMetadata(key string) {
	f.metadataMu.Lock()
	delete(f.metadata, key)
	f.metadataMu.Unlock()
}

// Event 通过指定事件名称触发状态转换
//
// 此方法接收可变数量的参数，这些参数将传递给已定义的回调函数（如果有）
//
// 如果状态更改正常，将返回 nil，否则返回以下错误之一：
//
// - InTransitionError：事件 X 不合时宜，因为之前的转换尚未完成
// - InvalidEventError：事件 X 在当前状态 Y 下不适用
// - UnknownEventError：事件 X 不存在
// - InternalError：状态转换期间的内部错误（理论上此错误在此情况下不应发生，表明存在内部错误，其他错误是可预知的，所以使用 SentinelError）
func (f *FSM) Event(ctx context.Context, event string, args ...interface{}) error {
	f.eventMu.Lock() // 事件互斥锁锁定

	// 为了始终解锁事件互斥锁（eventMu），此处添加了 defer 防止状态转换完成后执行 enter/after 回调时仍持有锁；
	// 因为这些回调可能触发新的状态转换，故在下方代码中需要显式解锁
	var unlocked bool // 标记是否已经解锁
	defer func() {
		if !unlocked { // 如果下方的逻辑已经显式操作过解锁，defer 中无需重复解锁
			f.eventMu.Unlock()
		}
	}()

	f.stateMu.RLock() // 获取状态读锁
	defer f.stateMu.RUnlock()

	// NOTE: 之前的转换尚未完成
	if f.transition != nil {
		// 上一次状态转换还未完成，返回"前一个转换未完成"错误
		return InTransitionError{event}
	}

	// NOTE: 事件 event 在当前状态 current 下是否适用，即是否在 transitions 表中
	dst, ok := f.transitions[eKey{event, f.current}]
	if !ok { // 无效事件
		for ekey := range f.transitions {
			if ekey.event == event {
				// 事件和当前状态不对应
				return InvalidEventError{event, f.current}
			}
		}
		// 未定义的事件
		return UnknownEventError{event}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// 构造一个事件对象
	e := &Event{f, event, f.current, dst, nil, args, false, false, cancel}

	// NOTE: 执行 before 钩子
	err := f.beforeEventCallbacks(ctx, e)
	if err != nil {
		return err
	}

	// NOTE: 当前状态等于目标状态，无需转换
	if f.current == dst {
		f.stateMu.RUnlock()
		defer f.stateMu.RLock()
		f.eventMu.Unlock()
		unlocked = true
		// NOTE: 执行 after 钩子
		f.afterEventCallbacks(ctx, e)
		return NoTransitionError{e.Err}
	}

	// 定义状态转换闭包函数
	transitionFunc := func(ctx context.Context, async bool) func() {
		return func() {
			if ctx.Err() != nil {
				if e.Err == nil {
					e.Err = ctx.Err()
				}
				return
			}

			f.stateMu.Lock()
			f.current = dst    // 状态转换
			f.transition = nil // NOTE: 标记状态转换完成
			f.stateMu.Unlock()

			// 显式解锁 eventMu 事件互斥锁，允许 enterStateCallbacks 回调函数触发新的状态转换操作（避免死锁）
			// 对于异步状态转换，无需显式解锁，锁已在触发异步操作时释放
			if !async {
				f.eventMu.Unlock()
				unlocked = true
			}
			// NOTE: 执行 enter 钩子
			f.enterStateCallbacks(ctx, e)
			// NOTE: 执行 after 钩子
			f.afterEventCallbacks(ctx, e)
		}
	}

	// 记录状态转换函数（这里标记为同步转换）
	f.transition = transitionFunc(ctx, false)

	// NOTE: 执行 leave 钩子
	if err = f.leaveStateCallbacks(ctx, e); err != nil {
		if _, ok := err.(CanceledError); ok {
			f.transition = nil // NOTE: 如果通过 ctx 取消了，则标记为 nil，无需转换
		} else if asyncError, ok := err.(AsyncError); ok { // NOTE: 如果是 AsyncError，说明是异步转换
			// 为异步操作创建独立上下文，以便异步状态转换正常工作
			// 这个新的 ctx 实际上已经脱离了原始 ctx，原 ctx 取消不会影响当前 ctx
			// 不过新的 ctx 保留了原始 ctx 的值，所有通过 ctx 传递的值还可以继续使用
			ctx, cancel := uncancelContext(ctx)
			e.cancelFunc = cancel                    // 绑定新取消函数
			asyncError.Ctx = ctx                     // 传递新上下文
			asyncError.CancelTransition = cancel     // 暴露取消接口
			f.transition = transitionFunc(ctx, true) // NOTE: 标记为异步转换状态
			// NOTE: 如果是异步转换，直接返回，不会同步调用 f.doTransition()，需要用户手动调用 f.Transition() 来触发状态转换
			return asyncError
		}
		return err
	}

	// Perform the rest of the transition, if not asynchronous.
	f.stateMu.RUnlock()
	defer f.stateMu.RLock()
	err = f.doTransition() // NOTE: 执行状态转换逻辑，即调用 f.transition()
	if err != nil {
		return InternalError{}
	}

	return e.Err
}

// Transition wraps transitioner.transition.
func (f *FSM) Transition() error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()
	return f.doTransition()
}

// doTransition wraps transitioner.transition.
func (f *FSM) doTransition() error {
	return f.transitionerObj.transition(f)
}

// 状态转换接口的默认实现
type transitionerStruct struct{}

// Transition completes an asynchronous state change.
//
// The callback for leave_<STATE> must previously have called Async on its
// event to have initiated an asynchronous state transition.
func (t transitionerStruct) transition(f *FSM) error {
	if f.transition == nil {
		return NotInTransitionError{}
	}
	f.transition()
	return nil
}

// beforeEventCallbacks calls the before_ callbacks, first the named then the
// general version.
func (f *FSM) beforeEventCallbacks(ctx context.Context, e *Event) error {
	if fn, ok := f.callbacks[cKey{e.Event, callbackBeforeEvent}]; ok {
		fn(ctx, e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	if fn, ok := f.callbacks[cKey{"", callbackBeforeEvent}]; ok {
		fn(ctx, e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	return nil
}

// leaveStateCallbacks calls the leave_ callbacks, first the named then the
// general version.
func (f *FSM) leaveStateCallbacks(ctx context.Context, e *Event) error {
	if fn, ok := f.callbacks[cKey{f.current, callbackLeaveState}]; ok {
		fn(ctx, e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async { // NOTE: 异步信号
			return AsyncError{Err: e.Err}
		}
	}
	if fn, ok := f.callbacks[cKey{"", callbackLeaveState}]; ok {
		fn(ctx, e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{Err: e.Err}
		}
	}
	return nil
}

// enterStateCallbacks calls the enter_ callbacks, first the named then the
// general version.
func (f *FSM) enterStateCallbacks(ctx context.Context, e *Event) {
	if fn, ok := f.callbacks[cKey{f.current, callbackEnterState}]; ok {
		fn(ctx, e)
	}
	if fn, ok := f.callbacks[cKey{"", callbackEnterState}]; ok {
		fn(ctx, e)
	}
}

// afterEventCallbacks calls the after_ callbacks, first the named then the
// general version.
func (f *FSM) afterEventCallbacks(ctx context.Context, e *Event) {
	if fn, ok := f.callbacks[cKey{e.Event, callbackAfterEvent}]; ok {
		fn(ctx, e)
	}
	if fn, ok := f.callbacks[cKey{"", callbackAfterEvent}]; ok {
		fn(ctx, e)
	}
}

const (
	// 未设置回调
	callbackNone int = iota
	// 事件触发前执行的回调
	callbackBeforeEvent
	// 离开旧状态前执行的回调
	callbackLeaveState
	// 进入新状态是执行的回调
	callbackEnterState
	// 事件完成时执行的回调
	callbackAfterEvent
)

// cKey is a struct key used for keeping the callbacks mapped to a target.
type cKey struct {
	// target is either the name of a state or an event depending on which
	// callback type the key refers to. It can also be "" for a non-targeted
	// callback like before_event.
	target string

	// callbackType is the situation when the callback will be run.
	callbackType int
}

// eKey is a struct key used for storing the transition map.
type eKey struct {
	// event is the name of the event that the keys refers to.
	event string

	// src is the source from where the event can transition.
	src string
}
