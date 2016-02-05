package main_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

type FakePersistenceBackup struct {
	ErrFake       error
	DumpCallCount int
	SpyImport     string
}

func (s *FakePersistenceBackup) Dump(w io.Writer) (err error) {
	s.DumpCallCount++
	return s.ErrFake
}

func (s *FakePersistenceBackup) Import(r io.Reader) (err error) {
	return
}
