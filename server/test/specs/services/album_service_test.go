package specs

import (
	"context"
	"fhs/app/internal/dto"
	"fhs/app/internal/service"
	"fhs/app/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("album_service", Label("services"), func() {
	var svc service.AlbumService
	var ctx = context.Background()

	BeforeEach(func() {
		repos, stop := test.Configure()
		DeferCleanup(stop)

		svc = service.NewAlbumService(repos, &test.Config)
	})

	It("should create an album", func() {
		album, err := svc.CreateAlbum(ctx, &dto.CreateAlbumDto{Name: "Test Album"})
		Expect(err).To(BeNil())
		Expect(album.Name).To(Equal("Test Album"))

		album2, err := svc.GetAlbum(ctx, album.ID)
		Expect(err).To(BeNil())
		Expect(album2.Name).To(Equal("Test Album"))
	})
})
