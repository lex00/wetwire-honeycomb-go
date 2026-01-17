// Command watch provides auto-rebuild on source file changes.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lex00/wetwire-honeycomb-go/internal/builder"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	var outputFile string
	var interval int
	var verbose bool

	cmd := &cobra.Command{
		Use:   "watch [packages]",
		Short: "Auto-rebuild on source file changes",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			fmt.Printf("Watching %s for changes (interval: %ds)\n", path, interval)
			fmt.Println("Press Ctrl+C to stop")
			fmt.Println()

			var lastModTime time.Time
			var lastHash string

			for {
				// Get current modification state
				currentModTime, currentHash, err := getDirectoryState(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error checking files: %v\n", err)
					time.Sleep(time.Duration(interval) * time.Second)
					continue
				}

				// Check if anything changed
				if !currentModTime.Equal(lastModTime) || currentHash != lastHash {
					if lastModTime.IsZero() {
						fmt.Printf("[%s] Initial build\n", time.Now().Format("15:04:05"))
					} else {
						fmt.Printf("[%s] Changes detected, rebuilding...\n", time.Now().Format("15:04:05"))
					}

					// Build
					b, err := builder.NewBuilder(path)
					if err != nil {
						fmt.Fprintf(os.Stderr, "  Error: %v\n", err)
					} else {
						result, err := b.Build()
						if err != nil {
							fmt.Fprintf(os.Stderr, "  Build failed: %v\n", err)
						} else {
							if verbose {
								fmt.Printf("  Found %d queries\n", result.QueryCount())
							}

							if result.QueryCount() > 0 && outputFile != "" {
								// Write output
								queries := result.Queries()
								var jsonData []byte

								if len(queries) == 1 {
									q := discoveredToQuery(queries[0])
									jsonData, err = serialize.ToJSONPretty(q)
								} else {
									queryMap := make(map[string]json.RawMessage)
									for _, dq := range queries {
										q := discoveredToQuery(dq)
										data, e := serialize.ToJSON(q)
										if e != nil {
											err = e
											break
										}
										queryMap[dq.Name] = data
									}
									if err == nil {
										jsonData, err = json.MarshalIndent(queryMap, "", "  ")
									}
								}

								if err != nil {
									fmt.Fprintf(os.Stderr, "  Serialization failed: %v\n", err)
								} else {
									if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
										fmt.Fprintf(os.Stderr, "  Failed to write output: %v\n", err)
									} else {
										fmt.Printf("  Wrote %s (%d bytes)\n", outputFile, len(jsonData))
									}
								}
							} else if result.QueryCount() > 0 {
								fmt.Printf("  Build succeeded (%d queries)\n", result.QueryCount())
							} else {
								fmt.Println("  No queries found")
							}
						}
					}

					lastModTime = currentModTime
					lastHash = currentHash
				}

				time.Sleep(time.Duration(interval) * time.Second)
			}
		},
	}

	cmd.Flags().StringVar(&outputFile, "output", "", "Output file")
	cmd.Flags().IntVar(&interval, "interval", 2, "Polling interval in seconds")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func getDirectoryState(dir string) (time.Time, string, error) {
	var latestTime time.Time
	var fileList []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and non-Go files
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
		}

		fileList = append(fileList, fmt.Sprintf("%s:%d", path, info.ModTime().UnixNano()))
		return nil
	})

	if err != nil {
		return latestTime, "", err
	}

	// Create a simple hash of file states
	hash := strings.Join(fileList, "|")
	return latestTime, hash, nil
}
