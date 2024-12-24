package gofile

const (
	UPLOAD_TYPE_LOCAL = 1
	UPLOAD_TYPE_OSS   = 2
	UPLOAD_TYPE_CHUNK = 3
)

type FileReceiveResult struct {
	OriginalFile      string //有后缀最终文件路径
	FileName          string //有后缀最终文件名
	FileMd5           string //最终文件Md5
	OriginalFileName  string //有后缀原始文件名
	OriginalShortName string //没有后缀的原始文件名
	ChunkCount        int64
}

// fmt.Sprintf("%s.part%d", fileName, i)
// mapreduce big file
// 大文件逻辑 for 把大文件并发分片处理，为了防止OOM超大文件边分片边处理的策略
type FileChunk struct {
	Data             []byte //分片数据
	Hash             string //分片Hash
	Index            int64  //分片顺序号
	OriginalFileName string //原文件名
	OriginalFileMd5  string //原文件Md5
	FileName         string //分片文件名称
}

type FileMergeReq struct {
	FileMd5     string `json:"fileMd5"`
	FileName    string `json:"fileName"`
	TotalChunks int64  `json:"totalChunks"`
}
