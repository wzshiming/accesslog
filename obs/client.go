package obs

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/wzshiming/accesslog"
)

type Client = obs.ObsClient

var (
	WithHttpClient = obs.WithHttpClient
	NewOBSClient   = obs.New
)

func ProcessAccessLogWithClient(client *Client, bucketName, dateStr string, callback func(entry accesslog.Entry[AccessLog], err error) error) error {
	bucketLoggingResult, err := client.GetBucketLoggingConfiguration(bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %w", err)
	}

	prefix := bucketLoggingResult.BucketLoggingStatus.TargetPrefix + dateStr

	continueToken := ""
	for {

		lsRes, err := client.ListObjects(&obs.ListObjectsInput{
			Marker: continueToken,
			Bucket: bucketLoggingResult.BucketLoggingStatus.TargetBucket,
			ListObjsInput: obs.ListObjsInput{
				Prefix: prefix,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		for _, object := range lsRes.Contents {
			err := (func() error {
				r, err := client.GetObject(&obs.GetObjectInput{
					GetObjectMetadataInput: obs.GetObjectMetadataInput{
						Bucket: bucketLoggingResult.BucketLoggingStatus.TargetBucket,
						Key:    object.Key,
					},
				})
				if err != nil {
					return fmt.Errorf("failed to get object: %w", err)
				}
				defer r.Body.Close()

				err = accesslog.ProcessEntries[AccessLog](r.Body, callback)
				if err != nil {
					return err
				}
				return nil
			})()
			if err != nil {
				return err
			}
		}
		if !lsRes.IsTruncated {
			break
		}
		continueToken = lsRes.NextMarker
	}
	return nil
}
