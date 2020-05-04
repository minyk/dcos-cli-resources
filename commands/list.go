package commands

import (
	"github.com/minyk/dcos-resources/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type resourceListHandler struct {
	q       *queries.ResourceList
	agentID string
	role    string
}

// HandleScheduleSection
func HandleListResourcesSection(app *kingpin.Application, q *queries.ResourceList) {
	HandleListResourcesCommands(app.Command("list", "List reserved resources from agent").Alias("listresources"), q)
}

func (cmd *resourceListHandler) handleListResources(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.ListResourcesFromNode(cmd.agentID, cmd.role)
}

func HandleListResourcesCommands(resources *kingpin.CmdClause, q *queries.ResourceList) {
	cmd := &resourceListHandler{q: q}
	listResources := resources.Action(cmd.handleListResources)
	listResources.Flag("agent-id", "Agent ID to list").Required().StringVar(&cmd.agentID)
	listResources.Flag("role", "Role for list").Required().StringVar(&cmd.role)
}
