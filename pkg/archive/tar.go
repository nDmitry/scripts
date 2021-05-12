package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Tar struct{}

func (t *Tar) TarGzip(src string, out string) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files: %w", err)
	}

	f, err := os.Create(out)

	if err != nil {
		return fmt.Errorf("Unable to create the output file: %w", err)
	}

	gzw := gzip.NewWriter(f)

	defer gzw.Close()

	tw := tar.NewWriter(gzw)

	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Error while walking the directory: %w", err)
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())

		if err != nil {
			return fmt.Errorf("Could not create file header: %w", err)
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("Could not write file header: %w", err)
		}

		f, err := os.Open(file)

		if err != nil {
			return fmt.Errorf("Could not open file for taring: %w", err)
		}

		if _, err := io.Copy(tw, f); err != nil {
			return fmt.Errorf("Could not copy file into the archive: %w", err)
		}

		// manually close here after each file operation
		// defering would cause each file close to wait until all operations have completed
		f.Close()

		return nil
	})
}
