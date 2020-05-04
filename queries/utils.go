package queries

import (
	"encoding/json"
	"errors"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	mastercalls "github.com/mesos/mesos-go/api/v1/lib/master/calls"
	"github.com/minyk/dcos-resources/client"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Struct for Mesos API V0
type AgentState struct {
	AgentReservedResourcesFull ReservedResourcesFull `json:"reserved_resources_full,omitempty"`
}

type ReservedResourcesFull map[string]ResourceRole

type ResourceRole []mesos.Resource

func getAgentList(masterUrl string) ([]string, error) {
	body := mastercalls.GetAgents()

	requestContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	responseContent, err := client.HTTPServicePostJSON(masterUrl, requestContent)
	if err != nil {
		return nil, err
	}

	agents := master.Response{}
	err = json.Unmarshal(responseContent, &agents)
	if err != nil {
		return nil, err
	}

	var list []string

	for _, agent := range agents.GetAgents.GetAgents() {
		list = append(list, agent.AgentInfo.ID.Value)
	}

	return list, nil
}

func getResourcesOnRole(urlPath string, role string, principal string) (ResourceRole, error) {
	resourcesFull, err := listResources(urlPath)
	if err != nil {
		return nil, err
	}

	resources := resourcesFull[role]
	if len(resources) == 0 {
		return nil, errors.New("no resources are reserved for role")
	}

	if principal == "" {
		return resources, nil
	}

	var resourcesOfPrinciapl ResourceRole
	for _, r := range resources {
		if r.GetReservation().GetPrincipal() == principal {
			resourcesOfPrinciapl = append(resourcesOfPrinciapl, r)
		}
	}

	return resourcesOfPrinciapl, nil
}
