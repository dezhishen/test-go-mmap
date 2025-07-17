# Test mmap with plugin
## run
### windows
```
rm -rf .\plugin\mmap_operator.exe && go build -o .\plugin\mmap_operator.exe .\plugin\main.go  && go run main.go
```
### Others
```
rm -rf .\plugin\mmap_operator && go build -o .\plugin\mmap_operator .\plugin\main.go  && go run main.go
```
