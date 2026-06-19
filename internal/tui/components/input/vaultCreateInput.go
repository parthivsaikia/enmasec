package input

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/tui/messages"
	"github.com/parthivsaikia/enmasec/internal/validation"
)

type VaultCreateForm struct {
	form *huh.Form
	open bool

	VaultName       string
	VaultPath       string
	Password        string
	confirmPassword string
}

func NewVaultCreateForm() *VaultCreateForm {
	vcf := &VaultCreateForm{
		open:      false,
		VaultPath: store.GetEnmasecDirLocation(),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("vaultName").
				Title("Vault Name").
				Validate(validation.ValidateVaultName),
			huh.NewInput().
				Key("vaultPath").
				Title("Vault Path").
				Value(&vcf.VaultPath),
			huh.NewInput().
				Key("password").
				Title("Password").
				Validate(validation.CheckPasswordValid),
			huh.NewInput().
				Title("Confirm Password").
				Validate(func(string) error {
					if vcf.confirmPassword != vcf.Password {
						return fmt.Errorf("passwords don't match")
					}
					return nil
				}),
		),
	)
	form.WithShowHelp(false)
	// form.WithShowErrors(false)

	vcf.form = form
	return vcf
}

func (v *VaultCreateForm) Init() tea.Cmd {
	return v.form.Init()
}

func (v *VaultCreateForm) Update(msg tea.Msg) tea.Cmd {
	form, cmd := v.form.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		v.form = f
	}

	if v.form.State == huh.StateCompleted {
		v.Close()
		return messages.VaultCreatedMsg(v.VaultPath, v.VaultName, v.Password)
	}

	return cmd
}

func (v *VaultCreateForm) View() string {
	return v.form.View()
}

func (v *VaultCreateForm) IsOpen() bool {
	return v.open
}

func (v *VaultCreateForm) Open() {
	v.form.State = huh.StateNormal
	v.open = true
}

func (v *VaultCreateForm) Close() {
	v.open = false
}
