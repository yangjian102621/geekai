package storage

import (
	"bytes"
	"crypto/sha1"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Recorder interface {
	// 新建或更新文件分片上传的进度
	Set(key string, data []byte) error

	// 获取文件分片上传的进度信息
	Get(key string) ([]byte, error)

	// 删除文件分片上传的进度文件
	Delete(key string) error

	// 根据给定的文件信息生成持久化纪录的 key
	GenerateRecorderKey(keyInfos []string, sourceFileInfo os.FileInfo) string
}

type FileRecorder struct {
	directoryPath string
}

func NewFileRecorder(directoryPath string) (fileRecorder *FileRecorder, err error) {
	err = os.MkdirAll(directoryPath, 0700)
	if err != nil {
		return
	}
	fileRecorder = &FileRecorder{directoryPath: directoryPath}
	return
}

func (fileRecorder *FileRecorder) Set(key string, data []byte) error {
	path := filepath.Join(fileRecorder.directoryPath, key)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(data))
	return err
}

func (fileRecorder *FileRecorder) Get(key string) ([]byte, error) {
	path := filepath.Join(fileRecorder.directoryPath, key)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if isRecorderFileOutOfDate(fileInfo) {
		err = fileRecorder.Delete(key)
		return nil, err
	}

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	return buffer.Bytes(), err
}

func (fileRecorder *FileRecorder) Delete(key string) error {
	path := filepath.Join(fileRecorder.directoryPath, key)
	return os.Remove(path)
}

func (fileRecorder *FileRecorder) GenerateRecorderKey(keyInfos []string, sourceFileInfo os.FileInfo) string {
	const delimiter = "*:|>?^ \b"
	buffer := new(bytes.Buffer)
	for _, keyInfo := range keyInfos {
		buffer.WriteString(keyInfo)
		buffer.WriteString(delimiter)
	}
	buffer.WriteString(sourceFileInfo.ModTime().String())
	return hashRecorderKey(buffer.Bytes())
}

func isRecorderFileOutOfDate(fileInfo os.FileInfo) bool {
	return fileInfo.ModTime().Add(24 * 5 * time.Hour).Before(time.Now())
}

func hashRecorderKey(base []byte) string {
	sha1edBase := sha1.Sum(base)

	var stringBuilder strings.Builder
	for i := 0; i < len(sha1edBase); i++ {
		b := (int64(sha1edBase[i]) & 0xff) + 0x100
		stringBuilder.WriteString(strconv.FormatInt(b, 16)[1:])
	}
	return stringBuilder.String()
}
