package files

import (
	"os"
)

// File contém informações sobre um arquivo.
type File struct {
	Name string
	Size int64
	Type string // file or directory
}

// ListDirectory lista os arquivos em um diretório e retorna uma lista de estruturas File.
func ListDirectory(path string) ([]File, error) {
	var fileList []File

	// Abre o diretório especificado
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Lê os arquivos e subdiretórios no diretório
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Itera sobre os arquivos e subdiretórios
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

// IsFolder verifica se o caminho especificado é um diretório.
func IsFolder(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
