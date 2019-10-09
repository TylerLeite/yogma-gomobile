# Ye Olde Greater Metropolitan Area

### Run
To run and hide a bunch of random deprecation errors (from Mac deprecating OpenGL)
```
go run main.go 2>&1 | grep -v deprecated | grep -v '# golang.org/x/mobile'
```
