package echoutil

import (
	"testing"

	"github.com/ironzhang/superlib/testutil"
)

func TestLogger(t *testing.T) {
	log := NewLogger()

	log.Print("hello", "world")
	log.Printf("hello, %s", "world")
	log.Printj(map[string]interface{}{
		"hello": "world",
	})

	log.Debug("hello", "world")
	log.Debugf("hello, %s", "world")
	log.Debugj(map[string]interface{}{
		"hello": "world",
	})

	log.Info("hello", "world")
	log.Infof("hello, %s", "world")
	log.Infoj(map[string]interface{}{
		"hello": "world",
	})

	log.Warn("hello", "world")
	log.Warnf("hello, %s", "world")
	log.Warnj(map[string]interface{}{
		"hello": "world",
	})

	log.Error("hello", "world")
	log.Errorf("hello, %s", "world")
	log.Errorj(map[string]interface{}{
		"hello": "world",
	})

	testutil.RecoverPanic(func() {
		log.Panic("hello", "world")
	})
	testutil.RecoverPanic(func() {
		log.Panicf("hello, %s", "world")
	})
	testutil.RecoverPanic(func() {
		log.Panicj(map[string]interface{}{
			"hello": "world",
		})
	})

	//log.Fatal("hello", "world")
	//log.Fatalf("hello, %s", "world")
	//log.Fatalj(map[string]interface{}{
	//	"hello": "world",
	//})
}
