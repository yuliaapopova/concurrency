package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(dir, filename string) (*os.File, error) {
	path, err := Path(dir)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(path, filename)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %w", err)
	}

	return file, nil
}

func Path(dir string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(cwd, dir)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("failed to create wal folder '%s': %v", dir, err)
	}
	return path, nil
}

func SegmentNameNext(dir, segmentName string) (string, error) {
	path, err := Path(dir)
	if err != nil {
		return "", err
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed reading directory: %w", err)
	}

	filenames := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filenames = append(filenames, f.Name())
	}

	for i := len(filenames) - 1; i >= 0; i-- {
		if filenames[i] > segmentName {
			return filenames[i], nil
		} else if filenames[i] == segmentName {
			return "", nil
		}
	}
	return "", fmt.Errorf("next segment not found")
}

func SegmentLastName(dir string) (string, error) {
	path, err := Path(dir)
	files, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed reading directory: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory '%s'", dir)
	}
	return files[len(files)-1].Name(), nil
}

func WriteFile(file *os.File, data []byte) error {
	_, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed sync to file: %w", err)
	}

	return nil
}
