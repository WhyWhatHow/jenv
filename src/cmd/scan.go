package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"
)

var (
	scanCmd = &cobra.Command{
		Aliases: []string{"sc"},
		Use:     "scan <dir>",
		Short:   "Scan a directory for JDKs (max depth: 5 subdirectories)",
		Long: `Scan a specified directory for JDK installations and add them to jenv's config.

Directory Depth Limit:
The scan will only check subdirectories up to 3 levels deeper than the start directory.
For example, if scanning from "C:":
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
  jenv scan C:\\Users\\Username\\.jdks
  jenv sc  C:\\Program Files\\Java `,
		Args: cobra.ExactArgs(1),
		Run:  runScan,
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) {
	dir := args[0]

	// 显示扫描标题
	header := style.Header.Render("🔍 Scanning directory: ") + style.Path.Render(dir)
	fmt.Println(header + "\n" + strings.Repeat("─", 50))

	// Show scanning progress message
	fmt.Println(style.Input.Render("⏳ Scanning for JDK installations..."))
	fmt.Println(style.Input.Render("   • Excluding already registered JDKs"))
	fmt.Println(style.Input.Render("   • Skipping system directories"))
	fmt.Println()

	// Use the new optimized scan with statistics
	result := java.ScanJDKWithStats(dir)

	// Display scan statistics
	fmt.Printf("%s\n", style.Header.Render("📊 Scan Results"))
	fmt.Printf("%s: %s\n", style.Name.Render("⏱️  Scan Duration"), style.Success.Render(result.Duration.String()))
	fmt.Printf("%s: %d\n", style.Name.Render("📁 Directories Scanned"), result.Scanned)
	fmt.Printf("%s: %d\n", style.Name.Render("⚠️  Directories Skipped"), result.Skipped)
	fmt.Printf("%s: %d\n", style.Name.Render("🚫 Paths Excluded (duplicates)"), result.Excluded)
	fmt.Printf("%s: %d\n", style.Name.Render("🎯 New JDKs Found"), len(result.JDKs))
	fmt.Println(strings.Repeat("─", 50))

	if len(result.JDKs) == 0 {
		fmt.Println(style.Input.Render("✨ No new JDK installations found."))
		if result.Excluded > 0 {
			fmt.Println(style.Input.Render("   All discovered JDKs are already registered."))
		}
		return
	}

	successCount := 0
	skipCount := 0

	for i, jdk := range result.JDKs {
		// 显示带编号的JDK发现信息
		fmt.Printf("\n%s %s\n",
			style.Name.Render(fmt.Sprintf("#%02d", i+1)),
			style.Path.Render(jdk.Path))

		// 带样式的输入提示
		prompt := style.Input.Render("⇨ Enter a name for this JDK (e.g. jdk11, jdk21-azul): ")
		fmt.Print(prompt + " ")

		var name string
		fmt.Scanln(&name)

		if name == "" {
			fmt.Println(style.Input.Render("↪ Skipping unnamed JDK"))
			skipCount++
			continue
		}

		if err := java.AddJDK(name, jdk.Path); err != nil {
			fmt.Printf("%s: %v\n",
				style.Error.Render("✖ Failed to add JDK"),
				style.Error.Render(err.Error()))
		} else {
			fmt.Printf("%s: %s → %s\n\n",
				style.Success.Render("✔ Added JDK"),
				style.Success.Render(name),
				style.Path.Render(jdk.Path))
			successCount++
		}
	}

	// 显示最终统计信息
	summary := fmt.Sprintf("\n%s\n%s: %d\n%s: %d\n%s: %d\n%s: %s",
		style.Header.Render("✅ Scan Complete!"),
		style.Name.Render("New JDKs Found"), len(result.JDKs),
		style.Success.Render("Successfully Added"), successCount,
		style.Error.Render("Skipped"), skipCount,
		style.Name.Render("Total Time"), style.Success.Render(result.Duration.String()))

	fmt.Println(summary)
}
