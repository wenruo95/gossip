package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	options := make([]zap.Option, 0)
	logger, err := zap.NewProduction(options...)
	if err != nil {
		t.Errorf("error:%v", err)
		return
	}

	sugar := logger.Sugar()
	sugar.Infof("xx")
	logger.Sync()

	time.Sleep(time.Second)
}

func TestZapLog(t *testing.T) {
	InitLogger(&Config{})
	core2 := zapCore.With([]zapcore.Field{zap.String("cmd", "test")})

	l := zap.New(core2, zap.AddCaller())
	l.Info("hello")
}

func TestStdLog(t *testing.T) {
	a, b, c := 10, "char_b", 1.332
	Debug("hello Debug")
	Debugf("hello Debug a:%v b:%v c:%v", a, b, c)
	Debugw("hello Debug", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	Info("hello Info")
	Infof("hello Info a:%v b:%v c:%v", a, b, c)
	Infow("hello Info", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	Warn("hello Warn")
	Warnf("hello Warn a:%v b:%v c:%v", a, b, c)
	Warnw("hello Warn", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	Error("hello Error")
	Errorf("hello Error a:%v b:%v c:%v", a, b, c)
	Errorw("hello Error", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	//Fatal("hello Fatal")
	//Fatalf("hello Fatal a:%v b:%v c:%v", a, b, c)
	//Fatalw("hello Fatal", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))
}

func TestZLogger(t *testing.T) {
	a, b, c := 10, "char_b", 1.332
	l := NewZLogger(ZapLogger(zap.AddCallerSkip(1)))
	l.Debug("hello Debug")
	l.Debugf("hello Debug a:%v b:%v c:%v", a, b, c)
	l.Debugw("hello Debug", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	l.Info("hello Info")
	l.Infof("hello Info a:%v b:%v c:%v", a, b, c)
	l.Infow("hello Info", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	l.Warn("hello Warn")
	l.Warnf("hello Warn a:%v b:%v c:%v", a, b, c)
	l.Warnw("hello Warn", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	l.Error("hello Error")
	l.Errorf("hello Error a:%v b:%v c:%v", a, b, c)
	l.Errorw("hello Error", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))

	//Fatal("hello Fatal")
	//Fatalf("hello Fatal a:%v b:%v c:%v", a, b, c)
	//Fatalw("hello Fatal", zap.Int("a", a), zap.String("b", b), zap.Float64("c", c))
}
