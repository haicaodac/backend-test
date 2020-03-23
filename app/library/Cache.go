package library

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

const (
	cacheDir = "tmp" // Cache directory
)

// CacheSet ...
func CacheSet(key string, data interface{}) error {

	key = regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllLiteralString(key, "")

	clean(key)

	file := "filecache." + key
	fpath := filepath.Join(cacheDir, file)

	serialized, err := serialize(data)
	if err != nil {
		return err
	}

	var fmutex sync.RWMutex

	fmutex.Lock()
	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err = fp.Write(serialized); err != nil {
		return err
	}
	defer fmutex.Unlock()
	return nil
}

// CacheGet ...
func CacheGet(key string, dst interface{}) error {

	key = regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllLiteralString(key, "")

	pattern := filepath.Join(cacheDir, "filecache."+key)
	files, err := filepath.Glob(pattern)
	if len(files) < 1 || err != nil {
		return errors.New("filecache: no cache file found")
	}

	if _, err = os.Stat(files[0]); err != nil {
		return err
	}

	fp, err := os.OpenFile(files[0], os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer fp.Close()

	var serialized []byte
	buf := make([]byte, 1024)
	for {
		var n int
		n, err = fp.Read(buf)
		serialized = append(serialized, buf[0:n]...)
		if err != nil || err == io.EOF {
			break
		}
	}

	if err = deserialize(serialized, dst); err != nil {
		return err
	}
	return nil
}

// CacheCleanAll ...
func CacheCleanAll() {
	files, _ := ioutil.ReadDir(cacheDir)
	for _, file := range files {
		pattern := filepath.Join(cacheDir, file.Name())
		os.Remove(pattern)
	}
}

func clean(key string) {
	pattern := filepath.Join(cacheDir, "filecache."+key)
	files, _ := filepath.Glob(pattern)
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}
}

// serialize encodes a value using binary.
func serialize(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// deserialize decodes a value using binary.
func deserialize(src []byte, dst interface{}) error {
	buf := bytes.NewReader(src)
	if err := gob.NewDecoder(buf).Decode(dst); err != nil {
		log.Println(dst)
		return err
	}
	return nil
}
