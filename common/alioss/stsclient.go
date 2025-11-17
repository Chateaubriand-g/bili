package alioss

import (
	"context"
	"fmt"
	"os"

	"github.com/Chateaubriand-g/bili/common/model"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type STSClient struct {
	STS        *sts.Client
	Ctx        context.Context
	BucketName string
	BucketUrl  string
	RoleArn    string
}

func getenv(key string) (string, error) {
	res := os.Getenv(key)
	if res == "" {
		return "", fmt.Errorf("%s required", key)
	}
	return res, nil
}

func initSTSClient() (*sts.Client, error) {
	accessKeyId, err := getenv("ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	accessKeySecret, err := getenv("ACCESS_KEY_SECRET")
	if err != nil {
		return nil, err
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}

	config.Endpoint = tea.String(endPoint)
	client, err := sts.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %w", err)
	}
	return client, nil
}

func NewSTSClient(ctx context.Context) (*STSClient, error) {
	stsclient, err := initSTSClient()
	if err != nil {
		return nil, fmt.Errorf("initSTSClient failed: %w", err)
	}
	rolearn, err := getenv("RAM_ROLE_ARN")
	if err != nil {
		return nil, fmt.Errorf("initSTSClient failed: %w", err)
	}

	return &STSClient{
		STS:        stsclient,
		Ctx:        ctx,
		BucketName: bucketName,
		BucketUrl:  bucketUrl,
		RoleArn:    rolearn,
	}, nil
}

func (stsc *STSClient) GetSTS(sessionname string) (*model.STSResponse, error) {
	request := &sts.AssumeRoleRequest{
		// 指定STS临时访问凭证过期时间为3600秒。
		DurationSeconds: tea.Int64(3600),
		// 从环境变量中获取步骤1.3生成的RAM角色的RamRoleArn。
		RoleArn: tea.String(stsc.RoleArn),
		// 指定自定义角色会话名称，这里使用和第一段代码一致的 examplename
		RoleSessionName: tea.String(sessionname),
	}
	response, err := stsc.STS.AssumeRoleWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		return &model.STSResponse{}, fmt.Errorf("Failed to assume role: %w", err)
	}

	credentials := response.Body.Credentials
	return &model.STSResponse{
		EndPoint:        endPoint,
		BucketName:      bucketName,
		AccessKeyId:     tea.StringValue(credentials.AccessKeyId),
		AccessKeySecret: tea.StringValue(credentials.AccessKeySecret),
		SecurityToken:   tea.StringValue(credentials.SecurityToken),
		Expiration:      tea.StringValue(credentials.SecurityToken),
	}, nil
}
