package main

import (
	"image"
	"image/color"
	"net/http"
	"sync"
	"unsafe"

	"github.com/puzpuzpuz/xsync/v2"
)

const MAX_LOCKS = 50

var (
	namedMutexPool = make([]sync.Mutex, MAX_LOCKS)
	imagesCache    = xsync.NewMapOf[image.Image]()
	transparentBox = boxImage(color.Transparent)
)

func boxImage(color color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color)
	return img
}

func imageFromURL(u string) (res image.Image) {
	res, ok := imagesCache.Load(u)
	if ok {
		return res
	}

	// this is so that we only try to load the same url once
	unlock := namedLock(u)
	defer unlock()

	// store result on cache (even if it's nil)
	defer func() {
		imagesCache.Store(u, res)
	}()

	// load url
	response, err := http.Get(u)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil
	}

	return img
}

func namedLock(name string) (unlock func()) {
	sptr := unsafe.StringData(name)
	idx := uint64(memhash(unsafe.Pointer(sptr), 0, uintptr(len(name)))) % MAX_LOCKS
	namedMutexPool[idx].Lock()
	return namedMutexPool[idx].Unlock
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr
