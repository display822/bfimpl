/*
* Auth : acer
* Desc : 文件上传
* Time : 2020/9/6 14:51
 */

package controllers

import (
	"bfimpl/services"
	"bfimpl/services/log"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go"
)

type FileController struct {
	BaseController
}

// @Title 上传文件
// @Description 上传的文件
// @Param  file form-data binary true "文件"
// @Success 200 {string} "success"
// @Failure 500 server internal err
// @router /upload [post]
func (m *FileController) Upload() {
	_ = m.Ctx.Input.ParseFormOrMulitForm(100 << 20)
	file, header, err := m.GetFile("file")
	bucket := m.GetString("bucket", "default")
	if err != nil {
		log.GLogger.Error(err.Error())
		m.Error("need file")
	}
	defer file.Close()
	//创建path
	dir := "static/" + bucket
	f, e1 := os.Open(dir)
	if e1 != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			m.Error(fmt.Sprintf("create %s err.", dir))
		}
	}
	f.Close()
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + "_" + header.Filename
	tmpFile, err := os.Create(dir + "/" + fileName)
	if err != nil {
		log.GLogger.Error("upload err:%s", err.Error())
		m.Error(err.Error())
	}
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		m.Error(err.Error())
	}
	go upload(fileName, bucket)
	m.Correct("/" + bucket + "/" + fileName)
}
func upload(filename, bucket string) {
	contentType := "application/octet-stream"
	// 查询bucket 是否存在
	services.InitBucket(bucket)
	// 使用FPutObject上传一个文件。
	filePath := "static/" + bucket + "/" + filename
	n, err := services.MinIOClient().FPutObject(bucket, filename, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.GLogger.Error(err.Error())
	} else {
		log.GLogger.Info("Successfully uploaded %s of size %d\n", filename, n)
	}
	os.Remove(filePath)
}
