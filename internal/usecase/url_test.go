package usecase

import (
	"errors"
	"fmt"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/domain/url"
	"leenwood/yandex-http/internal/domain/url/mocks"
	"leenwood/yandex-http/internal/usecase/dto"
	"testing"
	"time"

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

	Describe("GetUrlList", func() {
		Context("when URLs exist in the repository", func() {
			It("should return the paginated list of URLs", func() {
				pagination := dto.PaginationRequest{Page: 1, Limit: 2}
				createdDate := time.Now()
				mockUrls := []*url.Url{
					{
						Id:          "id1",
						OriginalUrl: "http://example1.com",
						ClickCount:  5,
						CreatedDate: createdDate,
					},
					{
						Id:          "id2",
						OriginalUrl: "http://example2.com",
						ClickCount:  10,
						CreatedDate: createdDate,
					},
				}

				mockRepo.EXPECT().FindAll(pagination.Page, pagination.Limit).Return(mockUrls, nil)

				response, err := urlUseCase.GetUrlList(pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(response).To(HaveLen(2))
				Expect(response[0].OriginalUrl).To(Equal("http://example1.com"))
				Expect(response[0].ShortUrl).To(Equal(fmt.Sprintf("%s:%s/%s", cfg.App.Hostname, cfg.App.Port, "id1")))
				Expect(response[0].CountClick).To(Equal(*getLink[uint64](5)))
				Expect(response[0].CreatedDate).To(Equal(createdDate))
				Expect(response[1].OriginalUrl).To(Equal("http://example2.com"))
				Expect(response[1].ShortUrl).To(Equal(fmt.Sprintf("%s:%s/%s", cfg.App.Hostname, cfg.App.Port, "id2")))
				Expect(response[1].CountClick).To(Equal(*getLink[uint64](10)))
				Expect(response[1].CreatedDate).To(Equal(createdDate))
			})
		})

		Context("when the repository returns an error", func() {
			It("should return an error", func() {
				pagination := dto.PaginationRequest{Page: 1, Limit: 2}

				mockRepo.EXPECT().FindAll(pagination.Page, pagination.Limit).Return(nil, errors.New("repository error"))

				response, err := urlUseCase.GetUrlList(pagination)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("repository error"))
				Expect(response).To(BeNil())
			})
		})
	})
})

func getLink[T any](variable T) *T {
	return &variable
}
