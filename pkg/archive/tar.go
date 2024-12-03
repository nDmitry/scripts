package archive

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
- Using buffered I/O using `bufio.NewReader` and `bufio.NewWriter` to reduce system calls
- Using a fixed-size buffer (32KB) for copying files instead of letting `io.Copy` handle it
- Processing one file at a time and closing it immediately after use
- Explicit buffer flushing and proper closing of all writers
- Using `gzip.DefaultCompression` which provides a good balance between CPU usage and compression ratio
*/

type Tar struct{}

func (t *Tar) TarGzip(src string, out string) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files: %w", err)
	}

	f, err := os.Create(out)

	if err != nil {
		return fmt.Errorf("Unable to create the output file: %w", err)
	}

	defer f.Close()

	bufW := bufio.NewWriter(f)
	defer bufW.Flush()

	gzw, err := gzip.NewWriterLevel(bufW, gzip.DefaultCompression)

	if err != nil {
		return fmt.Errorf("Unable to create gzip writer: %w", err)
	}

	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	const bufferSize = 32 * 1024 // 32KB buffer
	buffer := make([]byte, bufferSize)

	err = filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
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

		// use buffered reader for better performance
		bufR := bufio.NewReader(f)

		// copy file contents in chunks
		for {
			n, err := bufR.Read(buffer)

			if err != nil && err != io.EOF {
				f.Close()
				return fmt.Errorf("Error reading file: %w", err)
			}

			if n == 0 {
				break
			}

			if _, err := tw.Write(buffer[:n]); err != nil {
				f.Close()
				return fmt.Errorf("Error writing to tar: %w", err)
			}
		}

		// manually close here after each file operation
		// defering would cause each file close to wait until all operations have completed
		f.Close()

		return nil
	})

	if err != nil {
		return err
	}

	// ensure everything is written
	if err := tw.Close(); err != nil {
		return fmt.Errorf("Error closing tar writer: %w", err)
	}

	if err := gzw.Close(); err != nil {
		return fmt.Errorf("Error closing gzip writer: %w", err)
	}

	if err := bufW.Flush(); err != nil {
		return fmt.Errorf("Error flushing buffer: %w", err)
	}

	return nil
}
