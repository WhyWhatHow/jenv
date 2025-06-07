package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whywhathow/jenv/internal/java"
	"github.com/whywhathow/jenv/internal/style"
	"strings"
)

var (
	scanCmd = &cobra.Command{
		Aliases: []string{"sc"},
		Use:     "scan <dir>",
		Short:   "Scan a directory for JDKs (max depth: 3 subdirectories)",
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

		//TODO [whywhathow] [2025/3/13] [opt]  如果以后支持多平台,example 是不是需要根据os不同进行适配呢?
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

	jdks := java.ScanJDK(dir)
	successCount := 0
	skipCount := 0

	for i, jdk := range jdks {
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

	// 显示统计信息
	summary := fmt.Sprintf("\n%s\n%s: %d\n%s: %d\n%s: %d",
		style.Header.Render("Scan Complete!"),
		style.Name.Render("Total Found"), len(jdks),
		style.Success.Render("Successfully Added"), successCount,
		style.Error.Render("Skipped"), skipCount)

	fmt.Println(summary)
}
