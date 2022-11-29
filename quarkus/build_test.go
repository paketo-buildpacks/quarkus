package quarkus_test

import (
	"os"
	"testing"

	"github.com/paketo-buildpacks/quarkus/quarkus"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		build quarkus.Build
		ctx   libcnb.BuildContext
	)

	it("does nothing without plan", func() {
		Expect(build.Build(ctx)).To(Equal(libcnb.BuildResult{}))
	})

	it("sets envvars for JVM build", func() {
		ctx.Plan = libcnb.BuildpackPlan{
			Entries: []libcnb.BuildpackPlanEntry{
				{
					Name: "quarkus",
				},
			},
		}

		buildResult, err := build.Build(ctx)
		Expect(err).To(BeNil())
		Expect(len(buildResult.Layers)).To(Equal(1))

		layer, err := buildResult.Layers[0].Contribute(libcnb.Layer{BuildEnvironment: map[string]string{}})
		Expect(err).To(BeNil())
		Expect(layer.Build).To(Equal(true))
		Expect(allSet(layer.BuildEnvironment,
			"BP_MAVEN_BUILT_ARTIFACT.default",
			"BP_MAVEN_BUILD_ARGUMENTS.default")).To(Equal(true))
	})

	it("sets envvars for native build", func() {
		ctx.Plan = libcnb.BuildpackPlan{
			Entries: []libcnb.BuildpackPlanEntry{
				{
					Name: "quarkus",
				},
			},
		}

		defer withEnvVar("BP_NATIVE_IMAGE", "1")()

		buildResult, err := build.Build(ctx)
		Expect(err).To(BeNil())
		Expect(len(buildResult.Layers)).To(Equal(1))
		layer, err := buildResult.Layers[0].Contribute(libcnb.Layer{BuildEnvironment: map[string]string{}})
		Expect(err).To(BeNil())
		Expect(layer.Build).To(Equal(true))
		Expect(allSet(layer.BuildEnvironment,
			"BP_MAVEN_BUILT_ARTIFACT.default",
			"BP_MAVEN_BUILD_ARGUMENTS.default",
			"BP_NATIVE_IMAGE_BUILD_ARGUMENTS_FILE.default",
			"BP_NATIVE_IMAGE_BUILT_ARTIFACT.default")).To(Equal(true))
		t.Log(layer)
	})

}

func withEnvVar(key, value string) func() {
	oldNative, hadNative := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if hadNative {
			os.Setenv(key, oldNative)
		} else {
			os.Unsetenv(key)
		}
	}
}

func allSet(envs map[string]string, keys ...string) bool {
	for _, key := range keys {
		if _, ok := envs[key]; !ok {
			return false
		}
	}
	return true
}
