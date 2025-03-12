package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var now = time.Now

type Segment struct {
	file           *os.File
	directory      string
	segmentSize    int
	maxSegmentSize int
}

func NewSegment(directory string, maxSegmentSize int) *Segment {
	return &Segment{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
	}
}

func (s *Segment) SegmentNext() error {
	if s.file != nil {
		s.file.Close()
	}
	if err := s.createSegmentFile(s.directory); err != nil {
		return err
	}
	return nil
}

func (s *Segment) Write(data []byte) error {
	if s.segmentSize >= s.maxSegmentSize || s.file == nil {
		if err := s.SegmentNext(); err != nil {
			return err
		}
	}

	writtenBytes, err := s.file.Write(data)
	if err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}

	err = s.file.Sync()
	if err != nil {
		return err
	}

	s.segmentSize += writtenBytes
	return nil
}

func (s *Segment) LoadData() ([][]byte, error) {
	path, err := Path(s.directory)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading directory: %w", err)
	}

	var data [][]byte
	var filenameLastFile string
	var sizeLastFile int

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		file, err := os.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed opening file: %w", err)
		}
		filenameLastFile = f.Name()
		sizeLastFile = len(file)
		data = append(data, file)
	}

	if filenameLastFile != "" {
		lastFile, err := os.OpenFile(filepath.Join(path, filenameLastFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed opening file: %w", err)
		}
		s.file = lastFile
		s.segmentSize = sizeLastFile
	} else {
		if err = s.createSegmentFile(s.directory); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *Segment) createSegmentFile(dir string) error {
	filename := fmt.Sprintf("wal_%d.wal", now().UnixMilli())
	file, err := CreateFile(dir, filename)
	if err != nil {
		return err
	}

	s.segmentSize = 0
	s.file = file
	return nil
}

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
