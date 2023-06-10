package main

import (
	"log"
	"os"

	"github.com/curtisnewbie/fantahsea/web/controller"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/gocommon/server"
)

const (
	MNG_FILE_CODE = "manage-files"
	MNG_FILE_NAME = "Manage files"
)

func main() {
	ec := common.EmptyExecContext()
	server.DefaultBootstrapServer(os.Args, ec, func() error {

		if goauth.IsEnabled() {
			server.OnServerBootstrapped(func(c common.ExecContext) error {
				if e := goauth.AddResource(ec.Ctx, goauth.AddResourceReq{Code: MNG_FILE_CODE, Name: MNG_FILE_NAME}); e != nil {
					log.Fatalf("goauth.AddResource, %v", e)
				}
				return nil
			})

			goauth.ReportPathsOnBootstrapped()
		}

		// authenticated routes
		server.Get(server.OpenApiPath("/gallery/brief/owned"), controller.ListOwnedGalleryBriefsEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List owned gallery brief info"}))
		server.IPost(server.OpenApiPath("/gallery/new"), controller.CreateGalleryEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Create new gallery"}))
		server.IPost(server.OpenApiPath("/gallery/update"), controller.UpdateGalleryEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Update gallery"}))
		server.IPost(server.OpenApiPath("/gallery/delete"), controller.DeleteGalleryEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Delete gallery"}))
		server.IPost(server.OpenApiPath("/gallery/list"), controller.ListGalleriesEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List galleries"}))
		server.IPost(server.OpenApiPath("/gallery/access/grant"), controller.GrantGalleryAccessEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List granted access to the galleries"}))
		server.IPost(server.OpenApiPath("/gallery/images"), controller.ListImagesEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "List images of gallery"}))
		server.IPost(server.OpenApiPath("/gallery/image/transfer"), controller.TransferGalleryImageEndpoint,
			goauth.PathDocExtra(goauth.PathDoc{Type: goauth.PT_PROTECTED, Desc: "Host selected images on gallery"}))

		return nil
	})
}
