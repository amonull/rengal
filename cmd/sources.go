package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"

	"github.com/amonull/rengal/color"
	"github.com/amonull/rengal/constant"
	"github.com/amonull/rengal/key"
	"github.com/amonull/rengal/tui"
	"github.com/amonull/rengal/util"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/icon"
	"github.com/amonull/rengal/provider"
	"github.com/amonull/rengal/style"
	"github.com/amonull/rengal/where"
)

func init() {
	rootCmd.AddCommand(sourcesCmd)
}

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Manage sources",
}

func init() {
	sourcesCmd.AddCommand(sourcesListCmd)

	sourcesListCmd.Flags().BoolP("raw", "r", false, "do not print headers")
	sourcesListCmd.Flags().BoolP("custom", "c", false, "show only custom sources")
	sourcesListCmd.Flags().BoolP("builtin", "b", false, "show only builtin sources")

	sourcesListCmd.MarkFlagsMutuallyExclusive("custom", "builtin")
	sourcesListCmd.SetOut(os.Stdout)
}

var sourcesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List an available sources",
	Run: func(cmd *cobra.Command, args []string) {
		printHeader := !lo.Must(cmd.Flags().GetBool("raw"))
		headerStyle := style.New().Foreground(color.HiBlue).Bold(true).Render
		h := func(s string) {
			if printHeader {
				cmd.Println(headerStyle(s))
			}
		}

		printBuiltin := func() {
			h("Builtin:")
			for _, p := range provider.Builtins() {
				cmd.Println(p.Name)
			}
		}

		printCustom := func() {
			h("Custom:")
			for _, p := range provider.Customs() {
				cmd.Println(p.Name)
			}
		}

		switch {
		case lo.Must(cmd.Flags().GetBool("builtin")):
			printBuiltin()
		case lo.Must(cmd.Flags().GetBool("custom")):
			printCustom()
		default:
			printBuiltin()
			if printHeader {
				cmd.Println()
			}
			printCustom()
		}
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesRemoveCmd)

	sourcesRemoveCmd.Flags().StringArrayP("name", "n", []string{}, "name of the source to remove")
	lo.Must0(
		sourcesRemoveCmd.RegisterFlagCompletionFunc(
			"name",
			func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				sources, err := filesystem.Api().ReadDir(where.Sources())
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				return lo.FilterMap(sources, func(item os.FileInfo, _ int) (string, bool) {
					name := item.Name()
					if !strings.HasSuffix(name, provider.CustomProviderExtension) {
						return "", false
					}

					return util.FileStem(filepath.Base(name)), true
				}), cobra.ShellCompDirectiveNoFileComp
			},
		),
	)
}

var sourcesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a custom source",
	Run: func(cmd *cobra.Command, args []string) {
		for _, name := range lo.Must(cmd.Flags().GetStringArray("name")) {
			path := filepath.Join(where.Sources(), name+provider.CustomProviderExtension)
			handleErr(filesystem.Api().Remove(path))
			fmt.Printf("%s successfully removed %s\n", icon.Get(icon.Success), style.Fg(color.Yellow)(name))
		}
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesInstallCmd)
}

var sourcesInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Browse and install custom scrapers",
	Long: `Browse and install custom scrapers from official GitHub repo.
https://github.com/metafates/mangal-scrapers`,
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(tui.Run(&tui.Options{Install: true}))
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesGenCmd)

	sourcesGenCmd.Flags().StringP("name", "n", "", "name of the source")
	sourcesGenCmd.Flags().StringP("url", "u", "", "url of the website")

	lo.Must0(sourcesGenCmd.MarkFlagRequired("name"))
	lo.Must0(sourcesGenCmd.MarkFlagRequired("url"))
}

var sourcesGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new lua source",
	Long:  `Generate a new lua source.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)

		author := viper.GetString(key.GenAuthor)
		if author == "" {
			usr, err := user.Current()
			if err == nil {
				author = usr.Username
			} else {
				author = "Anonymous"
			}
		}

		s := struct {
			Name            string
			URL             string
			SearchMangaFn   string
			MangaChaptersFn string
			ChapterPagesFn  string
			Author          string
		}{
			Name:            lo.Must(cmd.Flags().GetString("name")),
			URL:             lo.Must(cmd.Flags().GetString("url")),
			SearchMangaFn:   constant.SearchMangaFn,
			MangaChaptersFn: constant.MangaChaptersFn,
			ChapterPagesFn:  constant.ChapterPagesFn,
			Author:          author,
		}

		funcMap := template.FuncMap{
			"repeat": strings.Repeat,
			"plus":   func(a, b int) int { return a + b },
			"max":    util.Max[int],
		}

		tmpl, err := template.New("source").Funcs(funcMap).Parse(constant.SourceTemplate)
		handleErr(err)

		target := filepath.Join(where.Sources(), util.SanitizeFilename(s.Name)+".lua")
		f, err := filesystem.Api().Create(target)
		handleErr(err)

		defer util.Ignore(f.Close)

		err = tmpl.Execute(f, s)
		handleErr(err)

		cmd.Println(target)
	},
}
