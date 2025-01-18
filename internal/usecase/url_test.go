package usecase

import (
	"errors"
	"fmt"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/domain/url"
	"leenwood/yandex-http/internal/domain/url/mocks"
	"leenwood/yandex-http/internal/usecase/dto"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNewUrlUseCase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Url UseCase Test Suite")
}

var _ = Describe("UrlUseCase", func() {
	var (
		ctrl       *gomock.Controller
		mockRepo   *mocks.MockRepositoryInterface
		cfg        config.Config
		urlUseCase UrlUseCaseInterface
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepositoryInterface(ctrl)
		cfg = config.Config{
			App: config.AppConfig{
				Hostname: "localhost",
				Port:     "8080",
			},
		}
		urlUseCase = &UrlUseCase{r: mockRepo, c: cfg}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("CreateShortUrl", func() {
		Context("when the URL already exists", func() {
			It("should return the existing short URL", func() {
				existingUrl := &url.Url{Id: "12345", OriginalUrl: "http://example.com", ClickCount: 10}

				mockRepo.EXPECT().FindByUrl("http://example.com").Return(existingUrl, nil)

				request := dto.CreateShortUrlUseCaseRequest{Url: "http://example.com"}
				response, err := urlUseCase.CreateShortUrl(request)

				Expect(err).NotTo(HaveOccurred())
				expectedUrl := fmt.Sprintf("%s:%s/%s", cfg.App.Hostname, cfg.App.Port, existingUrl.Id)
				Expect(response.Url).To(Equal(expectedUrl))
				Expect(response.ClickCount).To(Equal(existingUrl.ClickCount))
			})
		})

		Context("when the URL does not exist", func() {
			It("should create a new short URL", func() {
				mockRepo.EXPECT().FindByUrl("http://example.com").Return(nil, nil)
				newUrl := &url.Url{Id: "67890", OriginalUrl: "http://example.com", ClickCount: 0}
				mockRepo.EXPECT().Save("http://example.com", "").Return(newUrl, nil)

				request := dto.CreateShortUrlUseCaseRequest{Url: "http://example.com"}
				response, err := urlUseCase.CreateShortUrl(request)

				Expect(err).NotTo(HaveOccurred())
				expectedUrl := fmt.Sprintf("%s:%s/%s", cfg.App.Hostname, cfg.App.Port, newUrl.Id)
				Expect(response.Url).To(Equal(expectedUrl))
				Expect(response.ClickCount).To(Equal(newUrl.ClickCount))
			})
		})

		Context("when FindByUrl fails", func() {
			It("should return an error", func() {
				mockRepo.EXPECT().FindByUrl("http://example.com").Return(nil, errors.New("database error"))

				request := dto.CreateShortUrlUseCaseRequest{Url: "http://example.com"}
				response, err := urlUseCase.CreateShortUrl(request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
				Expect(response.Url).To(BeEmpty())
			})
		})

		Context("when Save fails", func() {
			It("should return an error", func() {
				mockRepo.EXPECT().FindByUrl("http://example.com").Return(nil, nil)
				mockRepo.EXPECT().Save("http://example.com", "").Return(nil, errors.New("save error"))

				request := dto.CreateShortUrlUseCaseRequest{Url: "http://example.com"}
				response, err := urlUseCase.CreateShortUrl(request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("save error"))
				Expect(response.Url).To(BeEmpty())
			})
		})
	})

	Describe("CreateShortUrlWithCustomId", func() {
		Context("when creating a new short URL with a custom ID", func() {
			It("should create the URL successfully", func() {
				mockRepo.EXPECT().FindByUrl("http://example.com").Return(nil, nil)
				newUrl := &url.Url{Id: "custom123", OriginalUrl: "http://example.com", ClickCount: 0}
				mockRepo.EXPECT().Save("http://example.com", "custom123").Return(newUrl, nil)

				request := dto.CreateShortUrlWithCustomIdRequest{Url: "http://example.com", Id: "custom123"}
				response, err := urlUseCase.CreateShortUrlWithCustomId(request)

				Expect(err).NotTo(HaveOccurred())
				expectedUrl := fmt.Sprintf("%s:%s/%s", cfg.App.Hostname, cfg.App.Port, newUrl.Id)
				Expect(response.Url).To(Equal(expectedUrl))
				Expect(response.ClickCount).To(Equal(newUrl.ClickCount))
			})
		})
	})
})
