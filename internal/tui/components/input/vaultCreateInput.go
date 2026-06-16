package input

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
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
				Title("Vault Name").
				Validate(validation.ValidateVaultName).
				Value(&vcf.VaultName),
			huh.NewInput().
				Title("Vault Path").
				Value(&vcf.VaultPath),
			huh.NewInput().
				Title("Password").
				Validate(validation.CheckPasswordValid).
				Value(&vcf.Password),
			huh.NewInput().
				Title("Confirm Password").
				Value(&vcf.confirmPassword).
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
