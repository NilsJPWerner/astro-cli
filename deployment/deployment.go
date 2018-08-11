package deployment

import (
	"fmt"

	"github.com/astronomerio/astro-cli/config"
	"github.com/astronomerio/astro-cli/houston"
	"github.com/astronomerio/astro-cli/messages"
	"github.com/astronomerio/astro-cli/pkg/httputil"
	"github.com/astronomerio/astro-cli/pkg/jsonstr"
)

var (
	http = httputil.NewHTTPClient()
	api  = houston.NewHoustonClient(http)
)

func Create(label, ws string) error {
	deployment, err := api.CreateDeployment(label, ws)
	if err != nil {
		return err
	}

	fmt.Printf(messages.HOUSTON_DEPLOYMENT_CREATE_SUCCESS, deployment.Id)

	fmt.Printf("\n"+messages.EE_LINK_AIRFLOW+"\n", deployment.ReleaseName, config.CFG.CloudDomain.GetString())
	fmt.Printf(messages.EE_LINK_FLOWER+"\n", deployment.ReleaseName, config.CFG.CloudDomain.GetString())

	return nil
}

func Delete(uuid string) error {
	resp, err := api.DeleteDeployment(uuid)
	if err != nil {
		return err
	}

	fmt.Printf(messages.HOUSTON_DEPLOYMENT_DELETE_SUCCESS, resp.Id)

	return nil
}

// List all airflow deployments
func List(ws string) error {
	r := "  %-30s %-50s"
	// colorFmt := "\033[33;m"
	// colorTrm := "\033[0m"

	deployments, err := api.GetDeployments(ws)
	if err != nil {
		return err
	}

	h := fmt.Sprintf(r, "NAME", "UUID")
	fmt.Println(h)

	for _, d := range deployments {
		fullStr := fmt.Sprintf(r, d.Label, d.Id)
		fmt.Println(fullStr)
	}
	return nil
}

// Update an airflow deployment
func Update(deploymentId string, args map[string]string) error {
	s := jsonstr.MapToJsonObjStr(args)

	dep, err := api.UpdateDeployment(deploymentId, s)
	if err != nil {
		return err
	}

	fmt.Printf(messages.HOUSTON_DEPLOYMENT_UPDATE_SUCCESS, dep.Id)

	return nil
}
