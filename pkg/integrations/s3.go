package integrations

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 can upload files to any S3 compatible object storage
type S3 struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}

func (s *S3) Upload(filePath string, bucketName string, objectName string) error {
	ctx := context.Background()
	useSSL := true

	minioClient, err := minio.New(s.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.AccessKeyID, s.AccessKeySecret, ""),
		Secure: useSSL,
	})

	if err != nil {
		return fmt.Errorf("could not initialize the client: %w", err)
	}

	f, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("could not open the file: %w", err)
	}

	defer f.Close()

	info, err := minioClient.FPutObject(
		ctx, bucketName, objectName, filePath,
		minio.PutObjectOptions{ContentType: GetFileContentType(f)},
	)

	if err != nil {
		return fmt.Errorf("could not upload the file: %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return nil
}

func GetFileContentType(out *os.File) string {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)

	if err != nil {
		log.Println(err)
		return ""
	}

	out.Seek(0, 0)

	return http.DetectContentType(buffer)
}
