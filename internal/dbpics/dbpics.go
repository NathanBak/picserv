package dbpics

import (
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Config struct {
	Token string `envvar:"TOKEN,required"`
	Seed  int64  `envvar:"SEED"`
}

type DbPics struct {
	cfg     Config
	fs      files.Client
	entries []*files.FileMetadata
	current int
}

func New(cfg Config) (*DbPics, error) {

	fs := files.New(dropbox.Config{
		Token:    cfg.Token,
		LogLevel: dropbox.LogInfo,
	})

	res, err := fs.ListFolder(&files.ListFolderArg{})
	if err != nil {
		return nil, err
	}

	db := DbPics{
		cfg: cfg,
		fs:  fs,
	}

	for _, e := range res.Entries {
		f, ok := e.(*files.FileMetadata)
		if !ok {
			continue
		}
		db.entries = append(db.entries, f)
	}

	if len(db.entries) == 0 {
		return nil, errors.New("no files found")
	}

	db.shuffleEntries()

	return &db, nil
}

func (db *DbPics) Next() ([]byte, error) {

	_, r, err := db.fs.Download(files.NewDownloadArg(db.entries[db.current].PathLower))
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	db.current++
	if db.current >= len(db.entries) {
		db.shuffleEntries()
		db.current = 0
	}

	return buf, nil
}

// shuffleEntries will randomize the order of the items in the entries slice.  It is the
// responsiblity of the caller to ensure that access to the slice is appropriately synchronized
// prior to calling this function.  If the Seed is 0 then "Now" is used as the seed.  If the Seed
// is negative, no shuffling will occurr.  Otherwise, the specified seed will be used when the
// shuffling is performed.
func (db *DbPics) shuffleEntries() {
	seed := db.cfg.Seed

	if seed < 0 {
		return
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)

	rand.Shuffle(len(db.entries),
		func(i, j int) { db.entries[i], db.entries[j] = db.entries[j], db.entries[i] })
}
