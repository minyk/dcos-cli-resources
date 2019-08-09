package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-resources/client"
)

type Resources struct {
	PrefixCb func() string
}

func NewResources() *Resources {
	return &Resources{}
}

func (q *Resources) ReserveResource(agentid string, role string, principal string, cpus float64, mem float64) error {

	resources := []mesos.Resource{}
	resources = append(resources, resourceCPU(role, principal, cpus))
	resources = append(resources, resourceMEM(role, principal, mem))

	body := master.Call{
		Type: master.Call_RESERVE_RESOURCES,
		ReserveResources: &master.Call_ReserveResources{
			AgentID:   mesos.AgentID{Value: agentid},
			Resources: resources,
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Reservation is successful.")
	}

	return nil
}

func (q *Resources) UnreserveResource(agentid string, role string, principal string, cpus float64, cpusLabel string, mem float64, memLabel string) error {

	resources := []mesos.Resource{}
	if cpus > 0 {
		resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpusLabel))
	}
	if mem > 0 {
		resources = append(resources, resourceWithLabel("mem", role, principal, mem, memLabel))
	}

	body := master.Call{
		Type: master.Call_UNRESERVE_RESOURCES,
		UnreserveResources: &master.Call_UnreserveResources{
			AgentID:   mesos.AgentID{Value: agentid},
			Resources: resources,
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Unreservation is successful.")
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

func resourceWithLabel(resourceType string, role string, principal string, cpus float64, resourceid string) mesos.Resource {

	label := mesos.Label{
		Key:   "resource_id",
		Value: &resourceid,
	}

	labels := []mesos.Label{}
	labels = append(labels, label)
	mesosLabels := mesos.Labels{Labels: labels}

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	return mesos.Resource{
		Type:        &scala,
		Name:        resourceType,
		Role:        &role,
		Reservation: &reservation,
		Scalar:      &mesos.Value_Scalar{Value: cpus},
	}
}
