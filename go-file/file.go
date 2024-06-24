package gofile

import (
	"fmt"
	"github.com/gif-gif/go.io/go-utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	EL = "\n"
)

// 文件名
func FILE() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// 行号
func LINE() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// 目录名称
func DIR() string {
	_, file, _, _ := runtime.Caller(1)
	return path.Dir(file) + "/"
}

// 写文件，支持路径创建
func WriteToFile(filename string, b []byte) error {
	dirname := path.Dir(filename)
	if _, err := os.Stat(dirname); err != nil {
		os.MkdirAll(dirname, 0755)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

// 追踪信息
func Trace(skip int) (arr []string) {
	arr = []string{}
	if skip == 0 {
		skip = 1
	}
	for i := skip; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" {
			continue
		}
		if strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "runtime/") ||
			strings.Contains(file, "src/testing") ||
			strings.Contains(file, "pkg/mod/") ||
			strings.Contains(file, "vendor/") {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s %dL", prettyFile(file), line))
	}
	return
}

func prettyFile(file string) string {
	index := strings.LastIndex(file, "/")
	if index < 0 {
		return file
	}

	index2 := strings.LastIndex(file[:index], "/")
	if index2 < 0 {
		return file[index+1:]
	}

	return file[index2+1:]
}

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 文件或目录存在
	}
	if os.IsNotExist(err) {
		return false, nil // 文件或目录不存在
	}
	return false, err // 发生了其他错误，无法确定
}

func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func GetFileHeaderMd5Name(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	name := goutils.MD5(body)

	return name, nil //+ filepath.Ext(fileHeader.Filename)
}

// 复制文件
func CopyFile(src, dst string) error {
	s, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, s, 0o600)
	if err != nil {
		return err
	}
	return nil
}

// 获取当前目录下所有文件
func GetFileList(path string) []string {
	var fileList []string
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if !f.IsDir() {
			fileList = append(fileList, f.Name())
		}
	}
	return fileList
}

func RemoveFile(file string) error {
	return os.Remove(file)
}
