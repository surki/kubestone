package nighthawk

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNighthawkController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nighthawk Controller Suite")
}
