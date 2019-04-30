package multiple_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/delegatebot/handler"
	. "github.com/dpb587/slack-delegate-bot/delegatebot/handler/handlers/multiple"
	"github.com/dpb587/slack-delegate-bot/delegatebot/handler/handlers/single"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditionfakes"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegatefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var subject Handler
	var firstHandler, secondHandler single.Handler
	var msg message.Message
	var del *delegatefakes.FakeDelegator

	BeforeEach(func() {
		del = &delegatefakes.FakeDelegator{}
		firstHandler = single.Handler{
			Delegator: del,
			Options: handler.Options{
				EmptyMessage: "fake-empty-message",
			},
		}
		secondHandler = single.Handler{
			Delegator: del,
			Options: handler.Options{
				EmptyMessage: "other-fake-empty-message",
			},
		}

		subject = Handler{
			Handlers: []handler.Handler{firstHandler, secondHandler},
		}
	})

	Describe("Execute", func() {
		Context("delegate errors", func() {
			BeforeEach(func() {
				del.DelegateReturns(nil, errors.New("fake-err1"))
			})

			It("errors", func() {
				_, err := subject.Execute(&msg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err1"))
			})
		})

		Context("with delegates", func() {
			BeforeEach(func() {
				del.DelegateReturns([]delegate.Delegate{
					delegate.Literal{Text: "something"},
				}, nil)
			})

			It("returns the delegates", func() {
				res, err := subject.Execute(&msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Delegates).To(ConsistOf(delegate.Literal{Text: "something"}))
			})
		})

		It("configures empty message", func() {
			res, err := subject.Execute(&msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Delegates).To(HaveLen(0))
			Expect(res.EmptyMessage).To(Equal("fake-empty-message"))
		})
	})

	Describe("IsApplicable", func() {
		var condition *conditionfakes.FakeCondition

		Context("condition configured", func() {
			BeforeEach(func() {
				condition = &conditionfakes.FakeCondition{}
				firstHandler.Condition = condition
				secondHandler.Condition = condition
				subject.Handlers = []handler.Handler{firstHandler, secondHandler}
			})

			Context("true", func() {
				BeforeEach(func() {
					condition.EvaluateReturns(true, nil)
				})

				It("applies", func() {
					b, err := subject.IsApplicable(msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(b).To(BeTrue())
				})
			})

			Context("false", func() {
				BeforeEach(func() {
					condition.EvaluateReturns(false, nil)
				})

				It("does not apply", func() {
					b, err := subject.IsApplicable(msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(b).To(BeFalse())
				})
			})
		})

		It("always applies", func() {
			b, err := subject.IsApplicable(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(BeTrue())
		})
	})
})
