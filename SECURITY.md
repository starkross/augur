# Security

## Supply chain

Every tagged release of `augur` ships with the following supply-chain artifacts:

| Artifact | Format | Where |
|----------|--------|-------|
| Binary archives | `tar.gz` / `zip` | GitHub release page |
| Linux packages | `deb`, `rpm`, `apk` | GitHub release page |
| Container images | OCI (linux/amd64, linux/arm64) | `ghcr.io/starkross/augur` |
| SBOM (per archive) | CycloneDX JSON | `<archive>.sbom.cdx.json` on the release page |
| SBOM (per image) | CycloneDX, attached as OCI referrer | discoverable via `cosign download sbom` |
| Checksums | SHA-256 | `checksums.txt` |
| Cosign signature (checksums) | sigstore keyless | `checksums.txt.sig` + `checksums.txt.pem` |
| Cosign signature (image) | sigstore keyless | OCI signature in registry |
| SLSA build provenance | in-toto bundle | `multiple.intoto.jsonl` on the release page |

Signatures are produced with **cosign keyless** (sigstore OIDC). The signing identity is the release workflow itself — there is no long-lived signing key to leak or rotate.

## Verifying a release

You will need [`cosign`](https://docs.sigstore.dev/cosign/installation/) (≥ v2.0). Replace `vX.Y.Z` with the release tag you want to verify.

### 1. Verify `checksums.txt` and then verify your binary

```sh
TAG=vX.Y.Z
BASE="https://github.com/starkross/augur/releases/download/${TAG}"

curl -fsSLO "${BASE}/checksums.txt"
curl -fsSLO "${BASE}/checksums.txt.sig"
curl -fsSLO "${BASE}/checksums.txt.pem"

cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature   checksums.txt.sig \
  --certificate-identity-regexp 'https://github.com/starkross/augur/\.github/workflows/release\.yml@refs/tags/v.*' \
  --certificate-oidc-issuer     'https://token.actions.githubusercontent.com' \
  checksums.txt

# Now verify your downloaded binary matches the signed checksum
curl -fsSLO "${BASE}/augur_${TAG#v}_linux_amd64.tar.gz"
sha256sum --ignore-missing --check checksums.txt
```

### 2. Verify the container image

```sh
cosign verify ghcr.io/starkross/augur:vX.Y.Z \
  --certificate-identity-regexp 'https://github.com/starkross/augur/\.github/workflows/release\.yml@refs/tags/v.*' \
  --certificate-oidc-issuer     'https://token.actions.githubusercontent.com'
```

### 3. Verify SLSA build provenance

```sh
# Install slsa-verifier: https://github.com/slsa-framework/slsa-verifier
slsa-verifier verify-artifact \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/starkross/augur \
  --source-tag  vX.Y.Z \
  augur_${TAG#v}_linux_amd64.tar.gz
```

### 4. Inspect the SBOM

```sh
# For an archive: download <archive>.sbom.cdx.json from the release page.

# For a container image:
cosign download sbom ghcr.io/starkross/augur:vX.Y.Z > augur.sbom.cdx.json
```

The SBOM is in [CycloneDX 1.5+](https://cyclonedx.org/) JSON. Feed it to your scanner of choice (`grype`, `trivy sbom`, Dependency-Track, etc.).

## Reporting a vulnerability

Please **do not open a public GitHub issue** for security vulnerabilities. Instead, use GitHub's private vulnerability reporting:

1. Go to https://github.com/starkross/augur/security/advisories
2. Click **Report a vulnerability**
3. Provide a description, reproduction steps, and the affected versions

We will acknowledge the report within a few working days and coordinate disclosure.
