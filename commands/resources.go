package commands

import (
	"github.com/minyk/dcos-resources/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type resourcesHandler struct {
	q         *queries.Resources
	agentID   string
	role      string
	principal string
	cpus      float64
	mem       float64

}

func (cmd *resourcesHandler) handleReserve(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.ReserveResource(cmd.agentID, cmd.role, cmd.principal, cmd.cpus, cmd.mem)
}

func (cmd *resourcesHandler) handleUnreserve(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.UnreserveResource(cmd.agentID, cmd.role, cmd.principal, cmd.cpus, cmd.mem)
}

// HandleScheduleSection
func HandleResourcesSection(app *kingpin.Application, q *queries.Resources) {
	HandleReserveResourcesCommands(app.Command("reserve", "Reserve resources").Alias("reserves"), q)
	HandleUnreserveResourcesCommands(app.Command("unreserve", "Unreserve resources").Alias("unreserves"), q)
}

// HandleScheduleCommand
func HandleReserveResourcesCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	reserve := resources.Action(cmd.handleReserve)
	reserve.Flag("agent-id", "Agent ID to reserve").Required().StringVar(&cmd.agentID)
	reserve.Flag("role", "Role for reserve").Required().StringVar(&cmd.role)
	reserve.Flag("principal", "Principal for reserve").Default("my-principal").StringVar(&cmd.principal)
	reserve.Flag("cpus", "Amount of cpus to reserve").Default("0").Float64Var(&cmd.cpus)
	reserve.Flag("mem", "Amount of memory to reserve. The unit is MB.").Default("0").Float64Var(&cmd.mem)
}

// HandleScheduleCommand
func HandleUnreserveResourcesCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	unReserve := resources.Action(cmd.handleUnreserve)
	unReserve.Flag("agent-id", "Agent ID to unreserve").Required().StringVar(&cmd.agentID)
	unReserve.Flag("role", "Role for unreserve").Required().StringVar(&cmd.role)
	unReserve.Flag("principal", "Principal for unreserve.").Default("my-principal").StringVar(&cmd.principal)
	unReserve.Flag("cpus", "Amount of cpus to unreserve").Default("0").Float64Var(&cmd.cpus)
	unReserve.Flag("mem", "Amount of memory to unreserve. The unit is MB.").Default("0").Float64Var(&cmd.mem)
}
