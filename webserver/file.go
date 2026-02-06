package webserver

import (
	"fmt"
	"io"
	"os"

	"path/filepath"
	"time"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

// uploadfile 通用文件上传接口。返回示例：
//
//	{"status":0, "msg": "", "data": {"value": "xxxx"}}
func uploadfile(ctx httpsvr.Context) {
	// 获取上传的文件
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError(err.Error()).Bytes())
		return
	}
	// fmt.Printf("----------upload-----filename(%s)----\n", header.Filename)
	defer file.Close()
	unixtime := time.Now().Unix()
	saveFilename := fmt.Sprintf("%s-%d.xlsx", time.Now().Format(time.DateOnly), unixtime)

	// 创建上传目录（如果不存在）
	uploadDir := "runtime/upload"
	if !IsPathExists(uploadDir) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			ctx.Writer.Write(response.NewApiDataServerError("创建上传目录失败：" + err.Error()).Bytes())
			return
		}
	}

	// 构建完整文件路径并保存文件
	filepath := filepath.Join(uploadDir, saveFilename)
	dst, err := os.Create(filepath)
	if err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError("创建目标文件失败: " + err.Error()).Bytes())
		return
	}
	defer dst.Close()

	// 将上传文件内容复制到目标文件
	if _, err := io.Copy(dst, file); err != nil {
		ctx.Writer.Write(response.NewApiDataQueryArgsError("保存文件失败: " + err.Error()).Bytes())
		return
	}
	ctx.Writer.Write(response.NewApiData(response.JsonObject{"value": filepath}, "success", 0).Bytes())
}

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// fmt.Println(stat.IsDir())
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
