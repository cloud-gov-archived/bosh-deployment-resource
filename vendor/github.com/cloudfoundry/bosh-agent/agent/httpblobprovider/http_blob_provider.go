package httpblobprovider

import (
	"fmt"
	"io"
	"net/http"
	"os"

	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

var DefaultCryptoAlgorithms = []boshcrypto.Algorithm{boshcrypto.DigestAlgorithmSHA1, boshcrypto.DigestAlgorithmSHA512}

type HTTPBlobImpl struct {
	fs               boshsys.FileSystem
	createAlgorithms []boshcrypto.Algorithm
	httpClient       *http.Client
}

func NewHTTPBlobImpl(fs boshsys.FileSystem, httpClient *http.Client) *HTTPBlobImpl {
	return NewHTTPBlobImplWithDigestAlgorithms(fs, httpClient, DefaultCryptoAlgorithms)
}

func NewHTTPBlobImplWithDigestAlgorithms(fs boshsys.FileSystem, httpClient *http.Client, algorithms []boshcrypto.Algorithm) *HTTPBlobImpl {
	return &HTTPBlobImpl{
		fs:               fs,
		createAlgorithms: algorithms,
		httpClient:       httpClient,
	}
}

func (h *HTTPBlobImpl) Upload(signedURL, filepath string) (boshcrypto.MultipleDigest, error) {
	digest, err := boshcrypto.NewMultipleDigestFromPath(filepath, h.fs, h.createAlgorithms)
	if err != nil {
		return boshcrypto.MultipleDigest{}, err
	}

	// Do not close the file in the happy path because the client.Do will handle that.
	file, err := h.fs.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return boshcrypto.MultipleDigest{}, err
	}

	stat, err := h.fs.Stat(filepath)
	if err != nil {
		defer file.Close()
		return boshcrypto.MultipleDigest{}, err
	}

	req, err := http.NewRequest("PUT", signedURL, file)
	if err != nil {
		defer file.Close()
		return boshcrypto.MultipleDigest{}, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Expect", "100-continue")
	req.ContentLength = stat.Size()

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return boshcrypto.MultipleDigest{}, err
	}
	if !isSuccess(resp) {
		return boshcrypto.MultipleDigest{}, fmt.Errorf("Error executing PUT to %s for %s, response was %#v", signedURL, file.Name(), resp)
	}

	return digest, nil
}

func (h *HTTPBlobImpl) Get(signedURL string, digest boshcrypto.Digest) (string, error) {
	file, err := h.fs.TempFile("bosh-http-blob-provider-GET")
	if err != nil {
		return "", bosherr.WrapError(err, "Creating temporary file")
	}
	defer file.Close()

	resp, err := h.httpClient.Get(signedURL)
	if err != nil {
		return file.Name(), err
	}

	if !isSuccess(resp) {
		return file.Name(), fmt.Errorf("Error executing GET to %s, response was %#v", signedURL, resp)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return file.Name(), err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return file.Name(), bosherr.WrapErrorf(err, "Rewinding file pointer to beginning")
	}

	err = digest.Verify(file)
	if err != nil {
		return file.Name(), bosherr.WrapErrorf(err, "Checking downloaded blob digest")
	}

	return file.Name(), nil
}

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
