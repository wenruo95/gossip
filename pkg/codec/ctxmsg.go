package codec

import (
	"context"

	"go.uber.org/zap"
)

type ContextMsgKey string

var ctxKey ContextMsgKey = ContextMsgKey("ctx_msg_key")

type Msg interface {
	Context() context.Context

	WithCMD(cmd string)
	WithUID(uid uint64)
	WithTraceID(traceID string)
	WithSourceID(sourceID string)

	CMD() string
	UID() uint64
	TraceID() string
	SourceID() string

	WithFields(field ...zap.Field)
	Fields() []zap.Field
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

	cmd      string
	uid      uint64
	sourceID string
	traceID  string
	fields   []zap.Field
}

func newMsg() *msg {
	return &msg{fields: make([]zap.Field, 0)}
}

func (msg *msg) Context() context.Context {
	return msg.context
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

func (msg *msg) WithSourceID(sourceID string) {
	msg.sourceID = sourceID
}

func (msg *msg) SourceID() string {
	return msg.sourceID
}

func (msg *msg) WithTraceID(traceID string) {
	msg.traceID = traceID
}

func (msg *msg) TraceID() string {
	return msg.traceID
}

func (msg *msg) WithFields(fields ...zap.Field) {
	msg.fields = append(msg.fields, fields...)
}

func (msg *msg) Fields() []zap.Field {
	if len(msg.fields) == 0 {
		return []zap.Field{
			zap.String("cmd", msg.cmd),
			zap.Uint64("uid", msg.uid),
			zap.String("trace_id", msg.traceID),
			zap.String("source_id", msg.SourceID()),
		}
	}

	f := make([]zap.Field, 0)
	f = append(f, msg.fields...)
	if len(msg.cmd) > 0 {
		f = append(f, zap.String("cmd", msg.cmd))
	}
	if msg.uid > 0 {
		f = append(f, zap.Uint64("uid", msg.uid))
	}
	if len(msg.traceID) > 0 {
		f = append(f, zap.String("trace_id", msg.traceID))
	}
	if len(msg.sourceID) > 0 {
		f = append(f, zap.String("source_id", msg.SourceID()))
	}
	return f
}
