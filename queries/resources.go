package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib"
	mastercalls "github.com/mesos/mesos-go/api/v1/lib/master/calls"
	"github.com/minyk/dcos-resources/client"
)

type Resources struct {
	PrefixMesosMasterApiV1 func() string
	PrefixMesosSlaveApiV0  func(string) string
	PrefixMesosSlaveApiV1  func(string) string
}

func NewResources() *Resources {
	return &Resources{
		PrefixMesosMasterApiV1: func() string { return "/mesos/api/v1/" },
		PrefixMesosSlaveApiV0:  func(agentid string) string { return "/agent/" + agentid },
		PrefixMesosSlaveApiV1:  func(agentid string) string { return "/agent/" + agentid + "/api/v1" },
	}
}

func (q *Resources) ReserveResource(agentid string, role string, principal string, cpus float64, mem float64) error {

	var resources []mesos.Resource
	resources = append(resources, resourceCPU(role, principal, cpus))
	resources = append(resources, resourceMEM(role, principal, mem))

	body := mastercalls.ReserveResources(mesos.AgentID{Value: agentid}, resources...)

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Reservation is successful.")
	}

	return nil
}

func resourceCPU(role string, principal string, cpus float64) mesos.Resource {
	return resource("cpus", role, principal, cpus)
}

func resourceMEM(role string, principal string, mem float64) mesos.Resource {
	return resource("mem", role, principal, mem)
}

func resource(resourceType string, role string, principal string, cpus float64) mesos.Resource {

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
	}

	return mesos.Resource{
		Type:        &scala,
		Name:        resourceType,
		Role:        &role,
		Reservation: &reservation,
		Scalar:      &mesos.Value_Scalar{Value: cpus},
	}
}

//func getResourcesOnRole(urlPath string, role string) (ResourceRole, error) {
//	resourcesFull, err := listResources(urlPath)
//	if err != nil {
//		return nil, err
//	}
//
//	resources := resourcesFull[role]
//	if len(resources) == 0 {
//		return nil, errors.New("no resources are reserved for role")
//	}
//
//	return resources, nil
//}
