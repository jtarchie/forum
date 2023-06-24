package cache_test

import (
	"errors"
	"sync"
	"time"

	"github.com/jtarchie/forum/cache"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Func", func() {
	It("passes the argument and returns the value", func() {
		invoked := false
		fun := func(value int) (bool, error) {
			Expect(value).To(Equal(10))
			invoked = true

			return true, nil
		}

		cached := cache.NewFunc(fun, time.Second)

		result, err := cached.Invoke(10)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeTrue())
		Expect(invoked).To(BeTrue())
	})

	It("returns the error", func() {
		someError := errors.New("some error")

		fun := func(_ int) (bool, error) {
			return false, someError
		}

		cached := cache.NewFunc(fun, time.Second)

		_, err := cached.Invoke(10)
		Expect(err).To(HaveOccurred())
	})

	When("waiting for a ttl", func() {
		It("does not invoke the original function", func() {
			invoked := 0
			fun := func(_ int) (int, error) {
				invoked += 1

				return invoked, nil
			}

			cached := cache.NewFunc(fun, 100*time.Millisecond)

			result, err := cached.Invoke(100)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(1))

			Consistently(func() int {
				result, _ := cached.Invoke(100)

				return result
			}, "50ms").Should(Equal(1))

			Eventually(func() int {
				result, _ := cached.Invoke(100)

				return result
			}).Should(Equal(2))

			Consistently(func() int {
				result, _ := cached.Invoke(100)

				return result
			}).Should(BeNumerically(">=", 2))
		})
	})

	When("being accessed across go routines", func() {
		It("doesn't have race conditions", func() {
			invoked := 0
			fun := func(_ int) (int, error) {
				invoked += 1

				return invoked, nil
			}

			cached := cache.NewFunc(fun, time.Millisecond)

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()

				for start := time.Now(); time.Since(start) < 10*time.Millisecond; {
					_, _ = cached.Invoke(1)
				}
			}()

			go func() {
				defer wg.Done()

				for start := time.Now(); time.Since(start) < 10*time.Millisecond; {
					_, _ = cached.Invoke(1)
				}
			}()

			wg.Wait()
		})
	})
})
