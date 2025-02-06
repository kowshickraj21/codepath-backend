package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)


func ReadFile(ctx context.Context, client *s3.Client, bucket, path string) string {

	result, err := client.GetObject(ctx,&s3.GetObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(path),
	})

	if err != nil { 
		fmt.Println("Err: ",err)
		return ""
	 }
	defer result.Body.Close()
	body,err := io.ReadAll(result.Body)
	if err != nil { return "" }
	return string(body);
}