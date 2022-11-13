package appstore

import (
	"github.com/golang/mock/gomock"
	"github.com/majd/ipatool/pkg/http"
	"github.com/majd/ipatool/pkg/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"os"
)

var _ = Describe("AppStore (Search)", func() {
	var (
		ctrl       *gomock.Controller
		mockClient *http.MockClient[SearchResult]
		mockLogger *log.MockLogger
		as         AppStore
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockClient = http.NewMockClient[SearchResult](ctrl)
		mockLogger = log.NewMockLogger(ctrl)
		as = &appstore{
			searchClient: mockClient,
			ioReader:     os.Stdin,
			logger:       mockLogger,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("country code is invalid", func() {
		It("returns error", func() {
			err := as.Search("", "XYZ", "", 0)
			Expect(err).To(MatchError(ContainSubstring(ErrorInvalidCountryCode.Error())))
		})
	})

	When("device family is invalid", func() {
		It("returns error", func() {
			err := as.Search("", "US", "XYZ", 0)
			Expect(err).To(MatchError(ContainSubstring(ErrorInvalidDeviceFamily.Error())))
		})
	})

	When("request fails", func() {
		var testErr = errors.New("test")

		BeforeEach(func() {
			mockClient.EXPECT().
				Send(gomock.Any()).
				Return(http.Result[SearchResult]{}, testErr)
		})

		It("returns error", func() {
			err := as.Search("", "US", DeviceFamilyPhone, 0)
			Expect(err).To(MatchError(ContainSubstring(testErr.Error())))
			Expect(err).To(MatchError(ContainSubstring(ErrorRequest.Error())))
		})
	})

	When("request returns bad status code", func() {
		BeforeEach(func() {
			mockLogger.EXPECT().
				Debug().
				Return(nil)

			mockClient.EXPECT().
				Send(gomock.Any()).
				Return(http.Result[SearchResult]{
					StatusCode: 400,
				}, nil)
		})

		It("returns error", func() {
			err := as.Search("", "US", DeviceFamilyPad, 0)
			Expect(err).To(MatchError(ContainSubstring(ErrorRequest.Error())))
		})
	})

	When("request is successful", func() {
		BeforeEach(func() {
			mockLogger.EXPECT().
				Info().
				Return(nil)

			mockClient.EXPECT().
				Send(gomock.Any()).
				Return(http.Result[SearchResult]{
					StatusCode: 200,
					Data: SearchResult{
						Count:   0,
						Results: []App{},
					},
				}, nil)
		})

		It("returns nil", func() {
			err := as.Search("", "US", DeviceFamilyPhone, 0)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
