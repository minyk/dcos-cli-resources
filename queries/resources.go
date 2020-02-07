package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-resources/client"
)

type Resources struct {
	PrefixMesosMasterApiV1 func() string
	PrefixMesosSlaveApiV0  func(string) string
}

func NewResources() *Resources {
	return &Resources{
		PrefixMesosMasterApiV1: func() string { return "/mesos/api/v1/" },
		PrefixMesosSlaveApiV0:  func(agentid string) string { return "/agent/" + agentid },
	}
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

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Reservation is successful.")
	}

	return nil
}

func (q *Resources) UnreserveResource(agentid string, role string, principal string, cpus float64, cpusLabel string, mem float64, memLabel string, disk float64, diskLabel string) error {

	resources := []mesos.Resource{}
	if cpus > 0 {
		resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpusLabel))
	}
	if mem > 0 {
		resources = append(resources, resourceWithLabel("mem", role, principal, mem, memLabel))
	}
	if disk > 0 {
		resources = append(resources, resourceWithLabel("disk", role, principal, disk, diskLabel))
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

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Unreservation is successful.")
	}

	return nil

}

func (q *Resources) DestroyVolume(agentid string, role string, principal string, disk float64, resourceid string, persistid string, containerpath string, hostpath string) error {

	resources := []mesos.Resource{}

	resources = append(resources, resourceDiskWithLabel(role, principal, disk, resourceid, persistid, containerpath, ""))

	request_body := master.Call{
		Type: master.Call_DESTROY_VOLUMES,
		DestroyVolumes: &master.Call_DestroyVolumes{
			AgentID: mesos.AgentID{Value: agentid},
			Volumes: resources,
		},
	}

	requestContent, err := json.Marshal(request_body)
	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Persistence Volume is successfully removed.")
	}

	return nil
}

func (q *Resources) UnreserveResourceAll(agentid string, role string, principal string) error {

	// query current state from agent
	state_body := agent.Call{
		Type: agent.Call_GET_STATE,
	}
	requestContent, err := json.Marshal(state_body)
	response, err := client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	}

	agentStateReponse := master.Response{}
	json.Unmarshal(response, &agentStateReponse)

	// TODO find all resources with role and principal
	cpus := float64(4)
	cpus_resource_id := ""
	mem := float64(4)
	mem_resource_id := ""

	// unreserve cpus, mem resources with resource_id
	resources := []mesos.Resource{}
	resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpus_resource_id))
	resources = append(resources, resourceWithLabel("mem", role, principal, mem, mem_resource_id))

	body := master.Call{
		Type: master.Call_UNRESERVE_RESOURCES,
		UnreserveResources: &master.Call_UnreserveResources{
			AgentID:   mesos.AgentID{Value: agentid},
			Resources: resources,
		},
	}

	requestContent, err = json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
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

func resourceDiskWithLabel(role string, principal string, disk float64, resourceid string, persistid string, containerPath string, hostPath string) mesos.Resource {

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

	persist := mesos.Resource_DiskInfo_Persistence{
		ID:        persistid,
		Principal: &principal,
	}

	rw := mesos.RW

	volume := mesos.Volume{
		Mode:          &rw,
		ContainerPath: containerPath,
	}

	diskinfo := mesos.Resource_DiskInfo{
		Persistence: &persist,
		Volume:      &volume,
		Source:      nil,
	}

	return mesos.Resource{
		Type:        &scala,
		Name:        "disk",
		Role:        &role,
		Reservation: &reservation,
		Scalar:      &mesos.Value_Scalar{Value: disk},
		Disk:        &diskinfo,
	}
}

func (q *Resources) ListResourcesFromNode(agentid string, role string) error {

	response, err := client.HTTPServiceGet(q.PrefixMesosSlaveApiV0(agentid) + "/state")
	if err != nil {
		return err
	}

	agentStateReponse := AgentState{}
	json.Unmarshal(response, &agentStateReponse)

	resources := agentStateReponse.AgentReservedResourcesFull[role]
	if len(resources) == 0 {
		client.PrintMessage("No resources are reserved for %s", role)
		return nil
	}

	client.PrintMessage("Role\t\tPrincipal\t\tName\t\tValue\t\tID\t\tPersistentID\t\tContainerPath")
	for i := range resources {
		resource := resources[i]

		if resource.Name == "disk" {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%f\t\t%s\t\t%s\t\t%s", resource.Role, resource.Reservation.Principal, resource.Name, resource.Scalar, resource.Reservation.Labels.Labels[0].Value, resource.Disk.Persistence.ID, resource.Disk.Volume.ContainerPath)
		} else {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%f\t\t%s", resource.Role, resource.Reservation.Principal, resource.Name, resource.Scalar, resource.Reservation.Labels.Labels[0].Value)
		}
	}

	return nil
}

type AgentState struct {
	AgentReservedResourcesFull ReservedResourcesFull `json:"reserved_resources_full,omitempty"`
}

type ReservedResourcesFull map[string]ResourceRole

type ResourceRole []Resource

type Resource struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Scalar struct {
		Value float64 `json:"value"`
	} `json:"scalar,omitempty"`
	Role        string `json:"role,omitempty"`
	Reservation struct {
		Principal string `json:"principal"`
		Labels    struct {
			Labels []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"labels"`
		} `json:"labels"`
	} `json:"reservation,omitempty"`
	Reservations []struct {
		Type      string `json:"type"`
		Role      string `json:"role"`
		Principal string `json:"principal"`
		Labels    struct {
			Labels []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"labels"`
		} `json:"labels"`
	} `json:"reservations,omitempty"`
	Disk struct {
		Persistence struct {
			ID        string `json:"id"`
			Principal string `json:"principal"`
		} `json:"persistence"`
		Volume struct {
			Mode          string `json:"mode"`
			ContainerPath string `json:"container_path"`
		} `json:"volume"`
	} `json:"disk,omitempty"`
}
