package mprof

import (
	"time"
	"testing"
	"os"
)


func TestNilNew(t *testing.T) {
	m := NewMprof("")
	if m != nil {t.Error("nil create not nil")}
}

func TestNew(t *testing.T) {
	m := NewMprof("tmprof")
	if m == nil {t.Error("create nil")}
}

func TestDump(t *testing.T) {
	m := NewMprof("tmprof")
	if m == nil {t.Error("create nil")}
	t.Log("create nil")
	m.Dump("test1")
	_, er := os.Stat("tmprof/test1")
	if er == nil {t.Error("no dump file")}
}

func TestIntervalDump(t *testing.T) {
	m := NewMprof("tmprof")
	if m == nil {t.Error("create nil")}
	m.StartDumps(10e5)
	c := time.After(10e6)
	<-c
	m.StopDumps()
}

