

### Linter run

```bash
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin latest

golangci-lint run ./steps/...
```
