package services

import (
	"fmt"
	"os/exec"
)

func SignMessage(message string) (string, error) {
	return callJavaFunction(message)
}

func callJavaFunction(message string) (string, error) {
	// Run Java program with message as argument
	cmd := exec.Command("java", "-jar", "/home/roger/projects/lb/mTLS/java/signer.jar", "-a", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Java: %v", err)
	}

	return string(output), nil
}
