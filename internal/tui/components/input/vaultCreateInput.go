	"charm.land/huh/v2"
type VaultCreateForm struct {
	form *huh.Form
	open bool

	VaultName       string
	VaultPath       string
	Password        string
	confirmPassword string
}
