package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/java"
)

var (
	scanCmd = &cobra.Command{
		Use:   "scan <dir>",
		Short: "Scan a directory for JDKs (max depth: 3 subdirectories)",
		Long: `Scan a specified directory for JDK installations and add them to jenv's config.

Directory Depth Limit:
The scan will only check subdirectories up to 3 levels deeper than the start directory.
For example, if scanning from "C:\":
  C:\                     (start)
  ├── Program Files      (depth 1)
  │   └── Java          (depth 2)
  │       └── jdk-21    (depth 3)
  └── Users             (depth 1)
      └── Username      (depth 2)
          └── .jdks     (depth 3)

This command will:
1. Search for JDKs in the specified directory and its subdirectories
2. Skip system directories (Windows, $Recycle.Bin, System Volume Information)
3. Add the JDKs to jenv's configuration`,
		Example: `  jenv scan C:\\
  jenv scan "C:\\Program Files\\Java"
  jenv scan C:\\Users\\Username\\.jdks`,
		Args: cobra.ExactArgs(1),
		Run:  runScan,
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) {
	dir := args[0]

	jdks := java.ScanJDK(dir)

	// Add each JDK to config
	for _, jdk := range jdks {
		fmt.Printf("Found JDK at: %s\n", jdk.Path)
		var name string
		fmt.Print("Enter a name for this JDK (e.g. jdk11, jdk21-azul): ")
		fmt.Scanln(&name)

		if name == "" {
			fmt.Println("Skipping unnamed JDK")
			continue
		}

		if err := java.AddJDK(name, jdk.Path); err != nil {
			fmt.Printf("Failed to add JDK: %v\n", err)
		} else {
			fmt.Printf("Successfully added JDK: %s -> %s\n", name, jdk.Path)
		}
	}
}
