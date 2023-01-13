package servers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"

	"GameJamPlatform/internal/log"
)

func (s *server) uploadImage(r *http.Request, formFileName string) (string, error) {
	file, handler, err := r.FormFile(formFileName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}

	defer file.Close()
	log.Info(fmt.Sprintf("Uploaded File: %+v", handler.Filename))
	log.Info(fmt.Sprintf("File Size: %+v", handler.Size))
	log.Info(fmt.Sprintf("MIME Header: %+v", handler.Header))

	err = os.MkdirAll(s.cfg.ImageDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	uid := generateUID()
	imageExt := path.Ext(handler.Filename)
	imagePath := fmt.Sprintf("%s/%s%s", s.cfg.ImageDir, uid, imageExt)
	dst, err := os.Create(imagePath)
	defer dst.Close()
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(dst, file); err != nil {
		err := os.Remove(imagePath)
		if err != nil {
			return "", err
		}
		return "", err
	}

	return uid + imageExt, nil
}

func generateUID() string {
	return uuid.New().String()
}
