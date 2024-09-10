package files

import (
	"os"
)

// File contains information about a file.
type File struct {
	Name string
	Size int64
	Type string // file or directory
}

// ListDirectory lists the files in a directory and returns a list of File structures.
func ListDirectory(path string) ([]File, error) {
	var fileList []File

	// Opens the specified directory
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Reads the files and subdirectories in the directory
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Iterates over files and subdirectories
	for _, fileInfo := range fileInfos {
		file := File{
			Name: fileInfo.Name(),
			Size: fileInfo.Size(),
		}
		if fileInfo.IsDir() {
			file.Type = "directory"
		} else {
			file.Type = "file"
		}
		fileList = append(fileList, file)
	}

	return fileList, nil
}

// IsFolder verify if a path is a folder.
func IsFolder(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// DeleteFile exclude a file.
func DeleteFile(path string) error {
	return os.Remove(path)
}
