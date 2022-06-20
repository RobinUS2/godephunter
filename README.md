# godephunter
Hunt down dependency graph

Example usage:
```
go mod graph | ./godephunter --find='github.com/your-repo/your-dep@v1.1.1'
```