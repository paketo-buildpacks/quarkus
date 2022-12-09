# `gcr.io/paketo-buildpacks/quarkus`

The Paketo Quarkus Buildpack is a Cloud Native Buildpack that sets environment variables for the `maven` and `native-image` buildpacks when building Quarkus app.

## Behavior

This buildpack will participate if all of the following conditions are met

* The project is Quarkus Maven project
* None of the following environment variables are explicitly set:
  * `BP_MAVEN_BUILD_ARGUMENTS`
  * `BP_MAVEN_BUILT_ARTIFACT`
  * `BP_NATIVE_IMAGE_BUILD_ARGUMENTS_FILE`
  * `BP_NATIVE_IMAGE_BUILT_ARTIFACT`

The buildpack will do the following:

* Sets required `BP_MAVEN_*` and `BP_NATIVE_IMAGE_BUILD_*` environment variables.

## Configuration

| Environment Variable | Description                                                                                                                                                                                                                        |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `$BP_MAVEN_POM_FILE` | Specifies a custom location to the project's `pom.xml` file. It should be a full path to the file under the `/workspace` directory or it should be relative to the root of the project (i.e. `/workspace'). Defaults to `pom.xml`. |

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

