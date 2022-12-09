package quarkus_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/quarkus/quarkus"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect quarkus.Detect
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = ioutil.TempDir("", "quarkus")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
	})

	it("passes", func() {
		pom, err := os.ReadFile(filepath.Join("testdata", "quarkus-pom.xml"))
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(ctx.Application.Path, "pom.xml"), pom, 0644)
		if err != nil {
			t.Fatal(err)
		}

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "quarkus"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "quarkus"},
					},
				},
			},
		}))
	})

	it("does not pass", func() {
		pom, err := os.ReadFile(filepath.Join("testdata", "spring-boot-pom.xml"))
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(ctx.Application.Path, "pom.xml"), pom, 0644)
		if err != nil {
			t.Fatal(err)
		}

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: false,
		}))
	})
}
