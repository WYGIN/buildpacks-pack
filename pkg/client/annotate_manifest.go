package client

import (
	"context"
	"fmt"
	"strings"

	ggcrName "github.com/google/go-containerregistry/pkg/name"
)

type ManifestAnnotateOptions struct {
	OS, OSVersion, OSArch, OSVariant string
	OSFeatures, Features             []string
	Annotations                      map[string]string
}

// AnnotateManifest implements commands.PackClient.
func (c *Client) AnnotateManifest(ctx context.Context, name string, image string, opts ManifestAnnotateOptions) error {
	manifestList, err := c.indexFactory.FindIndex(name)
	if err != nil {
		return err
	}

	digest, err := ggcrName.NewDigest(image)
	if err != nil {
		return err
	}

	if opts.OS != "" {
		if err := manifestList.SetOS(digest, opts.OS); err != nil {
			return err
		}
	}
	if opts.OSVersion != "" {
		if err := manifestList.SetOSVersion(digest, opts.OSVersion); err != nil {
			return err
		}
	}
	if len(opts.OSFeatures) != 0 {
		if err := manifestList.SetOSFeatures(digest, opts.OSFeatures); err != nil {
			return err
		}
	}
	if opts.OSArch != "" {
		if err := manifestList.SetArchitecture(digest, opts.OSArch); err != nil {
			return err
		}
	}
	if opts.OSVariant != "" {
		if err := manifestList.SetVariant(digest, opts.OSVariant); err != nil {
			return err
		}
	}
	if len(opts.Features) != 0 {
		if err := manifestList.SetFeatures(digest, opts.Features); err != nil {
			return err
		}
	}
	if len(opts.Annotations) != 0 {
		annotations := make(map[string]string)
		for _, annotationSpec := range opts.Annotations {
			spec := strings.SplitN(annotationSpec, "=", 2)
			if len(spec) != 2 {
				return fmt.Errorf("no value given for annotation %q", spec[0])
			}
			annotations[spec[0]] = spec[1]
		}
		if err := manifestList.SetAnnotations(digest, annotations); err != nil {
			return err
		}
	}

	err = manifestList.Save()
	if err == nil {
		fmt.Printf("%s annotated \n", digest.String())
	}

	return err
}
