// +build !windows

package find

import (
	"io/fs"
	"path/filepath"
	"runtime"

	model "github.com/adedayo/softaudit/pkg/model"
)

func Executables(roots model.ExecPaths) <-chan string {

	channels := make([]<-chan string, 0)

	for _, root := range roots.PathRoots {
		channels = append(channels, findExecs(root))
	}
	return mergeStringChannels(channels...)
}

func findExecs(root string) <-chan string {

	out := make(chan string)
	go func() {
		defer close(out)
		filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				if isExecutable(path) {
					out <- path
				}
			} else if runtime.GOOS == "darwin" { //macos .app files are "directories"
				// println(path)
				if isExecutable(path) {
					// println("Executable: ", path)
					out <- path
				}
			}
			return nil
		})

	}()
	return out
}
