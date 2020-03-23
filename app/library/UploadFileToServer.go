package library

import (
	"io/ioutil"
	"mime/multipart"
	"path/filepath"

	"github.com/gosimple/slug"
)

// UploadFileToServer ...
func UploadFileToServer(file multipart.File, handle *multipart.FileHeader) (string, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	random := RandomString(10)
	filename := handle.Filename
	extension := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(extension)]
	fileHandle := slug.Make(random+"-"+name) + extension

	err = ioutil.WriteFile("public/uploads/files/"+fileHandle, data, 0666)
	if err != nil {
		return "", err
	}
	url := fileHandle
	return url, nil
}
