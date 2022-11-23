package server

import (
	"sync"
	"time"
)

type flipper struct {
	picker    PicPicker
	picLife   time.Duration
	log       Logger
	pic       []byte
	ondeck    []byte
	flipLck   sync.Mutex
	updateLck sync.Mutex
	flipAt    time.Time
}

func newFlipper(picker PicPicker, picLife time.Duration, log Logger) (*flipper, error) {
	f := flipper{picker: picker, picLife: picLife, log: log}

	var err error

	f.pic, err = picker.Next()
	if err != nil {
		return nil, err
	}

	f.ondeck, err = picker.Next()
	if err != nil {
		return nil, err
	}

	f.flipAt = time.Now().Add(f.picLife)

	return &f, nil
}

func (f *flipper) Next() ([]byte, error) {
	if time.Now().After(f.flipAt) {
		f.flip()
	}

	return f.pic, nil
}

func (f *flipper) flip() {
	f.flipLck.Lock()
	defer f.flipLck.Unlock()

	if time.Now().After(f.flipAt) {
		f.log.Debug("flip()")
		f.pic = f.ondeck
		f.flipAt = time.Now().Add(f.picLife)
		go f.update()
	}
}

func (f *flipper) update() {
	if !f.updateLck.TryLock() {
		return
	}
	defer f.updateLck.Unlock()

	buf, err := f.picker.Next()
	if err != nil {
		f.log.Error(err.Error())
		return
	}
	f.ondeck = buf
}
