package xlog

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	rl := NewReqLogger(nil, ReqConfig{Level: LevelDebug})
	ctx := NewContext(context.Background(), rl)
	rl2, ok := FromContext(ctx)
	if !ok {
		t.Fatalf("rl is not in ctx")
	}
	if rl != rl2 {
		t.Fatalf("rl != rl2")
	}

	rl3 := FromContextSafe(context.Background())
	if rl3 == nil {
		t.Fail()
	}
}
