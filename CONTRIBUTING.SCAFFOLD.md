# Contributing

## Getting Started

1. Clone the repository
2. `$ make tidy`
3. `$ make run`
4. `$ open http://localhost:8080`

## Tasks

The `Makefile` in the project root contains commands to easily run common admin tasks:

| Command        | Meaning                                                                                               |
|----------------|-------------------------------------------------------------------------------------------------------|
| `$ make tidy`  | Format all code using `go fmt` and tidy the `go.mod` file.                                            |
| `$ make audit` | Run `go vet`, `staticheck`, execute all tests and verify required modules.                            |
| `$ make build` | Build a binary for the extension. Creates a file called `extension` in the repository root directory. |
| `$ make run`   | Build and then run the created binary.                                                                |

## Releasing the Code/Docker Image

To make a new release, do the following:

 1. Update the `CHANGELOG.md`
 2. Commit and push the changelog changes.
 3. Set the tag `git tag -a vX.X.X -m vX.X.X`
 4. Push the tag.

## Releasing Helm Chart Changes

 1. Update the version number in the [Chart.yaml](./charts/steadybit-extension-scaffold/Chart.yaml)
 2. Commit and push the changes.

Changing the Helm chart without bumping the version will result in the following error:

```
> Releasing charts...
    Error: error creating GitHub release steadybit-extension-scaffold-1.0.0: POST https://api.github.com/repos/steadybit/extension-scaffold/releases: 422 Validation Failed [{Resource:Release Field:tag_name Code:already_exists Message:}]
```

## Contributor License Agreement (CLA)

In order to accept your pull request, we need you to submit a CLA. You only need to do this once. If you are submitting a pull request for the first time, just submit a Pull Request and our CLA Bot will give you instructions on how to sign the CLA before merging your Pull Request.

All contributors must sign an [Individual Contributor License Agreement](https://github.com/steadybit/.github/blob/main/.github/cla/individual-cla.md).

If contributing on behalf of your company, your company must sign a [Corporate Contributor License Agreement](https://github.com/steadybit/.github/blob/main/.github/cla/corporate-cla.md). If so, please contact us via office@steadybit.com.

If for any reason, your first contribution is in a PR created by other contributor, please just add a comment to the PR
with the following text to agree our CLA: "I have read the CLA Document and I hereby sign the CLA".
