package mocks

import (
	"io"
	"log"
)

type MockLogger struct {
	l          *log.Logger
	debugLevel bool
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		l: log.Default(),
	}
}

func (m *MockLogger) SetDebugLevel(dl bool) {
	m.debugLevel = dl

}
func (m *MockLogger) SetWriter(w io.Writer) {
	m.l.SetOutput(w)

}
func (m *MockLogger) Error(v ...any) {
	m.l.Println("[ERROR]", v)

}
func (m *MockLogger) Info(v ...any) {
	m.l.Println("[INFO]", v)

}
func (m *MockLogger) Debug(v ...any) {
	if m.debugLevel {
		m.l.Println("[DEBUG]", v)
	}

}
func (m *MockLogger) Panic(v ...any) {
	m.l.Println("[PANIC]", v)

}
