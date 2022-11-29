package quarkus

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	d.Logger.Title(context.Buildpack)

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &d.Logger)
	if err != nil {
		return libcnb.DetectResult{}, err
	}

	// If any of environment variable that could be set by this buildpack is already set
	// this buildpack won't participate.
	if anySet(cr,
		"BP_MAVEN_BUILT_ARTIFACT",
		"BP_MAVEN_BUILD_ARGUMENTS",
		"BP_NATIVE_IMAGE_BUILD_ARGUMENTS_FILE",
		"BP_NATIVE_IMAGE_BUILT_ARTIFACT") {
		return libcnb.DetectResult{Pass: false}, nil
	}

	pomFile, ok := cr.Resolve("BP_MAVEN_POM_FILE")
	if !ok {
		pomFile = "pom.xml"
	}

	file := filepath.Join(context.Application.Path, pomFile)
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return libcnb.DetectResult{Pass: false}, nil
	} else if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to determine if %s exists\n%w", file, err)
	}

	project := struct {
		Dependencies []struct {
			GroupID    string `xml:"groupId"`
			ArtifactID string `xml:"artifactId"`
		} `xml:"dependencies>dependency"`
		Properties struct {
			QuarkusPlatformArtifactID string `xml:"quarkus.platform.artifact-id"`
		} `xml:"properties"`
	}{}

	data, err := os.ReadFile(file)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("failed to read pom file: %w", err)
	}

	err = xml.Unmarshal(data, &project)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("failed to parse pom file: %w", err)
	}

	isQuarkus := project.Properties.QuarkusPlatformArtifactID != ""

	if !isQuarkus {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
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
	}, nil
}

func anySet(cr libpak.ConfigurationResolver, vars ...string) bool {
	for _, s := range vars {
		if _, ok := cr.Resolve(s); ok {
			return true
		}
	}
	return false
}
