package install

import (
	_ "embed"
	"fmt"
	"io"
	"os"
)

//go:embed raspc.service
var serviceContent []byte

// Install performs the necessary tasks to install the raspc service
func Install() {
	// Create necessary directories
	err := createDirectories()
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}

	// Get the current binary path
	currentBinary, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting the current binary path:", err)
		return
	}

	// Copy the binary to /opt/raspc
	destBinary := "/opt/raspc/raspc"
	err = copyFile(currentBinary, destBinary)
	if err != nil {
		fmt.Println("Error copying the binary to /opt/raspc:", err)
		return
	}

	// Copy the conf.yml file to /etc/raspc
	err = copyFile("conf.yml", "/etc/raspc/conf.yml")
	if err != nil {
		fmt.Println("Error copying conf.yml to /etc/raspc:", err)
		return
	}

	// Write the systemd service file
	err = os.WriteFile("/etc/systemd/system/raspc.service", serviceContent, 0644)
	if err != nil {
		fmt.Println("Error creating the systemd service file:", err)
		return
	}

	// Installation completion message
	fmt.Println("Service installed successfully! Run the following commands:")
	fmt.Println("sudo systemctl enable raspc.service")
	fmt.Println("sudo systemctl start raspc.service")
}

// createDirectories creates the necessary directories for the installation
func createDirectories() error {
	directories := []string{
		"/opt/raspc",
		"/etc/raspc",
		"/etc/systemd/system",
	}

	for _, dir := range directories {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %w", src, dst, err)
	}
	return nil
}
