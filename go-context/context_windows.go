package gocontext

import (
	"context"
)

func WithCancelx() *Context {
	ctx, _ := context.WithCancel(context.TODO())
	//
	//// 创建事件句柄
	//event, err := windows.CreateEvent(nil, 1, 0, nil)
	//if err != nil {
	//	golog.Error("create event failed:", err)
	//	return WithParent(ctx)
	//}
	//
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			golog.Error(r)
	//		}
	//		windows.CloseHandle(event)
	//	}()
	//
	//	// 等待控制台控制事件
	//	kernel32 := windows.NewLazyDLL("kernel32.dll")
	//	setConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")
	//
	//	callback := windows.NewCallback(func(ctrlType uint) uintptr {
	//		switch ctrlType {
	//		case windows.CTRL_C_EVENT, // Ctrl+C
	//			windows.CTRL_BREAK_EVENT,    // Ctrl+Break
	//			windows.CTRL_CLOSE_EVENT,    // 关闭控制台窗口
	//			windows.CTRL_LOGOFF_EVENT,   // 用户注销
	//			windows.CTRL_SHUTDOWN_EVENT: // 系统关机
	//			cancel()
	//			windows.SetEvent(event)
	//			return 1 // 表示已处理事件
	//		}
	//		return 0
	//	})
	//
	//	_, _, _ = setConsoleCtrlHandler.Call(callback, 1)
	//
	//	// 等待事件触发
	//	windows.WaitForSingleObject(event, windows.INFINITE)
	//}()

	return WithParent(ctx)
}
