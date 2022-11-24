package dbpics

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Config struct {
	Seed         int64  `envvar:"SEED"`
	AppKey       string `envvar:"APP_KEY"`
	AppSecret    string `envvar:"APP_SECRET"`
	RefreshToken string `envvar:"USER_REFRESH_TOKEN"`
}

type DbPics struct {
	cfg     Config
	fs      files.Client
	entries []*files.FileMetadata
	current int
	lck     sync.Mutex
}

func New(cfg Config) (*DbPics, error) {

	db := DbPics{
		cfg: cfg,
	}

	fs, err := newFileClient(cfg.AppKey, cfg.AppSecret, cfg.RefreshToken)
	if err != nil {
		return nil, err
	}

	db.fs = fs

	db.entries, err = getAllFilesInDirectory(fs)
	if err != nil {
		return nil, err
	}

	if len(db.entries) == 0 {
		return nil, errors.New("no files found")
	}

	db.shuffleEntries()

	return &db, nil
}

// Next returns a slice of bytes representing the "next" image.
func (db *DbPics) Next() ([]byte, error) {

	db.lck.Lock()
	currentEntry := db.entries[db.current]
	db.lck.Unlock()

	_, r, err := db.fs.Download(files.NewDownloadArg(currentEntry.PathLower))
	if err != nil {
		fmt.Printf("error in dbpics.Next(): %v\nwill try refreshing access token\n", err)
		fs, fcerr := newFileClient(db.cfg.AppKey, db.cfg.AppSecret, db.cfg.RefreshToken)
		if fcerr != nil {
			return nil, fcerr
		}
		fmt.Println("access token refreshed")
		db.fs = fs
		_, r, err = db.fs.Download(files.NewDownloadArg(db.entries[db.current].PathLower))
		if err != nil {
			return nil, err
		}
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	db.lck.Lock()
	defer db.lck.Unlock()
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

// newFileClient uses the refreshToken to get a new accessToken and then uses the accessToken to
// create a new Dropbox files client.
func newFileClient(appKey, appSecret, refreshToken string) (files.Client, error) {
	oat, err := RefreshOfflineAccessToken(appKey, appSecret, refreshToken)
	if err != nil {
		return nil, err
	}

	return files.New(dropbox.Config{
		Token:    oat.AccessToken,
		LogLevel: dropbox.LogInfo,
	}), nil
}

// getAllFilesInDirectory returns all the FileMetadata associated with files in the path.
// Depending on the number of files in the directory, it may make several calls to DropBox.
func getAllFilesInDirectory(fs files.Client) ([]*files.FileMetadata, error) {

	const limit = 250
	var entries []*files.FileMetadata

	res, err := fs.ListFolder(&files.ListFolderArg{Limit: limit})
	if err != nil {
		return nil, err
	}

	addEntries := func(result *files.ListFolderResult) {
		for _, e := range result.Entries {
			f, ok := e.(*files.FileMetadata)
			if !ok {
				continue
			}
			entries = append(entries, f)
		}
	}

	addEntries(res)

	for res.HasMore {
		res, err = fs.ListFolderContinue(&files.ListFolderContinueArg{Cursor: res.Cursor})
		if err != nil {
			return nil, err
		}
		addEntries(res)
	}

	return entries, nil
}
