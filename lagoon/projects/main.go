package projects

import (
	"encoding/json"
	"fmt"

	"github.com/amazeeio/lagoon-cli/api"
	"github.com/amazeeio/lagoon-cli/graphql"
	"github.com/amazeeio/lagoon-cli/output"
)

// ListAllProjects will list all projects
func ListAllProjects() ([]byte, error) {
	// set up a lagoonapi client
	lagoonAPI, err := graphql.LagoonAPI()
	if err != nil {
		return []byte(""), err
	}
	allProjects, err := lagoonAPI.GetAllProjects(graphql.AllProjectsFragment)
	if err != nil {
		return []byte(""), err
	}
	returnResult, err := processAllProjects(allProjects)
	if err != nil {
		return []byte(""), err
	}
	return returnResult, nil
}

func processAllProjects(allProjects []byte) ([]byte, error) {
	var projects []api.Project
	err := json.Unmarshal([]byte(allProjects), &projects)
	if err != nil {
		return []byte(""), err
	}
	// process the data for output
	data := []output.Data{}
	for _, project := range projects {
		// count the current dev environments in a project
		var currentDevEnvironments = 0
		for _, environment := range project.Environments {
			if environment.EnvironmentType == "development" {
				currentDevEnvironments++
			}
		}
		data = append(data, []string{
			fmt.Sprintf("%v", project.ID),
			fmt.Sprintf("%v", project.Name),
			fmt.Sprintf("%v", project.GitURL),
			fmt.Sprintf("%v/%v", currentDevEnvironments, project.DevelopmentEnvironmentsLimit),
		})
	}
	dataMain := output.Table{
		Header: []string{"ID", "ProjectName", "GitURL", "DevEnvironments"},
		Data:   data,
	}
	return json.Marshal(dataMain)
}

// ListEnvironmentsForProject will list all environments for a project
func ListEnvironmentsForProject(projectName string) ([]byte, error) {
	// set up a lagoonapi client
	lagoonAPI, err := graphql.LagoonAPI()
	if err != nil {
		return []byte(""), err
	}
	// get project info from lagoon
	project := api.Project{
		Name: projectName,
	}
	projectByName, err := lagoonAPI.GetProjectByName(project, graphql.ProjectByNameFragment)
	if err != nil {
		return []byte(""), err
	}
	// @TODO do we need this data? I'm not sure
	// fmt.Println(fmt.Sprintf("%s: %s", aurora.Yellow("Project Name"), projects.Name))
	// fmt.Println(fmt.Sprintf("%s: %d", aurora.Yellow("Project ID"), projects.ID))
	// fmt.Println()
	// fmt.Println(fmt.Sprintf("%s: %s", aurora.Yellow("Git"), projects.GitURL))
	// fmt.Println(fmt.Sprintf("%s: %s", aurora.Yellow("Branches"), projects.Branches))
	// fmt.Println(fmt.Sprintf("%s: %s", aurora.Yellow("Pull Requests"), projects.Pullrequests))
	// fmt.Println(fmt.Sprintf("%s: %s", aurora.Yellow("Production Environment"), projects.ProductionEnvironment))
	// fmt.Println(fmt.Sprintf("%s: %d / %d", aurora.Yellow("Development Environments"), currentDevEnvironments, projects.DevelopmentEnvironmentsLimit))
	// fmt.Println()
	returnResult, err := processProjectInfo(projectByName)
	if err != nil {
		return []byte(""), err
	}
	return returnResult, nil
}

func processProjectInfo(projectByName []byte) ([]byte, error) {
	var projects api.Project
	err := json.Unmarshal([]byte(projectByName), &projects)
	if err != nil {
		return []byte(""), err
	}
	// count the current dev environments in a project
	var currentDevEnvironments = 0
	for _, environment := range projects.Environments {
		if environment.EnvironmentType == "development" {
			currentDevEnvironments++
		}
	}
	// process the data for output
	data := []output.Data{}
	for _, environment := range projects.Environments {
		data = append(data, []string{
			fmt.Sprintf("%d", environment.ID),
			environment.Name,
			string(environment.DeployType),
			string(environment.EnvironmentType),
			environment.Route,
			//fmt.Sprintf("ssh -p %s -t %s@%s", viper.GetString("lagoons."+cmdLagoon+".port"), environment.OpenshiftProjectName, viper.GetString("lagoons."+cmdLagoon+".hostname")),
		})
	}
	dataMain := output.Table{
		Header: []string{"ID", "Name", "DeployType", "Environment", "Route"}, //, "SSH"},
		Data:   data,
	}
	return json.Marshal(dataMain)
}

// AddProject .
func AddProject(projectName string, jsonPatch string) ([]byte, error) {
	lagoonAPI, err := graphql.LagoonAPI()
	if err != nil {
		return []byte(""), err
	}
	project := api.ProjectPatch{}
	err = json.Unmarshal([]byte(jsonPatch), &project)
	if err != nil {
		return []byte(""), err
	}
	project.Name = projectName
	projectAddResult, err := lagoonAPI.AddProject(project, graphql.ProjectByNameFragment)
	if err != nil {
		return []byte(""), err
	}
	return projectAddResult, nil
}

// DeleteProject .
func DeleteProject(projectName string) ([]byte, error) {
	lagoonAPI, err := graphql.LagoonAPI()
	if err != nil {
		return []byte(""), err
	}
	project := api.Project{
		Name: projectName,
	}
	returnResult, err := lagoonAPI.DeleteProject(project)
	return returnResult, err
}

// UpdateProject .
func UpdateProject(projectName string, jsonPatch string) ([]byte, error) {
	lagoonAPI, err := graphql.LagoonAPI()
	if err != nil {
		return []byte(""), err
	}
	// get the project id from name
	projectBName := api.Project{
		Name: projectName,
	}
	projectByName, err := lagoonAPI.GetProjectByName(projectBName, graphql.ProjectByNameFragment)
	if err != nil {
		return []byte(""), err
	}
	projectUpdate, err := processProjectUpdate(projectByName, jsonPatch)
	if err != nil {
		return []byte(""), err
	}
	returnResult, err := lagoonAPI.UpdateProject(projectUpdate, graphql.ProjectByNameFragment)
	if err != nil {
		return []byte(""), err
	}
	return returnResult, nil
}

func processProjectUpdate(projectByName []byte, jsonPatch string) (api.UpdateProject, error) {
	var projects api.Project
	var projectUpdate api.UpdateProject
	var project api.ProjectPatch
	err := json.Unmarshal([]byte(projectByName), &projects)
	if err != nil {
		return projectUpdate, err
	}
	projectID := projects.ID

	// patch the project by id
	err = json.Unmarshal([]byte(jsonPatch), &project)
	if err != nil {
		return projectUpdate, err
	}
	projectUpdate = api.UpdateProject{
		ID:    projectID,
		Patch: project,
	}
	return projectUpdate, nil
}