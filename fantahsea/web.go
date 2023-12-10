package fantahsea

import (
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/gin-gonic/gin"
)

const (
	ManageFileCode = "manage-files"
	ManageFileName = "Manage files"
)

func RegisterRoutes(rail miso.Rail) error {
	goauth.ReportResourcesOnBootstrapped(rail, []goauth.AddResourceReq{{Code: ManageFileCode, Name: ManageFileName}})
	goauth.ReportPathsOnBootstrapped(rail)

	miso.BaseRoute("/open/api").Group(
		miso.Get("/gallery/brief/owned",
			func(c *gin.Context, rail miso.Rail) (any, error) {
				user := common.GetUser(rail)
				return ListOwnedGalleryBriefs(rail, user)
			}).
			Extra(goauth.Protected("List owned gallery brief info", ManageFileCode)),

		miso.IPost("/gallery/new",
			func(c *gin.Context, rail miso.Rail, cmd CreateGalleryCmd) (any, error) {
				user := common.GetUser(rail)
				return CreateGallery(rail, cmd, user)
			}).
			Extra(goauth.Protected("Create new gallery", ManageFileCode)),

		miso.IPost("/gallery/update",
			func(c *gin.Context, rail miso.Rail, cmd UpdateGalleryCmd) (any, error) {
				user := common.GetUser(rail)
				e := UpdateGallery(rail, cmd, user)
				return nil, e
			}).
			Extra(goauth.Protected("Update gallery", ManageFileCode)),

		miso.IPost("/gallery/delete",
			func(c *gin.Context, rail miso.Rail, cmd DeleteGalleryCmd) (any, error) {
				user := common.GetUser(rail)
				e := DeleteGallery(rail, cmd, user)
				return nil, e
			}).
			Extra(goauth.Protected("Delete gallery", ManageFileCode)),

		miso.IPost("/gallery/list",
			func(c *gin.Context, rail miso.Rail, cmd ListGalleriesCmd) (any, error) {
				user := common.GetUser(rail)
				return ListGalleries(rail, cmd, user)
			}).
			Extra(goauth.Protected("List galleries", ManageFileCode)),

		miso.IPost("/gallery/access/grant",
			func(c *gin.Context, rail miso.Rail, cmd PermitGalleryAccessCmd) (any, error) {
				user := common.GetUser(rail)
				e := GrantGalleryAccessToUser(rail, cmd, user)
				return nil, e
			}).
			Extra(goauth.Protected("Grant access to the galleries", ManageFileCode)),

		miso.IPost("/gallery/access/list",
			func(c *gin.Context, rail miso.Rail, cmd ListGrantedGalleryAccessCmd) (any, error) {
				user := common.GetUser(rail)
				return ListedGrantedGalleryAccess(rail, miso.GetMySQL(), cmd, user)
			}).
			Extra(goauth.Protected("List granted access to the galleries", ManageFileCode)),

		miso.IPost("/gallery/images",
			func(c *gin.Context, rail miso.Rail, cmd ListGalleryImagesCmd) (any, error) {
				user := common.GetUser(rail)
				return ListGalleryImages(rail, cmd, user)
			}).
			Extra(goauth.Protected("List images of gallery", ManageFileCode)),

		miso.IPost("/gallery/image/transfer",
			func(c *gin.Context, rail miso.Rail, cmd TransferGalleryImageReq) (any, error) {
				user := common.GetUser(rail)
				return BatchTransferAsync(rail, cmd, user)
			}).
			Extra(goauth.Protected("Host selected images on gallery", ManageFileCode)),
	)
	return nil
}
