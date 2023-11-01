package client

import (
	"context"

	// "github.com/buildpacks/imgutil"
)

type PushManifestOptions struct {
	Format string
	Insecure, Purge bool
}
// PushManifest implements commands.PackClient.
func (c *Client) PushManifest(ctx context.Context, index string, opts PushManifestOptions) (imageID string, err error) {
	// img, err := c.indexFactory.NewIndex(index, parseFalgsForImgUtil(opts))

	// if err != nil {
	// 	return err
	// }

	// if opts.Purge {
	// 	if err := img.Delete(); err != nil {
	// 		return err
	// 	}
	// }

	return imageID, err
}

// func parseFalgsForImgUtil(opts PushManifestOptions) (idxOptions imgutil.IndexOptions) {
// 	return idxOptions
// }
