package find

import "os"

func isExecutable(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.Mode()&0111 != 0
}
