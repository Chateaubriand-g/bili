package alioss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     = "cn-guangzhou"
	bucketName = "bili-gz-a1"
	bucketUrl  = "bili-gz-a1.oss-cn-guangzhou.aliyuncs.com"
	endPoint   = "sts.cn-guangzhou.aliyuncs.com"
)

type OssClient struct {
	Oss        *oss.Client
	Ctx        context.Context
	BucketName string
	BucketUrl  string
}

type PresignRes struct {
	URL           string
	SignedHeaders map[string]string
}

func InitOssClient(cfg *config.Config) *oss.Client {
	osscfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion("cn-guangzhou").
		WithConnectTimeout(10 * time.Second).
		WithReadWriteTimeout(30 * time.Second)

	client := oss.NewClient(osscfg)
	return client
}

func NewOssClient(ctx context.Context, cfg *config.Config) *OssClient {
	return &OssClient{
		Oss:        InitOssClient(cfg),
		Ctx:        ctx,
		BucketName: bucketName,
		BucketUrl:  bucketUrl,
	}
}

func (ossc *OssClient) PutObject(objectKey, contextType string, reader io.Reader) error {
	ossrequest := &oss.PutObjectRequest{
		Bucket:      oss.Ptr(bucketName),
		Key:         oss.Ptr(objectKey),
		Body:        reader,
		ContentType: oss.Ptr(contextType),
	}

	_, err := ossc.Oss.PutObject(ossc.Ctx, ossrequest)
	if err != nil {
		return fmt.Errorf("failed to upload to OSS: %w", err)
	}

	return nil
}

func (ossc *OssClient) PutSignURL(objectName string) (*PresignRes, error) {
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(ossc.BucketName),
		Key:    oss.Ptr(objectName),
	}

	result, err := ossc.Oss.Presign(
		context.TODO(), request,
		oss.PresignExpiration(time.Now().Add(60*time.Minute)),
	)
	if err != nil {
		return &PresignRes{}, fmt.Errorf("oss presign error: %w", err)
	}
	return &PresignRes{URL: result.URL, SignedHeaders: result.SignedHeaders}, nil
}

func (ossc *OssClient) GetSignURL(objectName string) (*PresignRes, error) {
	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(ossc.BucketName),
		Key:    oss.Ptr(objectName),
	}

	result, err := ossc.Oss.Presign(
		ossc.Ctx, request,
		oss.PresignExpires(30*time.Minute),
	)
	if err != nil {
		return &PresignRes{}, fmt.Errorf("oss presign error: %w", err)
	}
	return &PresignRes{URL: result.URL, SignedHeaders: result.SignedHeaders}, nil
}

func (ossc *OssClient) MutipartUpload(objectName string, partSize int64, data []byte) {
	request := &oss.InitiateMultipartUploadRequest{
		Bucket: oss.Ptr(ossc.BucketName),
		Key:    oss.Ptr(objectName),
	}

	initResult, err := ossc.Oss.InitiateMultipartUpload(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed InitiateMultipartUpload %v", err)
	}

	fileSize := int64(len(data))
	partsNum := int(fileSize/partSize + 1)
	// 遍历每个分片，生成签名 URL 并上传分片
	for i := 0; i < partsNum; i++ {
		start := int64(i) * partSize
		end := start + partSize
		if end > fileSize {
			end = fileSize
		}
		signedResult, err := ossc.Oss.Presign(context.TODO(), &oss.UploadPartRequest{
			Bucket:     oss.Ptr(bucketName),
			Key:        oss.Ptr(objectName),
			PartNumber: int32(i + 1),
			Body:       bytes.NewReader(data[start:end]),
			UploadId:   initResult.UploadId,
		}, oss.PresignExpiration(time.Now().Add(1*time.Hour))) // 生成签名 URL，有效期为1小时
		if err != nil {
			log.Fatalf("failed to generate presigned URL %v", err)
		}
		fmt.Printf("signed url:%#v\n", signedResult.URL) // 打印生成的签名URL

		// 创建HTTP请求并上传分片
		req, err := http.NewRequest(signedResult.Method, signedResult.URL, bytes.NewReader(data[start:end]))
		if err != nil {
			log.Fatalf("failed to create HTTP request %v", err)
		}

		c := &http.Client{} // 创建HTTP客户端
		_, err = c.Do(req)
		if err != nil {
			log.Fatalf("failed to upload part by signed URL %v", err)
		}
	}

	// 列举已上传的分片
	partsResult, err := ossc.Oss.ListParts(context.TODO(), &oss.ListPartsRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: initResult.UploadId,
	})
	if err != nil {
		log.Fatalf("failed to list parts %v", err)
	}

	// 收集已上传的分片信息
	var parts []oss.UploadPart
	for _, p := range partsResult.Parts {
		parts = append(parts, oss.UploadPart{PartNumber: p.PartNumber, ETag: p.ETag})
	}

	// 完成分片上传
	result, err := ossc.Oss.CompleteMultipartUpload(context.TODO(), &oss.CompleteMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: initResult.UploadId,
		CompleteMultipartUpload: &oss.CompleteMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		log.Fatalf("failed to complete multipart upload %v", err)
	}

	// 打印完成分片上传的结果
	log.Printf("complete multipart upload result:%#v\n", result)
}
