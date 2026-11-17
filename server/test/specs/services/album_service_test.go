package specs

import (
	"fhs/app/internal/dto"
	"fhs/app/internal/service"
	"fhs/app/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AlbumService", Label("medium", "services"), func() {
	var sut service.AlbumService

	BeforeEach(func() {
		repos, stop := test.Configure()
		DeferCleanup(stop)

		sut = service.NewAlbumService(repos, &test.Config)
	})

	Describe("CreateAlbum", func() {
		It("persists the album name", func(ctx SpecContext) {
			created, err := sut.CreateAlbum(ctx, &dto.CreateAlbumDto{Name: "Test Album"})
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Name).To(Equal("Test Album"))

			saved, err := sut.GetAlbum(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(saved.Name).To(Equal("Test Album"))
		})
	})
})
