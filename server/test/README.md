# Tests

All tests in this folder follow a similar concept to Medium tests in [github.com/immich-app/immich](https://github.com/immich-app/immich/tree/main/server/test/medium). This allows us to test against a live database which helps catch bugs related to assumptions about the behavior of queries and the database.

## Running tests

To run the tests, you can use the following command from within the `server` folder:

```bash
mise test
```

## Test Structure

Test structure follows a similar format to vitest through the ginkgo framework. An example test suite that uses a real database is shown below:

```go
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
```
