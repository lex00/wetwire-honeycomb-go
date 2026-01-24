---
title: "Versioning"
---

This document explains the versioning system for wetwire-honeycomb-go.

## Semantic Versioning

The project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (X): Breaking API changes
- **MINOR** (Y): New features, backwards compatible
- **PATCH** (Z): Bug fixes, backwards compatible

**Location:** Git tags (e.g., `v0.3.0`)

**Detection:** Runtime via `debug.ReadBuildInfo()` or CLI `version` command

---

## Versioned Components

| Component | Tracked By |
|-----------|------------|
| **Package Version** | Git tags (vX.Y.Z) |
| **Query API Format** | Implicit in types |
| **Lint Rules** | Rule codes (WHC001, etc.) |

### Package Version

The main version for releases. Updated for:
- New features (calculations, filters, resource types)
- Bug fixes
- Breaking changes to public API

### Query API Format

The Honeycomb Query API format is implicitly tracked through the type definitions. When Honeycomb adds new query features, the types are updated in a minor version bump.

### Lint Rules

Lint rules are identified by code (WHC001, WHC002, etc.). New rules are added in minor versions. Rule behavior changes that affect existing code are documented in CHANGELOG.

---

## Version Resolution

The CLI determines its version using this priority:

1. **ldflags**: If built with `-ldflags "-X main.version=v0.3.0"`
2. **Build info**: If installed via `go install @version`
3. **Constant**: Fallback to hardcoded version in `main.go`

---

## Viewing Current Version

### From CLI

```bash
wetwire-honeycomb version
# or
wetwire-honeycomb --version
```

### From Go Code

```go
import "runtime/debug"

if info, ok := debug.ReadBuildInfo(); ok {
    fmt.Println(info.Main.Version)
}
```

---

## Bumping the Version

When releasing a new version:

1. Update version constant in `cmd/wetwire-honeycomb/main.go`:
   ```go
   const version = "0.4.0"
   ```

2. Update CHANGELOG.md with changes

3. Run tests:
   ```bash
   go test ./...
   golangci-lint run ./...
   ```

4. Commit and tag:
   ```bash
   git commit -am "chore: release v0.4.0"
   git tag v0.4.0
   git push && git push --tags
   ```

The tag triggers the release workflow in GitHub Actions.

---

## Release Process

1. **Update CHANGELOG.md**
   - Move items from `[Unreleased]` to new version section
   - Add release date

2. **Update version constant**
   ```go
   const version = "0.4.0"
   ```

3. **Create release commit**
   ```bash
   git add CHANGELOG.md cmd/wetwire-honeycomb/main.go
   git commit -m "chore: release v0.4.0"
   ```

4. **Tag the release**
   ```bash
   git tag v0.4.0
   git push origin main --tags
   ```

5. **GitHub Actions** automatically:
   - Builds binaries for multiple platforms
   - Creates GitHub release
   - Users can install via `go install @v0.4.0`

---

## Version Compatibility

| wetwire-honeycomb-go | Go Version | wetwire-core-go |
|----------------------|------------|-----------------|
| 0.x | 1.23+ | 0.x |

---

## See Also

- [Developer Guide](../developers/) - Full development guide
