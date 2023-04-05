package images

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GenerateObjectName(userID int64) string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return fmt.Sprintf("%s_%s.%s", strconv.Itoa(int(userID)), formatted, "webp")
}

func GenerateFileURL(fileName string, bucket string) (string, error) {
	return fmt.Sprintf("http://%s/%s/%s", os.Getenv("HOST_MINIO"), bucket, fileName), nil
}