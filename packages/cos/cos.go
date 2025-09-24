package cos

import (
	"context"
	"craapi/packages/log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var client *cos.Client

func COS_INIT() {
	u, err := url.Parse(viper.GetString("tencentCOS.CosDomain"))
	if err != nil {
		log.LOGE("Error when Parse COS DOMAIN ", "error: ", err)
		os.Exit(1)
	}
	b := &cos.BaseURL{BucketURL: u}
	// 2.临时密钥
	client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 如果使用临时密钥需要填入，临时密钥生成和使用指引参见 https://cloud.tencent.com/document/product/436/14048
			SecretID:  viper.GetString("tencentCOS.SecretId"),
			SecretKey: viper.GetString("tencentCOS.SecretKey"),
		},
	})
	if client == nil {
		log.LOGE("COS init error!")
		os.Exit(1)
	}
}

func GetCOSDownloadURL(file string, accessToken string) (string, error) {
	ctx := context.Background()
	var presignedURL *url.URL
	var err error
	if viper.GetBool("tencentCOS.EnhanceSecurity") {
		opt := &cos.PresignedURLOptions{
			// http 请求参数，传入的请求参数需与实际请求相同，能够防止用户篡改此 HTTP 请求的参数
			Query: &url.Values{},
			// http 请求头部，传入的请求头部需包含在实际请求中，能够防止用户篡改签入此处的 HTTP 请求头部
			Header: &http.Header{},
		}
		opt.Header.Add("x-craapi-Access-token", accessToken)
		presignedURL, err = client.Object.GetPresignedURL(ctx, http.MethodGet, file, viper.GetString("tencentCOS.SecretId"), viper.GetString("tencentCOS.SecretKey"), time.Hour, opt, true)
	} else {
		presignedURL, err = client.Object.GetPresignedURL(ctx, http.MethodGet, file, viper.GetString("tencentCOS.SecretId"), viper.GetString("tencentCOS.SecretKey"), time.Hour, nil)
	}
	if err != nil {
		log.LOGE("COS get file: ", file, " ERROR: ", err)
		return "", err
	}
	return presignedURL.String(), nil
}

func GetCOSUploadURL(file string, accessToken string) (string, error) {
	ctx := context.Background()
	var presignedURL *url.URL
	var err error
	if viper.GetBool("tencentCOS.EnhanceSecurity") {
		opt := &cos.PresignedURLOptions{
			// http 请求参数，传入的请求参数需与实际请求相同，能够防止用户篡改此 HTTP 请求的参数
			Query: &url.Values{},
			// http 请求头部，传入的请求头部需包含在实际请求中，能够防止用户篡改签入此处的 HTTP 请求头部
			Header: &http.Header{},
		}
		opt.Header.Add("x-craapi-Access-token", accessToken)
		presignedURL, err = client.Object.GetPresignedURL(ctx, http.MethodPut, file, viper.GetString("tencentCOS.SecretId"), viper.GetString("tencentCOS.SecretKey"), time.Hour, opt, true)
	} else {
		presignedURL, err = client.Object.GetPresignedURL(ctx, http.MethodPut, file, viper.GetString("tencentCOS.SecretId"), viper.GetString("tencentCOS.SecretKey"), time.Hour, nil)
	}
	if err != nil {
		log.LOGE("COS get file: ", file, " ERROR: ", err)
		return "", err
	}
	return presignedURL.String(), nil
}
