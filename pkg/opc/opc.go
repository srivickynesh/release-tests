package opc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

type PacInfoInstall struct {
	PipelinesAsCode   PipelinesAsCodeSection
	GithubApplication GithubApplicationSection
	RepositoriesCR    RepositoriesCRSection
}

type PipelinesAsCodeSection struct {
	InstallVersion   string
	InstallNamespace string
}

type GithubApplicationSection struct {
	Name               string
	URL                string
	HomePage           string
	Description        string
	Created            string
	InstallationsCount string
	WebhookURL         string
}

type RepositoriesCRSection struct {
	Count        int
	Repositories []Repository
}

type Repository struct {
	Namespace string
	URL       string
}

type PipelineRunList struct {
	Name     string
	Started  string
	Duration string
	Status   string
}

func GetOpcPacInfoInstall() (*PacInfoInstall, error) {
	result := cmd.MustSucceed("opc", "pac", "info", "install")
	output := result.Stdout()
	lines := strings.Split(output, "\n")

	var pacInfo PacInfoInstall
	section := "" // current section: "pipelines", "github", or "repositories"
	tableHeaderParsed := false

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		switch line {
		case "Pipelines as Code:":
			section = "pipelines"
			continue
		case "Github Application:":
			section = "github"
			continue
		}

		if strings.HasPrefix(line, "Repositories CR:") {
			section = "repositories"
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				countStr := strings.TrimSpace(parts[1])
				count, err := strconv.Atoi(countStr)
				if err != nil {
					return nil, fmt.Errorf("failed to parse repository count: %v", err)
				}
				pacInfo.RepositoriesCR.Count = count
			}
			continue
		}

		if section == "pipelines" || section == "github" {
			if !strings.Contains(line, ":") {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if section == "pipelines" {
				switch key {
				case "Install Version":
					pacInfo.PipelinesAsCode.InstallVersion = value
				case "Install Namespace":
					pacInfo.PipelinesAsCode.InstallNamespace = value
				}
			} else if section == "github" {
				switch key {
				case "Name":
					pacInfo.GithubApplication.Name = value
				case "URL":
					pacInfo.GithubApplication.URL = value
				case "HomePage":
					pacInfo.GithubApplication.HomePage = value
				case "Description":
					pacInfo.GithubApplication.Description = value
				case "Created":
					pacInfo.GithubApplication.Created = value
				case "Installations Count":
					pacInfo.GithubApplication.InstallationsCount = value
				case "Webhook URL":
					pacInfo.GithubApplication.WebhookURL = value
				}
			}
			continue
		}

		if section == "repositories" {
			if !tableHeaderParsed && strings.Contains(line, "Namespace") && strings.Contains(line, "URL") {
				tableHeaderParsed = true
				continue
			}
			if strings.HasPrefix(line, "- ") {
				line = strings.TrimPrefix(line, "-")
				line = strings.TrimSpace(line)
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			repo := Repository{
				Namespace: fields[0],
				URL:       fields[1],
			}
			pacInfo.RepositoriesCR.Repositories = append(pacInfo.RepositoriesCR.Repositories, repo)
		}
	}

	if pacInfo.PipelinesAsCode.InstallVersion == "" {
		return nil, fmt.Errorf("output of 'opc pac info install' is empty or missing Pipelines as Code information")
	}

	return &pacInfo, nil
}

func GetOpcPrList() ([]PipelineRunList, error) {
	result := cmd.MustSucceed("opc", "pr", "ls")
	output := strings.TrimSpace(result.Stdout())
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected output: %s", output)
	}

	var runs []PipelineRunList
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			return nil, fmt.Errorf("unexpected row format: %s", line)
		}

		nameFieldCount := len(fields) - 3
		run := PipelineRunList{
			Name:     strings.Join(fields[:nameFieldCount], " "),
			Started:  fields[nameFieldCount],
			Duration: fields[nameFieldCount+1],
			Status:   fields[nameFieldCount+2],
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func GetOpcPr(pipelineRunName string) (*PipelineRun, error) {
	runs, err := GetOpcPrList()
	if err != nil {
		return nil, err
	}
	for _, run := range runs {
		if run.Name == pipelineRunName {
			return &run, nil
		}
	}
	return nil, fmt.Errorf("pipeline run %q not found", pipelineRunName)
}

// func GetOpcClusterTriggerBinding() ([]string, error) {
// 	cmd := exec.Command(opcPath, "clustertriggerbinding")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("opc clustertriggerbinding failed: %s", string(output))
// 	}
// 	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results, nil
// }

// func StartPipelineUseDefaults(pipelineName string) ([]string, error) {
// 	cmd := exec.Command(opcPath, "pipeline", "start", "--use-defaults", pipelineName)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("opc pipeline start --use-defaults failed: %s", string(output))
// 	}
// 	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results, nil
// }

// func StartPipelineWithWorkspace(pipelineName, workspace string) ([]string, error) {
// 	info, err := os.Stat(workspace)
// 	if err != nil || !info.IsDir() {
// 		return nil, fmt.Errorf("invalid workspace directory: %s", workspace)
// 	}
// 	cmd := exec.Command(opcPath, "pipeline", "start", pipelineName)
// 	cmd.Dir = workspace
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("opc pipeline start (workspace) failed: %s", string(output))
// 	}
// 	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results, nil
// }

// func GetOpcEventListenerList() ([]string, error) {
// 	cmd := exec.Command(opcPath, "eventlistener", "list")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("opc eventlistener list failed: %s", string(output))
// 	}
// 	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results, nil
// }

// // GetOpcHubSearch runs "opc hub search <query>" and returns its output.
// func GetOpcHubSearch(query string) ([]string, error) {
// 	cmd := exec.Command(opcPath, "hub", "search", query)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("opc hub search failed: %s", string(output))
// 	}
// 	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results, nil
// }

// func GetOpcTriggerTemplateList() ([]string, error) {
// 	cmd := exec.Command(opcPath, "triggertemplate", "ls")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("opc triggertemplate ls failed: %s", string(output))
// 	}
// 	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results, nil
// }

// ---------------------------

// func splitNonEmptyLines(text string) []string {
// 	lines := strings.Split(text, "\n")
// 	var results []string
// 	for _, line := range lines {
// 		if trimmed := strings.TrimSpace(line); trimmed != "" {
// 			results = append(results, trimmed)
// 		}
// 	}
// 	return results
// }

// // parseKeyValuePairs converts lines of the form "Key: Value" into a dictionary.
// func parseKeyValuePairs(text string) map[string]string {
// 	infoMap := make(map[string]string)
// 	lines := strings.Split(text, "\n")
// 	for _, line := range lines {
// 		trimmed := strings.TrimSpace(line)
// 		if trimmed == "" {
// 			continue
// 		}
// 		parts := strings.SplitN(trimmed, ":", 2)
// 		if len(parts) == 2 {
// 			key := strings.TrimSpace(parts[0])
// 			value := strings.TrimSpace(parts[1])
// 			infoMap[key] = value
// 		} else {
// 			infoMap[trimmed] = ""
// 		}
// 	}
// 	return infoMap
// }

// // GetOpcClusterTriggerBinding runs "opc clustertriggerbinding" and returns its output as a slice of strings.
// func GetOpcClusterTriggerBinding() ([]string, error) {
// 	result := cmd.MustSucceed(opcPath, "clustertriggerbinding")
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// // StartPipelineUseDefaults runs "opc pipeline start --use-defaults <pipelineName>"
// // and returns its output as a slice of strings.
// func StartPipelineUseDefaults(pipelineName string) ([]string, error) {
// 	result := cmd.MustSucceed(opcPath, "pipeline", "start", "--use-defaults", pipelineName)
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// // StartPipelineWithWorkspace runs "opc pipeline start <pipelineName>" in the specified workspace directory.
// // The workspace must exist (e.g. an empty directory). Returns the output as a slice of strings.
// func StartPipelineWithWorkspace(pipelineName, workspace string) ([]string, error) {
// 	// Verify the workspace directory exists.
// 	info, err := os.Stat(workspace)
// 	if err != nil || !info.IsDir() {
// 		return nil, fmt.Errorf("invalid workspace directory: %s", workspace)
// 	}

// 	// Change working directory temporarily.
// 	originalDir, err := os.Getwd()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get current directory: %v", err)
// 	}
// 	if err := os.Chdir(workspace); err != nil {
// 		return nil, fmt.Errorf("failed to change directory to workspace: %v", err)
// 	}
// 	// Ensure we change back to the original directory.
// 	defer os.Chdir(originalDir)

// 	result := cmd.MustSucceed(opcPath, "pipeline", "start", pipelineName)
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// // GetOpcEventListenerList runs "opc eventlistener list" and returns its output as a slice of strings.
// func GetOpcEventListenerList() ([]string, error) {
// 	result := cmd.MustSucceed(opcPath, "eventlistener", "list")
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// // GetOpcHubSearch runs "opc hub search <query>" and returns its output as a slice of strings.
// func GetOpcHubSearch(query string) ([]string, error) {
// 	result := cmd.MustSucceed(opcPath, "hub", "search", query)
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// // GetOpcTriggerBindingList runs "opc triggerbinding ls" and returns its output as a slice of strings.
// func GetOpcTriggerBindingList() ([]string, error) {
// 	result := cmd.MustSucceed(opcPath, "triggerbinding", "ls")
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// // GetOpcTriggerTemplateLs runs "opc triggertemplate ls" and returns its output as a slice of strings.
// func GetOpcTriggerTemplates() ([]string, error) {
// 	result := cmd.MustSucceed(opcPath, "triggertemplate", "ls")
// 	output := strings.TrimSpace(result.Stdout())
// 	return splitNonEmptyLines(output), nil
// }

// -----------------------------------------------------------

// ----------------------------------------
