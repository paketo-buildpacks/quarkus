package quarkus

import (
	"fmt"
	"strconv"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	r := libpak.PlanEntryResolver{Plan: context.Plan}
	_, ok, err := r.Resolve("quarkus")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve buildpack plan entry quarkus: %w", err)
	} else if !ok {
		return libcnb.BuildResult{}, nil
	}

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	var isNativeBuild bool
	if n, ok := cr.Resolve("BP_NATIVE_IMAGE"); ok {
		if n, err := strconv.ParseBool(n); err == nil {
			isNativeBuild = n
		}
	}

	result.Layers = append(result.Layers, &quarkusEnvVarLayer{isNativeBuild: isNativeBuild})
	for _, s := range []string{"lib", "quarkus-run.jar", "app", "quarkus"} {
		result.Slices = append(result.Slices, libcnb.Slice{Paths: []string{s}})
	}

	return result, nil
}

type quarkusEnvVarLayer struct {
	isNativeBuild bool
}

func (q *quarkusEnvVarLayer) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	layer.Build = true
	if q.isNativeBuild {
		layer.BuildEnvironment.Default("BP_MAVEN_BUILD_ARGUMENTS", "package -DskipTests=true -Dmaven.javadoc.skip=true -Dquarkus.package.type=native-sources")
		layer.BuildEnvironment.Default("BP_MAVEN_BUILT_ARTIFACT", "target/native-sources/*")
		layer.BuildEnvironment.Default("BP_NATIVE_IMAGE_BUILD_ARGUMENTS_FILE", "native-image.args")
		layer.BuildEnvironment.Default("BP_NATIVE_IMAGE_BUILT_ARTIFACT", "*-runner.jar")
	} else {
		layer.BuildEnvironment.Default("BP_MAVEN_BUILT_ARTIFACT", "target/quarkus-app/lib/ target/quarkus-app/*.jar target/quarkus-app/app/ target/quarkus-app/quarkus/")
		layer.BuildEnvironment.Default("BP_MAVEN_BUILD_ARGUMENTS", "package -DskipTests=true -Dmaven.javadoc.skip=true -Dquarkus.package.type=fast-jar")
	}

	return layer, nil
}

func (q *quarkusEnvVarLayer) Name() string {
	return "quarkus-build-envvars"
}
