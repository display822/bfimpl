/*
* Auth : acer
* Desc : minio
* Time : 2020/9/6 14:44
 */

package services

import (
	"bfimpl/services/log"

	"github.com/astaxie/beego"
	"github.com/minio/minio-go"
)

var client *minio.Client

func SetMinIOClient(c *minio.Client) {
	client = c
}
func MinIOClient() *minio.Client {
	if client == nil {
		InitIO()
	}
	return client
}
func InitIO() {
	endpoint := beego.AppConfig.String("endpoint")
	accessKeyID := beego.AppConfig.String("accessKeyID")
	secretAccessKey := beego.AppConfig.String("secretAccessKey")
	// 初使化 minio client对象。
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		log.GLogger.Critical("init minio client err:%s\n", err.Error())
	}
	SetMinIOClient(minioClient)
	InitBucket("default")
	// minio Client初使化成功
	log.GLogger.Info("minio client init success.")
}
func InitBucket(bucket string) {
	// 创建一个default bucket
	location := "us-east-1"
	// 检查存储桶是否已经存在。
	exists, err := MinIOClient().BucketExists(bucket)
	if err == nil && exists {
		log.GLogger.Info("We already own %s\n", bucket)
		return
	}
	err = MinIOClient().MakeBucket(bucket, location)
	if err != nil {
		log.GLogger.Error("create bucket err:%s", err.Error())
	}
}
