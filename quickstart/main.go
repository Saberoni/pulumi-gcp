package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a GCP resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, "my-bucket", &storage.BucketArgs{
			Location: pulumi.String("US"),
			Website: storage.BucketWebsiteArgs{
				MainPageSuffix: pulumi.String("index.html"),
			},
			UniformBucketLevelAccess: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Add index.html to bucket
		bucketObject, err := storage.NewBucketObject(ctx, "index.html", &storage.BucketObjectArgs{
			Bucket: bucket.Name,
			Source: pulumi.NewFileAsset("index.html"),
		})
		if err != nil {
			return err
		}

		// Create IAM for anonymous access to bucket
		_, err = storage.NewBucketIAMBinding(ctx, "my-bucket-binding", &storage.BucketIAMBindingArgs{
			Bucket: bucket.Name,
			Role:   pulumi.String("roles/storage.objectViewer"),
			Members: pulumi.StringArray{
				pulumi.String("allUsers"),
			},
		})
		if err != nil {
			return err
		}
		// Export the DNS name of the bucket
		ctx.Export("bucketName", bucket.Url)
		ctx.Export("bucketEndpoint", pulumi.Sprintf("http://storage.googleapis.com/%s/%s", bucket.Name, bucketObject.Name))
		return nil

	})
}
