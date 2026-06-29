package input

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/parthivsaikia/enmasec/internal/core"
	"github.com/parthivsaikia/enmasec/internal/store"
	"github.com/parthivsaikia/enmasec/internal/tui/messages"
	"github.com/parthivsaikia/enmasec/internal/validation"
)

type VaultCreateForm struct {
	form      *huh.Form
	open      bool
	VaultName string
	VaultPath string
	Password  string
}

// cmd for creating vault
func vaultCreateCmd(v *VaultCreateForm) tea.Cmd {
	return func() tea.Msg {
		v.VaultName = v.form.GetString("vaultName")
		v.VaultPath = v.form.GetString("vaultPath")
		v.Password = v.form.GetString("password")
		vault, err := core.CreateVault(v.VaultPath, v.VaultName, v.Password)
		if err != nil {
			log.Print(err)
			return messages.ErrMsg(err)
		}
		log.Print("vault created successfully")
		return messages.VaultCreateMsg(vault)
	}
}

func NewVaultCreateForm() *VaultCreateForm {
	vcf := &VaultCreateForm{
		open:      false,
		VaultPath: store.GetEnmasecDirLocation(),
	}

	Form := huh.NewForm(
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
				EchoMode(huh.EchoModePassword).
				Validate(validation.CheckPasswordValid),
			huh.NewInput().
				Key("confirmPassword").
				Title("Confirm Password").
				EchoMode(huh.EchoModePassword).
				Validate(func(confirmPassword string) error {
					password := vcf.form.GetString("password")
					if password != confirmPassword {
						log.Print("password: ", password)
						log.Print("confirmPassword: ", confirmPassword)
						return fmt.Errorf("passwords don't match")
					}
					return nil
				}),
		),
	)
	Form.WithShowHelp(false)
	// Form.WithShowErrors(false)

	vcf.form = Form
	return vcf
}

func (v *VaultCreateForm) Init() tea.Cmd {
	return v.form.Init()
}

func (v *VaultCreateForm) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	Form, cmd := v.form.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if f, ok := Form.(*huh.Form); ok {
		v.form = f
	}

	if v.form.State == huh.StateCompleted {
		cmd = vaultCreateCmd(v)
		log.Print("cmd: ", cmd)
		cmds = append(cmds, cmd)
		v.Close()
	}

	return tea.Batch(cmds...)
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
