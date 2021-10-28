package gsignal

/*
https://github.com/vrecan/death
https://github.com/codemodus/sigmon
*/

import (
	"context"
	"github.com/SentimensRG/sigctx"
)

// 避免恶意结束进程
// 允许带密码方式结束进程

func WaitExitSignal() {
	ctx := sigctx.New() // returns a regular context.Context

	// Against this simple pattern, your goroutines are guaranteed to terminate correctly
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	<-ctx.Done() // will unblock on SIGINT and SIGTERM
}
