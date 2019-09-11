package util

import (
	"fmt"
	"os"
	"path"
)

func RelativePath(relpath string) string {
	fmt.Println(path.Join(os.Getenv("root_file_path"), relpath))
	return path.Join(os.Getenv("root_file_path"), relpath)
}
