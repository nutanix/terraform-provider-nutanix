package monitoringv2_test

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Running monitoring v2 tests")
	os.Exit(m.Run())
}
