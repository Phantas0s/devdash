package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/Phantas0s/devdash/internal"
	"github.com/spf13/cobra"
)

var configType string

func generateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate default dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stdout, generate(args))
		},
	}

	generateCmd.Flags().StringVarP(&configType, "type", "t", "", "Debug Mode - doesn't display graph")
	return generateCmd
}

func generate(args []string) string {
	switch configType {
	case "ga":
		return createBlogDefaultConfig()
	default:
		return createBlogDefaultConfig()
	}

	return "oops"
}

func createBlogDefaultConfig() string {
	wizardIntro("blogs")
	keyfile := askKeyfile()
	address := askSiteAddress()
	viewID := askViewID()

	ut, err := template.New("Ga").Parse(internal.GA())
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer([]byte{})
	err = ut.Execute(b, internal.CreateBlogConfig(keyfile, address, viewID))

	return b.String()
}

func wizardIntro(name string) {
	fmt.Fprintf(os.Stdout, "You're now generating the dashboard for %s.\n", name)
	fmt.Fprintln(os.Stdout, "You can let all of these fields blank and fill them later in the file itself.")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "------------------------------------------------------------------------------")
}

func askKeyfile() string {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "The Keyfile is required to connect to Google Analytics and Google Search Console")
	fmt.Fprintln(os.Stdout, "See: https://thedevdash.com/reference/services/google-analytics/")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprint(os.Stdout, "Enter the keyfile path for Google Analytics and Google Search Console: ")

	var keyfile string
	fmt.Scanf("%s", &keyfile)

	return keyfile
}

func askSiteAddress() string {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprint(os.Stdout, "Enter the address of your blog (beginning with http or https): ")

	var address string
	fmt.Scanf("%s", &address)

	return address
}

func askViewID() string {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "The View ID is required for Google Analytics.")
	fmt.Fprintln(os.Stdout, "See https://thedevdash.com/reference/services/google-analytics/")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprint(os.Stdout, "Enter the view ID for Google Analytics: ")

	var viewID string
	fmt.Scanf("%s", &viewID)

	return viewID
}
