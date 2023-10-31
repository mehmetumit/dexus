package stdlog

import (
	"io"
	"log"
)

type StdLog struct {
	l          *log.Logger
	debugLevel bool
}

func NewStdLog(debugLevel bool) *StdLog {
	log.Default()
	return &StdLog{
		//New(os.Stderr, "", log.LstdFlags)
		l:          log.Default(),
		debugLevel: debugLevel,
	}
}

func (sl *StdLog) SetDebugLevel(dl bool) {
	sl.debugLevel = dl
}

func (sl *StdLog) SetWriter(w io.Writer) {
	sl.l.SetOutput(w)
}
func (sl *StdLog) Error(v ...any) {
	sl.l.Println("[ERROR]", v)
}
func (sl *StdLog) Debug(v ...any) {
	if sl.debugLevel {
		sl.l.Println("[DEBUG]", v)
	}
}
func (sl *StdLog) Info(v ...any) {
	sl.l.Println("[INFO]", v)
}
func (sl *StdLog) Panic(v ...any) {
	sl.l.Println("[PANIC]", v)
}
