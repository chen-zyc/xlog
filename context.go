package xlog

import "context"

type key int

const logKey key = 0

// NewContext returns a new context, and puts rl inside.
func NewContext(parent context.Context, rl ReqLogger) context.Context {
	return context.WithValue(parent, logKey, rl)
}

// FromContext takes ReqLogger from ctx. If there is no ReqLogger inside, the `ok` will be false.
func FromContext(ctx context.Context) (rl ReqLogger, ok bool) {
	v := ctx.Value(logKey)
	if v == nil {
		return nil, false
	}
	rl, ok = v.(ReqLogger)
	return
}

// FromContextSafe is the same as FromContext, except if there is no ReqLogger, return a new one.
func FromContextSafe(ctx context.Context) (rl ReqLogger) {
	rl, ok := FromContext(ctx)
	if !ok {
		rl = NewReqLogger(nil, ReqConfig{})
	}
	return
}
