package utils

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/archive"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"net/http"
	"runtime"
)

// ExtractOCIImage will extract a given targetImage into a given targetDestination and pull from the local repo if set.
func ExtractOCIImage(targetImage, targetDestination string, isLocal bool) error {
	platform, err := v1.ParsePlatform(fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}

	ref, err := name.ParseReference(targetImage)
	if err != nil {
		return err
	}

	img, err := getimage(ref, *platform, isLocal)
	if err != nil {
		return err
	}

	reader := mutate.Extract(img)

	_, err = archive.Apply(context.Background(), targetDestination, reader)
	if err != nil {
		return err
	}
	return nil
}

// image returns the proper image to pull with transport and auth
func getimage(ref name.Reference, platform v1.Platform, local bool) (v1.Image, error) {
	if local {
		return daemon.Image(ref)
	}

	return remote.Image(ref,
		remote.WithTransport(http.DefaultTransport),
		remote.WithPlatform(platform),
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
	)
}
