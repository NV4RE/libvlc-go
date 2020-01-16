package vlc

// #cgo LDFLAGS: -lvlc
// #include <vlc/vlc.h>
import "C"
import "errors"

// MediaList represents a collection of media files.
type MediaList struct {
	list *C.libvlc_media_list_t
}

// NewMediaList creates an empty media list.
func NewMediaList() (*MediaList, error) {
	if inst == nil {
		return nil, errors.New("module must be initialized first")
	}

	var list *C.libvlc_media_list_t
	if list = C.libvlc_media_list_new(inst.handle); list == nil {
		return nil, getError()
	}

	return &MediaList{list: list}, nil
}

// Release destroys the media list instance.
func (ml *MediaList) Release() error {
	if ml.list == nil {
		return nil
	}

	C.libvlc_media_list_release(ml.list)
	ml.list = nil

	return getError()
}

// AddMedia adds a Media instance to the media list.
func (ml *MediaList) AddMedia(m *Media) error {
	if ml.list == nil {
		return errors.New("media list must be initialized first")
	}
	if m.media == nil {
		return errors.New("media must be initialized first")
	}

	C.libvlc_media_list_add_media(ml.list, m.media)
	return getError()
}

// AddMediaFromPath loads a media file from path and adds it
// to the the media list.
func (ml *MediaList) AddMediaFromPath(path string) error {
	media, err := NewMediaFromPath(path)
	if err != nil {
		return err
	}

	return ml.AddMedia(media)
}

// AddMediaFromURL loads a media file from url and adds it
// to the the media list.
func (ml *MediaList) AddMediaFromURL(url string) error {
	media, err := NewMediaFromURL(url)
	if err != nil {
		return err
	}

	return ml.AddMedia(media)
}

// Lock makes the caller the current owner of the media list.
func (ml *MediaList) Lock() error {
	if ml.list == nil {
		return errors.New("media list must be initialized first")
	}

	C.libvlc_media_list_lock(ml.list)
	return getError()
}

// Unlock releases ownership of the media list.
func (ml *MediaList) Unlock() error {
	if ml.list == nil {
		return errors.New("media list must be initialized first")
	}

	C.libvlc_media_list_unlock(ml.list)
	return getError()
}

// EventManager returns the event manager responsible for the media list.
func (ml *MediaList) EventManager() (*EventManager, error) {
	if ml.list == nil {
		return nil, errors.New("media list must be initialized first")
	}

	manager := C.libvlc_media_list_event_manager(ml.list)
	if manager == nil {
		return nil, errors.New("could not retrieve media list event manager")
	}

	return newEventManager(manager), nil
}


// Size get list size
// List must be locked
func (ml *MediaList) Size() (int, error) {
	cResult := C.libvlc_media_list_count(ml.list)

	return int(cResult), getError()
}

// RemoveAtIndex remove item from list
// List must be locked
func (ml *MediaList) RemoveAtIndex(index int) error {
	cIndex := C.int(index)
	cResult := C.libvlc_media_list_remove_index(ml.list ,cIndex)

	if (int(cResult) == 0) {
		return nil
	}

	return errors.New("could not remove from list")
}

// ClearList cleanup list
// List must be locked
func (ml *MediaList) ClearList(keep uint) error {
	size, err := ml.Size()
	if (err != nil) {
		return err
	}

	if (size <= int(keep)) {
		return nil
	}

	for i := size - 1; i >= int(keep); i-- {
		if err := ml.RemoveAtIndex(i); err != nil {
			return err
		}
	}

	return nil
}
