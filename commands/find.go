package commands

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	directory   string
	deleteFiles bool
	priority    string
	force       bool
)

func NewCmdSearch() *cobra.Command {

	cmd := &cobra.Command{
		Use:    "search",
		Short:  "searches for duplicate files",
		Long:   "searches for duplicate files in a directory",
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(cmd, args)
		},
	}
	cmd.Flags().StringVarP(&directory, "directory", "d", ".", "The root directory to start the deduplication process")
	cmd.Flags().BoolVarP(&deleteFiles, "delete", "x", false, "Delete duplicate files (use with caution)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Do not prompt for confirmation when deleting files (use with caution)")
	cmd.Flags().StringVarP(&priority, "priority", "p", "", "Configuration file for authority priorities")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) {

	fileMap := make(map[string][]string)

	// Count the number of files to process
	totalFiles := countFiles(directory)

	// Create a new progress bar
	bar := pb.StartNew(int(totalFiles))

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			hash, err := calculateSHA256(path)
			if err != nil {
				return err
			}
			fileMap[hash] = append(fileMap[hash], path)
			bar.Increment()
		}
		return nil
	})
	bar.Finish()

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
		return
	}
	totalFilesToDelete := 0
	totalFilesDeleted := 0

	// Identify and handle identical files
	for hash, paths := range fileMap {
		if len(paths) > 1 {
			fmt.Printf("Identical files (SHA256: %s):\n", hash)

			paths := ClosestPath((path.Join(directory, priority)), paths)

			for i, path := range paths {
				deleteIndicator := " "
				if i > 0 {
					deleteIndicator = "x"
					totalFilesToDelete++
					color.New(color.FgRed).Add(color.Bold).Printf("%s   %s\n", deleteIndicator, path)
					if deleteFiles && confirmDelete() {
						err := deleteDuplicates(path)
						if err != nil {
							fmt.Printf("Error deleting file %s: %v\n", path, err)
						}
						totalFilesDeleted++
					}
				} else {
					deleteIndicator = "-"
					color.New(color.FgGreen).Add(color.Bold).Printf("%s   %s\n", deleteIndicator, path)
				}
			}
			fmt.Println()
		}
	}
	if deleteFiles {
		fmt.Println("Deleted files:", totalFilesDeleted)
	} else {
		fmt.Println("Files to delete:", totalFilesToDelete)
	}
}

// CommonPrefixLength calculates the length of the common prefix between two paths.
func CommonPrefixLength(a, b string) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

// ClosestPath finds the path closest to the specified path in the given list of paths.
func ClosestPath(specifiedPath string, paths []string) []string {
	sort.SliceStable(paths, func(i, j int) bool {
		prefixLengthI := CommonPrefixLength(specifiedPath, paths[i])
		prefixLengthJ := CommonPrefixLength(specifiedPath, paths[j])
		return prefixLengthI > prefixLengthJ
	})

	return paths
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}

func countFiles(rootDir string) int {
	var count int

	_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Mode().IsRegular() {
			count++
		}
		return nil
	})

	return count
}

func confirmDelete() bool {
	if force {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to delete these files? (yes/no): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}

func deleteDuplicates(path string) error {

	err := os.Remove(path)
	if err != nil {
		return err
	}
	fmt.Println("Deleted:", path)

	return nil
}
