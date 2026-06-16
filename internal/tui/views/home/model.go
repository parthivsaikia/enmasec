type currentPane int

const (
	vaultPane = iota
	servicePane
	accountPane
)

type model struct {
	// data
	runtimeIndex *models.RuntimeIndex
	// panes
	vaultModel vault.Model
	// cursor
	pane currentPane
}
