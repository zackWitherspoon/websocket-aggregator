package aggregate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAggregate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aggregate Suite")
}
