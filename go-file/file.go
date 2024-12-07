package gofile

import (
	"bufio"
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/go-utils"
	"io"
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

// GetFileInfo 获取文件信息
func GetFileInfo(filename string) (size int64, mode os.FileMode, err error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, 0, err
	}
	return info.Size(), info.Mode(), nil
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

	body, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	name := goutils.Md5(body)

	return name, nil //+ filepath.Ext(fileHeader.Filename)
}

// ReadEntireFile 读取整个文件到字节切片一样
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		golog.WithTag("file").Error(err)
		return nil, err
	}

	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ReadLines(filePath string, lineFunc func(line string) error) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	// 创建一个新的扫描器
	scanner := bufio.NewScanner(file)
	// 按行扫描文件
	for scanner.Scan() {
		line := scanner.Text()
		e := lineFunc(line)
		if e != nil {
			return e
		}
	}
	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading file: %v", err)
	}

	return nil
}

func GetFileContentString(filePath string) (string, error) {
	body, err := GetFileContent(filePath)
	if err != nil {
		return "", err
	}

	return string(body), nil
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
	files, _ := os.ReadDir(path)
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
