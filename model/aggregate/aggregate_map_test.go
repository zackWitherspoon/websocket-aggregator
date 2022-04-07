package aggregate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregate Test", func() {
	startTimestamp := int64(1536036818784)
	initialDuration := Duration{
		StartTime: startTimestamp,
		EndTime:   startTimestamp + 20,
	}

	Describe("Given a timestamp", func() {
		Context("When that timestamp is outside the duration window", func() {
			It("Should return false", func() {
				expected := false
				actual := initialDuration.Between(startTimestamp - 2000)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that timestamp is inside the duration window", func() {
			It("Should return true", func() {
				expected := true
				actual := initialDuration.Between(startTimestamp + 3)
				Expect(actual).To(Equal(expected))
			})
		})
		Context("When that timestamp is equal to the startTime", func() {
			It("Should return false", func() {
				expected := false
				actual := initialDuration.Between(startTimestamp)
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
