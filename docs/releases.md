# Releases

The CLI and the image are released separately, from tags on `main`:

```bash
git tag -a cli-v0.2.0 -m "cli-v0.2.0"   # GitHub release
git tag -a img-v0.1.0 -m "img-v0.1.0"   # ghcr.io/notfoundy/dotfiles
git push origin cli-v0.2.0 img-v0.1.0
```

- `cli-v*`: GoReleaser builds, `gh` publishes the release with the commits touching `cli/` since
  the previous CLI tag. A `-rc.1` style suffix makes a pre-release.
- `img-v*`: the image is tagged `0.1.0`, `latest`, `44` and `44-YYYYMMDD`. No GitHub release,
  which would take `/releases/latest` away from the CLI.

Asset names carry no version: `self-update` downloads from `/releases/latest/download/`.
Renaming them in `.goreleaser.yaml` breaks it.

Dry run:

```bash
goreleaser release --snapshot --clean --skip=publish
```
