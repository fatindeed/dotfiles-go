package fs

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

var (
	s3Client   *s3.Client
	s3Instance *s3Adapter
)

func getS3Client() *s3.Client {
	if s3Client != nil {
		return s3Client
	}
	options := s3.Options{
		Credentials: credentials.NewStaticCredentialsProvider(viper.GetString("s3.access_key_id"), viper.GetString("s3.secret_access_key"), ""),
		Region:      viper.GetString("s3.region"),
	}
	if url := viper.GetString("s3.endpoint"); url != "" {
		options.EndpointResolver = s3.EndpointResolverFromURL(url)
	}
	s3Client = s3.New(options)
	return s3Client
}

type s3Adapter struct {
	*s3.Client
}

func (c *s3Adapter) ReadFile(name string) ([]byte, error) {
	part := strings.SplitN(name, "/", 2)
	result, err := c.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &part[0],
		Key:    &part[1],
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

func getS3Adapter() (*s3Adapter, error) {
	if s3Instance != nil {
		return s3Instance, nil
	}
	s3Instance = new(s3Adapter)
	s3Instance.Client = getS3Client()
	return s3Instance, nil
}
