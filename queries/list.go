package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	agentcalls "github.com/mesos/mesos-go/api/v1/lib/agent/calls"
	"github.com/minyk/dcos-resources/client"
)

type ResourceList struct {
	PrefixMesosMasterApiV1 func() string
	PrefixMesosSlaveApiV0  func(string) string
	PrefixMesosSlaveApiV1  func(string) string
}

func NewResourceList() *ResourceList {
	return &ResourceList{
		PrefixMesosMasterApiV1: func() string { return "/mesos/api/v1/" },
		PrefixMesosSlaveApiV0:  func(agentid string) string { return "/agent/" + agentid },
		PrefixMesosSlaveApiV1:  func(agentid string) string { return "/agent/" + agentid + "/api/v1" },
	}
}

func (q *ResourceList) ListResourcesFromNode(agentid string, role string) error {

	resources, err := getResourcesOnRole(q.PrefixMesosSlaveApiV0(agentid), role, "")
	if err != nil {
		return err
	}

	client.PrintMessage("Role\t\tPrincipal\t\tFrameworkID\t\tType\t\tValue\t\tID\t\tPersistentID\t\tContainerPath")
	for i := range resources {
		resource := resources[i]
		rid, fid := getIDsFromLabels(resource.GetReservation().GetLabels().GetLabels())
		if resource.GetName() == "disk" {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s\t\t%s\t\t%s", resource.GetRole(), resource.GetReservation().GetPrincipal(), fid, resource.GetName(), resource.GetScalar().GetValue(), rid, resource.GetDisk().GetPersistence().GetID(), resource.GetDisk().GetVolume().GetContainerPath())
		} else if resource.GetName() == "ports" {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", resource.GetRole(), resource.GetReservation().GetPrincipal(), fid, resource.GetName(), resource.GetRanges().GoString(), rid)
		} else {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", resource.GetRole(), resource.GetReservation().GetPrincipal(), fid, resource.GetName(), resource.GetScalar().GetValue(), rid)
		}
	}

	getResourceOnExecutors(q.PrefixMesosSlaveApiV1(agentid), role)

	return nil
}

func getResourceOnExecutors(urlPath string, role string) ([]agent.Response_GetExecutors_Executor, error) {
	allExec, err := getExecutors(urlPath)
	if err != nil {
		return nil, err
	}

	client.PrintMessage("ExecutorID\t\tRole\t\tPrincipal\t\tFrameworkID\t\tType\t\tValue\t\tID\t\tPersistentID\t\tContainerPath")
	for _, exec := range allExec {
		execInfo := exec.GetExecutorInfo()
		for _, r := range execInfo.GetResources() {
			if r.GetAllocationInfo().GetRole() == role && len(r.GetReservations()) > 0 {
				rid, fid := getIDsFromLabels(r.GetReservations()[0].GetLabels().GetLabels())
				if r.GetName() == "disk" {
					client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s\t\t%s\t\t%s", execInfo.GetExecutorID(), r.GetRole(), r.GetReservations()[0].GetPrincipal(), fid, r.GetName(), r.GetScalar().GetValue(), rid, r.GetDisk().GetPersistence().GetID(), r.GetDisk().GetVolume().GetContainerPath())
				} else if r.GetName() == "ports" {
					client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", execInfo.GetExecutorID(), r.GetRole(), r.GetReservations()[0].GetPrincipal(), fid, r.GetName(), r.GetRanges().GoString(), rid)
				} else {
					client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", execInfo.GetExecutorID(), r.GetRole(), r.GetReservations()[0].GetPrincipal(), fid, r.GetName(), r.GetScalar().GetValue(), rid)
				}
			} else if r.GetAllocationInfo().GetRole() == role && len(r.GetReservations()) <= 0 {
				client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t", execInfo.GetExecutorID(), r.GetRole(), "", "", r.GetName(), r.GetScalar().GetValue())
			}
		}
	}

	return nil, nil
}

func getExecutors(urlPath string) ([]agent.Response_GetExecutors_Executor, error) {
	body := agentcalls.GetExecutors()
	requestContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := client.HTTPServicePostJSON(urlPath, requestContent)
	if err != nil {
		return nil, err
	}

	executors := agent.Response{}
	err = json.Unmarshal(resp, &executors)
	if err != nil {
		return nil, err
	}

	return executors.GetExecutors.Executors, nil
}
