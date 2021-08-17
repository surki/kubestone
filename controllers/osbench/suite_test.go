package osbench

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOsbenchController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Osbench Controller Suite")
}
