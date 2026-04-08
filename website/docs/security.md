---
id: security
title: Security
sidebar_position: 7
---

# Security

Every augur release is reproducibly built and signed:

- **CycloneDX SBOMs** for every archive and every container image
- **Cosign keyless signatures** (sigstore OIDC) on `checksums.txt` and on the published OCI image
- **SLSA Level 3 build provenance** for binary artifacts

## Verify a container image

```sh
cosign verify ghcr.io/starkross/augur:vX.Y.Z \
  --certificate-identity-regexp 'https://github.com/starkross/augur/\.github/workflows/release\.yml@refs/tags/v.*' \
  --certificate-oidc-issuer     'https://token.actions.githubusercontent.com'
```

See [SECURITY.md](https://github.com/starkross/augur/blob/main/SECURITY.md) in the repository for the full verification recipe covering binaries, images, SLSA provenance, and SBOMs.

## Reporting a vulnerability

Please report suspected vulnerabilities using [GitHub private vulnerability reporting](https://github.com/starkross/augur/security/advisories/new) rather than opening a public issue. See [SECURITY.md](https://github.com/starkross/augur/blob/main/SECURITY.md) for the disclosure process.
