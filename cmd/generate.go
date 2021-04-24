package cmd

// TODO need to test all of that!

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/Phantas0s/devdash/internal"
	"github.com/spf13/cobra"
)

var configType, templateFile string

func generateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate dashboard templates",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stdout, generate(args))
		},
	}

	generateCmd.Flags().StringVarP(&configType, "type", "t", "", "Debug Mode - doesn't display graph")
	generateCmd.Flags().StringVarP(&templateFile, "file", "f", "", "Save the template into a file ($XDG_CONFIG_HOME/<name_of_file>")
	return generateCmd
}

func generate(args []string) string {
	switch configType {
	case "blog":
		return createBlogDefaultConfig()
	case "project", "githubProject", "github_project":
		return createGitHubDefaultConfig()
	default:
		return createBlogDefaultConfig()
	}

	return "oops"
}

func createBlogDefaultConfig() string {
	wizardIntro("a blog")
	keyfile := askKeyfile()
	address := askSiteAddress()
	viewID := askViewID()

	ut, err := template.New("blog").Parse(internal.Blog())
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer([]byte{})
	err = ut.Execute(b, internal.CreateBlogConfig(keyfile, address, viewID))

	return createOrDisplayTemplate(b.String())
}

func createGitHubDefaultConfig() string {
	wizardIntro("a project on GitHub")
	token := askGitHubToken()
	owner := askGitHubOwner()
	repo := askGitHubRepo()

	ut, err := template.New("github").Parse(internal.GitHubProject())
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer([]byte{})
	err = ut.Execute(b, internal.CreateGitHubProjectConfig(token, owner, repo))

	return createOrDisplayTemplate(b.String())
}

func createOrDisplayTemplate(content string) string {
	if templateFile != "" {
		cfg := createTemplateFile(templateFile, content)
		return fmt.Sprintf("The file %s has been created.", cfg)
	} else {
		return content
	}
}

func createTemplateFile(filename string, template string) string {
	return createConfig(dashPath(), filename+".yml", template)
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

func askGitHubToken() string {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "A token is required for GitHub.")
	fmt.Fprintln(os.Stdout, "See https://thedevdash.com/reference/services/github/")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprint(os.Stdout, "Enter the GitHub token: ")

	var token string
	fmt.Scanf("%s", &token)

	return token
}

func askGitHubOwner() string {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprint(os.Stdout, "Enter your GitHub username: ")

	var owner string
	fmt.Scanf("%s", &owner)

	return owner
}

func askGitHubRepo() string {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprint(os.Stdout, "Enter the name of the GitHub repo: ")

	var repo string
	fmt.Scanf("%s", &repo)

	return repo
}
