package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GetSSHKeyAuth(sshKeyPath string) (ssh.AuthMethod, error) {
	// Load the private key file
	key, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

func SCPCopyFolder(client *ssh.Client, localPath, remotePath string) error {
	// Start an SCP session using `scp -r` for recursive copy
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Use `scp` to copy the folder
	cmd := fmt.Sprintf("scp -r -t %s", remotePath)
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	session.Stdout = os.Stdout
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	err = session.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start SCP command: %v", err)
	}

	// Write files and directories recursively to SCP
	err = writeFolderToSCP(stdin, localPath)
	if err != nil {
		return fmt.Errorf("failed to write folder to SCP: %v", err)
	}

	// Create a channel to monitor errors
	errorChan := make(chan error, 1)
	stderrChan := make(chan string, 1)
	writeErr := make(chan error, 1)

	// Goroutine to read from stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stderrChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("error reading stderr: %v", err)
		}
		close(stderrChan)
	}()

	// Goroutine to monitor session closure
	go func() {
		err := session.Wait()
		errorChan <- err
		close(errorChan)
	}()

	// Write files and directories recursively to SCP
	go func() {
		writeErr <- writeFolderToSCP(stdin, localPath)
		stdin.Close()
		close(writeErr)
	}()

	// Wait for either stderr output, session closure, or write completion
	for {
		select {
		case line, ok := <-stderrChan:
			if ok {
				fmt.Printf("STDERR: %s\n", line)
			}
		case err := <-errorChan:
			if err != nil {
				return fmt.Errorf("session closed with error: %v", err)
			}
			return nil
		case err := <-writeErr:
			if err != nil {
				return fmt.Errorf("failed to write folder to SCP: %v", err)
			}
		}
	}
}

// writeFolderToSCP writes files and folders recursively to the SCP stdin pipe
func writeFolderToSCP(stdin io.WriteCloser, localPath string) error {
	return filepath.WalkDir(localPath, func(path string, file os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %q: %v", path, err)
		}

		relativePath, err := filepath.Rel(localPath, path)
		if err != nil {
			return err
		}
		folder := filepath.Dir(relativePath)
		// Skip the root folder itself
		if relativePath == "." {
			return nil
		}

		if file.IsDir() {
			return nil
		}

		folders := strings.Split(folder, "/")
		if folders[0] == "." {
			folders = folders[1:]
		}

		for _, folder := range folders {
			fmt.Fprintf(stdin, "D0755 0 %s\n", folder)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %v", path, err)
		}

		// Send file header and contents
		fmt.Fprintf(stdin, "C0644 %d %s\n", len(data), file.Name())
		stdin.Write(data)
		fmt.Fprint(stdin, "\x00") // End of file

		for i := 0; i < len(folders); i++ {
			fmt.Fprintf(stdin, "E\n")
		}

		return nil
	})
}
