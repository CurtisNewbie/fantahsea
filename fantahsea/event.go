package fantahsea

import "github.com/curtisnewbie/miso/miso"

const (
	AddDirGalleryImageEventBus = "fantahsea.dir.gallery.image.add"
	NotifyFileDeletedEventBus  = "fantahsea.notify.file.deleted"
)

func PrepareEventBus(rail miso.Rail) error {
	if e := miso.NewEventBus(AddDirGalleryImageEventBus); e != nil {
		return e
	}
	if e := miso.NewEventBus(NotifyFileDeletedEventBus); e != nil {
		return e
	}

	miso.SubEventBus(AddDirGalleryImageEventBus, 2, func(rail miso.Rail, evt CreateGalleryImgEvent) error {
		return OnCreateGalleryImgEvent(rail, evt)
	})

	miso.SubEventBus(NotifyFileDeletedEventBus, 2, func(rail miso.Rail, evt NotifyFileDeletedEvent) error {
		return OnNotifyFileDeletedEvent(rail, evt)
	})
	return nil
}
