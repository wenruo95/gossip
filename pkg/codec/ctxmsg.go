package codec

import (
	"context"

	"go.uber.org/zap"
)

type ContextMsgKey string

var ctxKey ContextMsgKey = ContextMsgKey("ctx_msg_key")

type Msg interface {
	Context() context.Context

	WithBiz(biz string)
	WithCMD(cmd string)
	WithUID(uid uint64)
	WithTraceID(traceID string)

	Biz() string
	CMD() string
	UID() uint64
	TraceID() string

	Fileds() []zap.Field
}

func NewMsg(ctx context.Context) (context.Context, Msg) {
	msg := newMsg()
	newCtx := context.WithValue(ctx, ctxKey, msg)
	msg.context = newCtx
	return newCtx, msg
}

func EnsureMsg(ctx context.Context) (context.Context, Msg) {
	msgI := ctx.Value(ctxKey)
	if msg, ok := msgI.(*msg); ok {
		return ctx, msg
	}

	return NewMsg(ctx)
}

func Message(ctx context.Context) Msg {
	if msg, ok := ctx.Value(ctxKey).(*msg); ok {
		return msg
	}
	return &msg{context: ctx}
}

type msg struct {
	context context.Context

	biz     string
	cmd     string
	uid     uint64
	traceID string
}

func newMsg() *msg {
	return &msg{}
}

func (msg *msg) Context() context.Context {
	return msg.context
}

func (msg *msg) WithBiz(biz string) {
	msg.biz = biz
}

func (msg *msg) Biz() string {
	return msg.biz
}

func (msg *msg) WithCMD(cmd string) {
	msg.cmd = cmd
}

func (msg *msg) CMD() string {
	return msg.cmd
}

func (msg *msg) WithUID(uid uint64) {
	msg.uid = uid
}

func (msg *msg) UID() uint64 {
	return msg.uid
}

func (msg *msg) WithTraceID(traceID string) {
	msg.traceID = traceID
}

func (msg *msg) TraceID() string {
	return msg.traceID
}

func (msg *msg) Fileds() []zap.Field {
	return []zap.Field{
		zap.String("biz", msg.biz),
		zap.String("cmd", msg.cmd),
		zap.Uint64("uid", msg.uid),
		zap.String("trace_id", msg.traceID),
	}
}
