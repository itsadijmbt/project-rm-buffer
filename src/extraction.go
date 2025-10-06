package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func extractFromContainer(filename string) error {
	fmt.Printf("\n\n ################## EXTRACTING ARTIFACT ###################################\n\n")

	// /root/backup/<filename>/.<DEPEND>
	bashCommand := filepath.Join(BACKUP, filename, DEPEND)
	fmt.Printf(" the file is %s\n", bashCommand)
	extractIt := exec.Command("bash", bashCommand)
	//extractIt.Stdout = os.Stdout
	//extractIt.Stderr = os.Stderr

	if err := extractIt.Run(); err != nil {
		fmt.Printf(" failed to extract the file\n\n")
		return err
	}

	fmt.Printf("\n\n ######################### RESOURCE CLEANUP ############################\n\n")

	bashCommand = filepath.Join(BACKUP, filename)
	cleanResource := exec.Command("rm", "-rf", bashCommand)

	if err := cleanResource.Run(); err != nil {
		fmt.Printf("failed resource cleanup\n\n")
		return err
	}

	return nil
}

func extractList() error {
	fmt.Printf("\n\n ######################## AVAILABLE CONTAINERS ###################################")

	listIt := exec.Command("ls", "-l", BACKUP)
	listIt.Stdout = os.Stdout
	listIt.Stderr = os.Stderr

	if err := listIt.Run(); err != nil {
		fmt.Printf(" failed to list\n")
		return err
	}
	return nil
}
