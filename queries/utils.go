package queries

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
