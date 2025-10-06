package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var BACKUP string

const DEPEND string = "dependency.sh"

/*
func ScriptShellWritter(path string, filename string)error{

	//creates a simple scipt file
	fmt.Printf("\n\n################## INJECTION A DEPENDENCY SCRIPT IN CONTAINER ###################\n\n");

	nameLOC := BACKUP+filename;
        scriptPath := filepath.Join(nameLOC,DEPEND)
        f,err := os.Create(scriptPath)
        if err != nil{
                fmt.Printf("Error in File Creation @ %s\n",scriptPath)
        }
        defer f.Close()

	//appending data into the scipt file


	fmt.Printf("\n\n################### ARTIFACT INFO IN SCRIPT INJECTION #######################\n\n");


	artifactName := filepath.Join(nameLOC,filename+".tar.gz")
	extractionPath  := filepath.Dir(path)

	script := fmt.Sprintf(`#!/bin/bash
	set -euo pipefail
	ARTIFACT="%s"
	EXP="%s"
	tar -xzvf "$ARTIFACT" -C "$EXP"
	`,artifactName , extractionPath)

	 if _, err = f.WriteString(script); err != nil {
        	return fmt.Errorf("write %s: %w", scriptPath, err)
         }
         if err := os.Chmod(scriptPath, 0o755); err != nil {
            return fmt.Errorf("chmod +x %s: %w", scriptPath, err)
         }

	fmt.Printf("\n\n################### SCRIPTING SUCCESSFULL!!! #############\n\n")

	return nil


}


func packer(path string, filename string) error {

	fmt.Printf("\n\n#################### SANDBOXING STARTED ######################################\n\n")

	//creating a directory @Config->backup with name @filename
   	createADir := exec.Command("mkdir", "-p", BACKUP+filename)
	if err := createADir.Run(); err!=nil{
		fmt.Printf(" \nError in Creating a Sandbox\n")
		log.Fatal(err)
	}


	fmt.Printf("\n\n################### CREATING ARTIFACT FOR CONTAINER ###########################\n\n")

	zipperName := filepath.Join(filepath.Dir(path), filename+".tar.gz")
	cmd := exec.Command("tar", "-czvf", zipperName, "-C", filepath.Dir(path), filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
    		fmt.Printf("tar failed: %v\n%s\n", err, string(out))
	}

//	fmt.Printf("zipped file name is %s\n",zipperName)
//	fmt.Printf(" file dir is %s\n",filepath.Dir(path))

	//################################## SCRIPT EXE #############################################
	err = ScriptShellWritter(path,filename)
	if err !=nil{
	 fmt.Printf("\nScripting issue in Shell\n")
	 return err
        }

	fmt.Printf("\n\n################### INJECTING ARTIFACT INTO CONTAINER #############\n\n")

	//this saves it as /root/backup/t5/t5.tar.gz
	//this saves it as /root/backup/test.txt/test.txt.tar.gz
	//this saves it as /root/backup/filename.ext/filename.ext.tar.gz
	dstDir := BACKUP+filename+"/";
	cmdForMove := exec.Command("mv", "-t",dstDir,zipperName)
	if err := cmdForMove.Run(); err != nil {
		fmt.Printf("src=%s\n", zipperName)
		fmt.Printf("destination is %s\n",dstDir)
	}

	fmt.Printf("\n\n####################### RESOURCE CLEANUP ###############\n\n")

	removecmd := exec.Command("rm", "-rf",path)
        err = removecmd.Run()
        if err != nil {
                fmt.Println("filename is %s\n",filename)
              //  fmt.Println("Error removing :", err)
                fmt.Printf("filename is %s\n",filename)

        }

	fmt.Printf("\n\n##################### SANBOXING COMPLETED ##########################\n\n")

	return nil

}

*/

func ScriptShellWritter(path string, filename string) error {

	fmt.Printf("\n\n################## INJECTION A DEPENDENCY SCRIPT IN CONTAINER ###################\n\n")

	// ensure path and filename are normalized
	if !filepath.IsAbs(path) {
		path = "/" + path
	}
	path = filepath.Clean(path)
	baseName := filepath.Base(filename)

	nameLOC := filepath.Clean(BACKUP + baseName)
	// ensure backup directory exists
	if err := os.MkdirAll(nameLOC, 0o755); err != nil {
		return fmt.Errorf("failed to create backup dir %s: %w", nameLOC, err)
	}

	scriptPath := filepath.Join(nameLOC, DEPEND)
	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("Error in File Creation @ %s: %w", scriptPath, err)
	}
	defer f.Close()

	fmt.Printf("\n\n################### ARTIFACT INFO IN SCRIPT INJECTION #######################\n\n")

	artifactName := filepath.Join(nameLOC, baseName+".tar.gz")
	extractionPath := filepath.Dir(path)
	extractionPath = filepath.Clean(extractionPath)

	script := fmt.Sprintf(`#!/bin/bash
set -euo pipefail
ARTIFACT="%s"
EXP="%s"
tar -xzvf "$ARTIFACT" -C "$EXP"
`, artifactName, extractionPath)

	if _, err = f.WriteString(script); err != nil {
		return fmt.Errorf("write %s: %w", scriptPath, err)
	}
	if err := os.Chmod(scriptPath, 0o755); err != nil {
		return fmt.Errorf("chmod +x %s: %w", scriptPath, err)
	}

	fmt.Printf("\n\n################### SCRIPTING SUCCESSFULL!!! #############\n\n")
	return nil
}

func packer(path string, filename string) error {

	fmt.Printf("\n\n#################### SANDBOXING STARTED ######################################\n\n")

	if !filepath.IsAbs(path) {
		path = "/" + path
	}
	path = filepath.Clean(path)
	if filename == "" {
		filename = filepath.Base(path)
	}
	baseName := filepath.Base(filename)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source does not exist: %s", path)
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}

	backupDir := filepath.Clean(BACKUP + baseName)
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return fmt.Errorf("failed to create backup dir %s: %w", backupDir, err)
	}

	fmt.Printf("\n\n################### CREATING ARTIFACT FOR CONTAINER ###########################\n\n")

	dir := filepath.Dir(path)
	zipperName := filepath.Join(dir, baseName+".tar.gz")

	cmd := exec.Command("tar", "-czvf", zipperName, "-C", dir, filepath.Base(baseName))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("tar failed: %v\n%s\n", err, string(out))
		return fmt.Errorf("tar failed: %v: %s", err, string(out))
	}

	//################################## SCRIPT EXE #############################################
	if err := ScriptShellWritter(path, baseName); err != nil {
		fmt.Printf("\nScripting issue in Shell\n")
		return err
	}

	fmt.Printf("\n\n################### INJECTING ARTIFACT INTO CONTAINER #############\n\n")

	dstDir := filepath.Clean(backupDir) + string(os.PathSeparator)
	if err := os.Rename(zipperName, filepath.Join(backupDir, filepath.Base(zipperName))); err != nil {
		cmdForMove := exec.Command("mv", "-t", dstDir, zipperName)
		if err := cmdForMove.Run(); err != nil {
			fmt.Printf("src=%s\n", zipperName)
			fmt.Printf("destination is %s\n", dstDir)
			return fmt.Errorf("mv failed from %s to %s: %w", zipperName, dstDir, err)
		}
	}

	fmt.Printf("\n\n####################### RESOURCE CLEANUP ###############\n\n")

	if err := os.RemoveAll(path); err != nil {
		fmt.Printf("Error removing %s: %v\n", path, err)
		return fmt.Errorf("cleanup rm failed: %w", err)
	}

	fmt.Printf("\n\n##################### SANDBOXING COMPLETED ##########################\n\n")

	return nil
}

func absClean(baseDir, p string) (string, error) {
	if p == "" {
		return "", errors.New("empty path")
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(baseDir, p)
	}
	return filepath.Clean(p), nil
}

func userBackupBase() string {
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return filepath.Join(h, "backupForRm")
	}
	if h := os.Getenv("HOME"); h != "" {
		return filepath.Join(h, "backupForRm")
	}
	return "/root/backupForRm"
}

func init() {
	//invokes before main
	BACKUP = userBackupBase() + "/"
	fmt.Printf("back %s\n", BACKUP)
}
func main() {

	if len(os.Args) > 3 {
		fmt.Printf("\n Entered the block\n")
		// extractionLocation := os.Args[1]
		//if strings.HasPrefix(extractionLocation, "{") && strings.HasSuffix(extractionLocation, "}") && len(extractionLocation) >= 2 {
		//    extractionLocation = extractionLocation[1 : len(extractionLocation)-1]
		// }
		//	if !strings.HasPrefix(extractionLocation, "/") {
		//    extractionLocation = "/" + extractionLocation
		//}
		extractionLocation := os.Args[1]

		if strings.HasPrefix(extractionLocation, "{") && strings.HasSuffix(extractionLocation, "}") && len(extractionLocation) >= 2 {
			extractionLocation = extractionLocation[1 : len(extractionLocation)-1]
		}

		if !filepath.IsAbs(extractionLocation) {
			extractionLocation = "/" + extractionLocation
		}
		extractionLocation = filepath.Clean(extractionLocation)

		if os.Args[2] == "-extract" || os.Args[2] == "-Extract" || os.Args[2] == "-E" {
			filename := os.Args[3]
			if strings.HasPrefix(filename, "{") && strings.HasSuffix(filename, "}") && len(filename) >= 2 {
				filename = filename[1 : len(filename)-1]
			}
			if err := extractFromContainer(filename); err != nil {
				os.Exit(1)
			}
		} else if os.Args[2] == "-list" || os.Args[2] == "-L" || os.Args[2] == "-List" {
			if err := extractList(); err != nil {
				os.Exit(1)
			}
		} else if os.Args[2] == "-buffer" || os.Args[2] == "-b" {

			rmArgs := os.Args[3:]

			_, fileOperands := splitRmArgs(rmArgs)

			if len(fileOperands) == 0 {
				fmt.Fprintln(os.Stderr, "rm-buffer: no file operands to buffer")
				os.Exit(0)
			}

			pool := NewWorkerPool(4)
			locker := NewPathLocker()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			for _, op := range fileOperands {
				if strings.HasPrefix(op, "{") && strings.HasSuffix(op, "}") && len(op) >= 2 {
					op = op[1 : len(op)-1]
				}

				absolutePath, err := absClean(extractionLocation, op)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error resolving path %s: %v\n", op, err)
					continue
				}

				fmt.Printf("\n\n  File path is \n%s,", absolutePath)

				extractedFileName := filepath.Base(absolutePath)
				lockKey := extractionLocation + "::" + op

				if err := pool.Do(ctx, func() error {
					return locker.With(lockKey, func() error {
						return packer(absolutePath, extractedFileName)
					})
				}); err != nil {
					fmt.Fprintf(os.Stderr, "packer error for %s: %v\n", absolutePath, err)
				}
			}
		}
	} else if len(os.Args) < 3 {
		fmt.Printf("rm {command} {fileAddress/Name}")
		os.Exit(2)
	}

}

/*
func main() {

       //  fmt.Printf("is this even loading\n", len(os.Args));

	//fmt.Printf("arg1 = %s\n",os.Args[1])
	//fmt.Printf("arg2 = %s\n",os.Args[2])
	//fmt.Printf("arg3 = %s\n",os.Args[3])

	//we want arg1 and arg2



	if len(os.Args) > 3 {
                fmt.Printf("\n Entered the block\n");
		extractionLocation := os.Args[1]
		filename := os.Args[3]
		extractionLocation = extractionLocation[1 : len(extractionLocation)-1]
		filename = filename[1 : len(filename)-1]

		if os.Args[2] == "-extract" || os.Args[2] == "-Extract" || os.Args[2] == "-E" {
			if err := extractFromContainer(filename); err != nil {
				os.Exit(1)
			}
		} else if os.Args[2] == "-list" || os.Args[2] == "-L" || os.Args[2] == "-List" {
			if err := extractList(); err != nil {
				os.Exit(1)
			}
		} else if os.Args[2] == "-buffer" || os.Args[2] == "-b" {

			absolutePath, err := absClean(extractionLocation, filename)
			fmt.Printf("\n\n  File path is \n%s,", absolutePath)
			if err != nil {
				fmt.Printf("Error in Path\nPlease Check!!!!\n")
					os.Exit(1)
			}

			fmt.Println("\n\n")

			checker := strings.Split(absolutePath, "/")

			// a mutex for performance lockup
			pool := NewWorkerPool(4)
			locker := NewPathLocker()

			lockKey := extractionLocation + "::" + filename
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			// extract the new filename from the absolutePath field
			extractedFileName := checker[len(checker)-1]

			if err := pool.Do(ctx, func() error {
				return locker.With(lockKey, func() error {
					return packer(absolutePath, extractedFileName)
				})
			}); err != nil {
				fmt.Fprintf(os.Stderr, "packer error: %v\n", err)
				os.Exit(1)
			}
		}
	} else if len(os.Args) < 3 {
		fmt.Printf("rm {command} {fileAddress/Name}")
		os.Exit(2)
	}

}


*/
