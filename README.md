

### Linter run

```bash
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.8.0

golangci-lint run ./steps/...
```
