The Paketo Quarkus Buildpack is a Cloud Native Buildpack
that sets environment variables for maven/native-image buildpacks
when building Quarkus app.

## Behavior

This buildpack will participate if all of the following conditions are met

* The project is Quarkus maven project
* None of the following envvars is explicitly set:
  * `BP_MAVEN_BUILD_ARGUMENTS`
  * `BP_MAVEN_BUILT_ARTIFACT`
  * `BP_NATIVE_IMAGE_BUILD_ARGUMENTS_FILE`
  * `BP_NATIVE_IMAGE_BUILT_ARTIFACT`

The buildpack will do the following:

* Sets required `BP_MAVEN_*` and `BP_NATIVE_IMAGE_BUILD_*` environment variables.

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

