# SplitJoin CLI
SplitJoin is a productivity focused service for helping developers communicate faster and more effectively.

## CLI Development
### Releasing
Releases are automated with goreleaser. Make sure to push tags to trigger builds.

```
git commit -m "update xyz"
git push
git tag -a v0.2 -m "release v0.2"
git push origin v0.2
```