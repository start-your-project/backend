package composites

import (
	"context"
	"log"
	"main/internal/constants"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioComposite struct {
	client *minio.Client
}

func NewMinioComposite() (*MinioComposite, error) {
	minioClient, err := minio.New(os.Getenv("MINIOURL"), &minio.Options{
		Creds:              credentials.NewStaticV4(os.Getenv("MINIOUSER"), os.Getenv("MINIOPASSWORD"), ""),
		Secure:             false,
		Transport:          nil,
		Region:             "",
		BucketLookup:       0,
		CustomRegionViaURL: nil,
		TrailingHeaders:    false,
		CustomMD5:          nil,
		CustomSHA256:       nil,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(context.Background(), constants.UserObjectsBucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), constants.UserObjectsBucketName, minio.MakeBucketOptions{
			Region:        "",
			ObjectLocking: false,
		})
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(constants.DefaultImage)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	opts := minio.PutObjectOptions{
		UserMetadata:            map[string]string{"x-amz-acl": "public-read"},
		UserTags:                nil,
		Progress:                nil,
		ContentType:             "image/png",
		ContentEncoding:         "",
		ContentDisposition:      "",
		ContentLanguage:         "",
		CacheControl:            "",
		Mode:                    "",
		RetainUntilDate:         time.Time{},
		ServerSideEncryption:    nil,
		NumThreads:              0,
		StorageClass:            "",
		WebsiteRedirectLocation: "",
		PartSize:                0,
		LegalHold:               "",
		SendContentMd5:          false,
		DisableContentSha256:    false,
		DisableMultipart:        false,
		ConcurrentStreamParts:   false,
		Internal: minio.AdvancedPutOptions{
			SourceVersionID:    "",
			SourceETag:         "",
			ReplicationStatus:  "",
			SourceMTime:        time.Time{},
			ReplicationRequest: false,
			RetentionTimestamp: time.Time{},
			TaggingTimestamp:   time.Time{},
			LegalholdTimestamp: time.Time{},
		},
	}

	_, err = minioClient.PutObject(
		context.Background(),
		constants.UserObjectsBucketName, // Константа с именем бакета
		constants.DefaultImage,
		file,
		stat.Size(),
		opts,
	)

	if err != nil {
		return nil, err
	}

	return &MinioComposite{client: minioClient}, nil
}
