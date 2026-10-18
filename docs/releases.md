# Releases

`dotfiles self-update` downloads the latest release published on GitHub, checks its
SHA-256 and swaps the running binary. Until a release exists it stops with
_no release published yet_ — there is nothing to update to.

The CLI and the image are released separately, each from its own tag on `main`:

```bash
git tag -a cli-v0.2.0 -m "cli-v0.2.0"   # CLI: GitHub release
git tag -a img-v0.1.0 -m "img-v0.1.0"   # image: ghcr.io/notfoundy/dotfiles
git push origin cli-v0.2.0 img-v0.1.0
```

A `cli-v*` tag runs `.github/workflows/release.yml`. GoReleaser OSS only accepts semver tags,
so it builds against a local `v0.2.0` tag that is never pushed, and `gh` publishes the
`cli-v0.2.0` release, its notes listing the commits that touched `cli/` since the previous
CLI tag (a `-rc.1` style suffix makes it a pre-release). An `img-v*` tag also pushes the image
as `0.1.0`, next to `latest`, `44` and `44-YYYYMMDD`; it creates no GitHub release, which would
take `/releases/latest` away from the CLI.

Assets are published **without a version in their name** (`dotfiles_linux_amd64.tar.gz`,
`checksums.txt`). That is deliberate: it keeps the `/releases/latest/download/` URLs
stable, so `self-update` needs no GitHub API call. Renaming them in `.goreleaser.yaml`
breaks `cli/internal/selfupdate`.

To rehearse a release without publishing anything:

```bash
goreleaser release --snapshot --clean --skip=publish
```
