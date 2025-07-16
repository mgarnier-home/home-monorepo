package cmd

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"mgarnier11.fr/go/orchestrator-cli/config"
	"mgarnier11.fr/go/orchestrator-cli/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Action struct {
	Title        string
	Icon         string
	Shell        string
	PopupOnStart string
	Timeout      int
}

type Content struct {
	Title    string     `yaml:"title,omitempty"`
	Type     string     `yaml:"type,omitempty"`
	Contents []*Content `yaml:"contents,omitempty"`
}

type Dashboard struct {
	Title    string
	Contents []*Content
}

type oliveTinConfig struct {
	Actions    []*Action
	Dashboards []*Dashboard
}

func createActionCommands(stack string, host string) []*Action {
	actions := []*Action{}

	for _, actionStr := range utils.ActionList {

		shell := "home-cli"

		if stack != "all" {
			shell += " " + stack
		}
		if host != "all" {
			shell += " " + host
		}
		shell += " " + actionStr

		actions = append(actions, &Action{
			Title:        stack + " " + host + " " + actionStr,
			Icon:         "box",
			Shell:        shell,
			PopupOnStart: "execution-dialog",
			Timeout:      150,
		})
	}

	return actions
}

func getActions(compose *utils.ComposeConfig) (actions []*Action, stacksActions map[string][]*Action, hostsActions map[string][]*Action) {
	actions = []*Action{}

	for stack, hosts := range compose.Stacks {
		if len(hosts) > 1 {
			actions = append(actions, createActionCommands(stack, "all")...)
		}

		for _, host := range hosts {
			actions = append(actions, createActionCommands(stack, host)...)
		}
	}

	for host, stacks := range compose.Hosts {
		actions = append(actions, createActionCommands("all", host)...)

		for _, stack := range stacks {
			actions = append(actions, createActionCommands(stack, host)...)
		}
	}

	actions = append(actions, createActionCommands("all", "all")...)

	slices.SortFunc(actions, func(a, b *Action) int {
		return strings.Compare(a.Title, b.Title)
	})

	actions = slices.CompactFunc(actions, func(a *Action, b *Action) bool {
		return a.Title == b.Title
	})

	stacksActions = make(map[string][]*Action)
	hostsActions = make(map[string][]*Action)

	for _, action := range actions {
		for _, stack := range utils.StackList {
			if strings.Contains(action.Title, stack) {
				stacksActions[stack] = append(stacksActions[stack], action)
			}
		}

		for _, host := range utils.HostList {
			if strings.Contains(action.Title, host) {
				hostsActions[host] = append(hostsActions[host], action)
			}
		}
	}

	return actions, stacksActions, hostsActions
}

func printAction(action *Action) {
	fmt.Println(action.Title + " => " + action.Shell)
}

func printActions(actions []*Action) {
	for _, action := range actions {
		printAction(action)
	}
}

func getContents(action string, actions []*Action) []*Content {
	contents := []*Content{
		{
			Type:  "display",
			Title: action,
		},
	}

	filteredActions := utils.FindAll(actions, func(a *Action) bool {
		return strings.Contains(a.Title, action)
	})

	for _, a := range filteredActions {
		contents = append(contents, &Content{
			Title: a.Title,
		})
	}

	return contents
}

func saveConfig(compose *utils.ComposeConfig) {
	actions, stacksActions, hostsActions := getActions(compose)

	fmt.Println("===========Actions===========")
	for _, action := range actions {
		fmt.Println(action.Title + " => " + action.Shell)
	}

	for stack, actions := range stacksActions {
		fmt.Println("*****************Stack", stack)
		printActions(actions)
	}

	for host, actions := range hostsActions {
		fmt.Println("*****************Host", host)
		printActions(actions)
	}

	os.Mkdir(config.Env.OliveTinConfigDir, 0755)
	file, err := os.OpenFile(path.Join(config.Env.OliveTinConfigDir, "config.yaml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		panic("File error while saving config")
	}
	defer file.Close()

	dashboardContents := []*Content{}

	for host, actions := range hostsActions {
		for i, action := range utils.ActionList {
			if i == 0 {
				dashboardContents = append(dashboardContents, &Content{
					Title:    host,
					Type:     "fieldset",
					Contents: getContents(action, actions),
				})
			} else {
				dashboardContents = append(dashboardContents, &Content{
					Type:     "fieldset",
					Contents: getContents(action, actions),
				})
			}
		}
	}

	err = yaml.NewEncoder(file).Encode(oliveTinConfig{
		Actions: actions,
		Dashboards: []*Dashboard{
			{
				Title:    "Dashboard",
				Contents: dashboardContents,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		panic("Yaml while saving config")
	}

	fmt.Println("Config saved successfully")
}

var buildOlivetinCmd = &cobra.Command{
	Use:   "build-olivetin",
	Short: "Build olivetin config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building olivetin config")

		saveConfig(utils.GetCompose())

	},
}

func init() {
	buildOlivetinCmd.Hidden = true
	rootCmd.AddCommand(buildOlivetinCmd)
}
