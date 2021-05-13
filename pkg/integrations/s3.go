package integrations

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 can upload files to any S3 compatible object storage
type S3 struct {
	Endpoints        string
	AccessKeyIDs     string
	AccessKeySecrets string
}

func (s *S3) Upload(filePath string, bucketNames string, objectName string) error {
	ctx := context.Background()
	useSSL := true
	endpoints := strings.Split(s.Endpoints, ",")
	accessKeyIDs := strings.Split(s.AccessKeyIDs, ",")
	accessKeySecrets := strings.Split(s.AccessKeySecrets, ",")
	buckets := strings.Split(bucketNames, ",")

	for i, endpoint := range endpoints {
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyIDs[i], accessKeySecrets[i], ""),
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
			ctx, buckets[i], objectName, filePath,
			minio.PutObjectOptions{ContentType: GetFileContentType(f)},
		)

		if err != nil {
			return fmt.Errorf("could not upload the file: %w", err)
		}

		log.Printf("Successfully uploaded %s of size %d to %s\n", objectName, info.Size, endpoint)
	}

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
