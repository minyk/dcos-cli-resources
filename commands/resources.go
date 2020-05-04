package commands

import (
	"github.com/minyk/dcos-resources/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type reserveResourcesHandler struct {
	q             *queries.Resources
	agentID       string
	role          string
	principal     string
	frameworkID   string
	cpus          float64
	cpuLabel      string
	mem           float64
	memLabel      string
	disk          float64
	diskLabel     string
	persistid     string
	containerpath string
	hostpath      string
}

func (cmd *reserveResourcesHandler) handleReserve(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.ReserveResource(cmd.agentID, cmd.role, cmd.principal, cmd.cpus, cmd.mem)
}

// HandleScheduleSection
func HandleResourcesSection(app *kingpin.Application, q *queries.Resources) {
	HandleReserveResourcesCommands(app.Command("reserve", "Reserve resources").Alias("reserves"), q)
}

// HandleScheduleCommand
func HandleReserveResourcesCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &reserveResourcesHandler{q: q}
	reserve := resources.Action(cmd.handleReserve)
	reserve.Flag("agent-id", "Agent ID to reserve").Required().StringVar(&cmd.agentID)
	reserve.Flag("role", "Role for reserve").Required().StringVar(&cmd.role)
	reserve.Flag("principal", "Principal for reserve").Default("my-principal").StringVar(&cmd.principal)
	reserve.Flag("cpus", "Amount of cpus to reserve").Default("0").Float64Var(&cmd.cpus)
	reserve.Flag("mem", "Amount of memory to reserve. The unit is MB.").Default("0").Float64Var(&cmd.mem)
}
