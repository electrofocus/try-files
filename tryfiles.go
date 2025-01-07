package tryfiles

import (
	"errors"
	"io/fs"
	"net/http"
	"strings"
)

type fileSystem struct {
	fs         http.FileSystem
	extensions []string
}

// FileSystem returns http.FileSystem interface that implements nginx try_files-like behavior.
// Read more: https://nginx.org/en/docs/http/ngx_http_core_module.html#try_files
func FileSystem(fs http.FileSystem, extensions ...string) http.FileSystem {
	return fileSystem{
		fs:         fs,
		extensions: extensions,
	}
}

func (tf fileSystem) Open(name string) (http.File, error) {
	f, err := tf.fs.Open(name)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		for _, extension := range tf.extensions {
			f, err := tf.fs.Open(strings.TrimSuffix(name, ".") + "." + extension)
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					return nil, err
				}

				continue
			}

			return f, nil
		}

		return nil, err
	}

	return f, nil
}
