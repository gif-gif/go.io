package adapters

import (
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/go-utils/gotime"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"
)

type FileAdapter struct {
	filepath string
	filename string
	fh       *os.File

	maxSize  int64
	count    int
	keepDays int
	ch       chan []byte
	mu       sync.Mutex
	timer    *gotime.Timer
}

func NewFileLog(opts ...FileOption) *golog.Logger {
	return golog.New(NewFileAdapter(opts...))
}

func NewFileAdapter(opt ...FileOption) *FileAdapter {
	opts := defaultFileOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	fa := &FileAdapter{
		filepath: opts.Filepath,
		maxSize:  opts.MaxSize,
		keepDays: opts.KeepDays,
		ch:       make(chan []byte, runtime.NumCPU()*2),
	}

	if l := len(fa.filepath); fa.filepath[l-1:] != "/" {
		fa.filepath += "/"
	}

	if _, err := os.Stat(fa.filepath); os.IsNotExist(err) {
		os.MkdirAll(fa.filepath, 0755)
	}

	var (
		nw  = time.Now()
		ymd = nw.Format("20060102")
	)

	files, _ := filepath.Glob(fa.filepath + ymd + "_*.log")
	fa.count = len(files)

	fa.timer = gotime.NewTimer(opts.ClearLogInterval) // 每天12点清理一次
	fa.timer.Start(func() {
		fa.CleanOldLogs()
	})

	goutils.AsyncFunc(func() {
		select {
		case <-gocontext.WithCancel().Done():
			fa.timer.Stop()
			fmt.Println("日志文件清理定时器已停止")
		}
	})

	return fa
}

func (fa *FileAdapter) Write(msg *golog.Message) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		c := <-fa.ch
		fa.writeHandle(c)
	}()

	fa.ch <- msg.JSON()
}

func (fa *FileAdapter) writeHandle(b []byte) {
	fa.mu.Lock()
	defer func() { fa.mu.Unlock() }()

	var (
		nw       = time.Now()
		ymd      = nw.Format("20060102")
		filename = ymd + ".log"
	)

	if filename != fa.filename {
		fa.filename = filename
		fa.closeFile()
	}

	if fa.fh == nil {
		if err := fa.openFile(); err != nil {
			return
		}
	}

	if fa.maxSize > 0 {
		if err := fa.cutFile(ymd); err != nil {
			return
		}
	}

	if fa.fh != nil {
		fa.fh.Write(b)
		fa.fh.Write([]byte("\n"))
	}
}

func (fa *FileAdapter) closeFile() {
	if fa.fh == nil {
		return
	}

	fa.fh.Close()
	fa.fh = nil
}

func (fa *FileAdapter) openFile() (err error) {
	fa.fh, err = os.OpenFile(fa.filepath+fa.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (fa *FileAdapter) cutFile(ymd string) (err error) {
	var info os.FileInfo

	info, err = os.Stat(fa.filepath + fa.filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if info.Size() < fa.maxSize {
		return
	}

	fa.count++

	filename := fmt.Sprintf("%s_%d.log", fa.filepath+ymd, fa.count)
	if err = os.Rename(fa.filepath+fa.filename, filename); err != nil {
		log.Println(err.Error())
		return
	}

	if fa.fh != nil {
		fa.fh.Close()
		fa.fh = nil
	}

	if err = fa.openFile(); err != nil {
		return
	}
	return
}

func (fa *FileAdapter) CleanOldLogs() error {
	// 获取当前时间
	now := time.Now()

	// 计算截止时间（7天前）
	cutoffTime := now.AddDate(0, 0, -fa.keepDays)

	// 读取目录中的所有文件
	files, err := ioutil.ReadDir(fa.filepath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	// 日志文件名正则表达式
	// 匹配格式: 20250525.log 或 20250528_1.log
	logPattern := regexp.MustCompile(`^(\d{8})(_\d+)?\.log$`)

	deletedCount := 0
	keptCount := 0

	for _, file := range files {
		// 跳过目录
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		// 检查是否是日志文件
		matches := logPattern.FindStringSubmatch(fileName)
		if matches == nil {
			continue
		}

		// 解析日期
		dateStr := matches[1]
		fileDate, err := time.Parse("20060102", dateStr)
		if err != nil {
			//log.Printf("解析日期失败 %s: %v", fileName, err)
			continue
		}

		// 获取文件的完整路径
		filePath := filepath.Join(fa.filepath, fileName)

		// 判断是否需要删除（文件日期早于截止日期）
		if fileDate.Before(cutoffTime) {
			// 删除文件
			if err := os.Remove(filePath); err != nil {
				//log.Printf("删除文件失败 %s: %v", filePath, err)
			} else {
				deletedCount++
				//log.Printf("已删除: %s (日期: %s)", fileName, fileDate.Format("2006-01-02"))
			}
		} else {
			keptCount++
			//log.Printf("保留: %s (日期: %s)", fileName, fileDate.Format("2006-01-02"))
		}
	}

	//log.Printf("清理完成: 删除 %d 个文件, 保留 %d 个文件", deletedCount, keptCount)
	return nil
}
