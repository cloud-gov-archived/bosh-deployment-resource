package applyspec_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"encoding/json"

	. "github.com/cloudfoundry/bosh-agent/agent/applier/applyspec"
	"github.com/cloudfoundry/bosh-agent/agent/applier/models"
	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
)

var _ = Describe("RenderedTemplatesArchive", func() {
	var (
		r RenderedTemplatesArchiveSpec
		d *boshcrypto.MultipleDigest
	)

	BeforeEach(func() {
		a := boshcrypto.MustNewMultipleDigest(boshcrypto.NewDigest(boshcrypto.DigestAlgorithmSHA1, "abc"))
		d = &a
	})

	JustBeforeEach(func() {
		r = RenderedTemplatesArchiveSpec{
			Sha1:        d,
			BlobstoreID: "123",
			SignedURL:   "signed/url",
		}
	})

	It("provides jobs as source", func() {
		asSource := r.AsSource(models.Job{Name: "foo"})
		Expect(asSource.Sha1.String()).To(Equal("abc"))
		Expect(asSource.BlobstoreID).To(Equal("123"))
		Expect(asSource.SignedURL).To(Equal("signed/url"))
		Expect(asSource.PathInArchive).To(Equal("foo"))
	})

	Context("when the digest is nil", func() {
		BeforeEach(func() {
			d = nil
		})

		It("does not try to derefence it", func() {
			Expect(func() { r.AsSource(models.Job{Name: "foo"}) }).ShouldNot(Panic())
		})
	})

	Context("unmarshalling JSON", func() {
		DescribeTable("unmarshalling", func(blobstoreID, signedURL, sha1, errorMsg string, expected *RenderedTemplatesArchiveSpec) {
			data := []byte(fmt.Sprintf(`{"blobstore_id": "%s", "signed_url": "%s", "sha1": "%s"}`, blobstoreID, signedURL, sha1))
			var rendered *RenderedTemplatesArchiveSpec
			rendered = &RenderedTemplatesArchiveSpec{}
			err := json.Unmarshal(data, rendered)
			if errorMsg == "" {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(errorMsg))
			}
			if expected == nil {
				Expect(rendered).To(BeNil())
			} else {
				Expect(*rendered).To(Equal(*expected))
			}
		},
			Entry("When signedURL and Sha1Sum are set", "", "signed/url", "abc", "", &RenderedTemplatesArchiveSpec{
				Sha1:        newMultipleSha1Digest("abc"),
				BlobstoreID: "",
				SignedURL:   "signed/url",
			}),
			Entry("When signedURL is set", "", "signed/url", "", "", &RenderedTemplatesArchiveSpec{
				Sha1:        nil,
				BlobstoreID: "",
				SignedURL:   "",
			}),
			Entry("When only sha1Sum is set", "", "", "abc", "", &RenderedTemplatesArchiveSpec{
				Sha1:        newMultipleSha1Digest("abc"),
				BlobstoreID: "",
				SignedURL:   "",
			}),
			Entry("When nothing is set", "", "", "", "", &RenderedTemplatesArchiveSpec{
				Sha1:        nil,
				BlobstoreID: "",
				SignedURL:   "",
			}),
			Entry("When everything is set", "123", "signed/url", "abc", "", &RenderedTemplatesArchiveSpec{
				Sha1:        newMultipleSha1Digest("abc"),
				BlobstoreID: "123",
				SignedURL:   "signed/url",
			}),
			Entry("When blobstoreID and signedURL are set", "123", "signed/url", "", "No digest algorithm found. Supported algorithms: sha1, sha256, sha512", &RenderedTemplatesArchiveSpec{
				Sha1:        nil,
				BlobstoreID: "",
				SignedURL:   "",
			}),
			Entry("When blobstoreID and Sha1Sum are set", "123", "", "abc", "", &RenderedTemplatesArchiveSpec{
				Sha1:        newMultipleSha1Digest("abc"),
				BlobstoreID: "123",
				SignedURL:   "",
			}),
			Entry("When only blobstoreID is set", "123", "", "", "", &RenderedTemplatesArchiveSpec{
				Sha1:        nil,
				BlobstoreID: "",
				SignedURL:   "",
			}),
		)
	})
})

func newMultipleSha1Digest(sha1 string) *boshcrypto.MultipleDigest {
	a := boshcrypto.MustNewMultipleDigest(boshcrypto.NewDigest(boshcrypto.DigestAlgorithmSHA1, sha1))
	return &a
}
