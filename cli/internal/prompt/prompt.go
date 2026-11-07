package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/notfoundy/dotfiles/cli/internal/onepass"
)

func SelectRoles(available, preselected []string) ([]string, error) {
	checked := make(map[string]bool, len(preselected))
	for _, r := range preselected {
		checked[r] = true
	}

	options := make([]huh.Option[string], 0, len(available))
	for _, role := range available {
		options = append(options, huh.NewOption(role, role).Selected(checked[role]))
	}

	selected := append([]string(nil), preselected...)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Roles").
				Description("space to toggle, enter to run").
				Options(options...).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return nil, err
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("no role selected")
	}
	return selected, nil
}

func SelectAccount(accounts []onepass.Account) (onepass.Account, error) {
	options := make([]huh.Option[string], 0, len(accounts))
	byRef := make(map[string]onepass.Account, len(accounts))
	for _, a := range accounts {
		options = append(options, huh.NewOption(a.Label(), a.Ref()))
		byRef[a.Ref()] = a
	}

	var ref string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("1Password account").
				Options(options...).
				Value(&ref),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return onepass.Account{}, err
	}
	return byRef[ref], nil
}

func AskDomain(defaultDomain string) (string, error) {
	domain := defaultDomain

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("1Password sign-in address").
				Value(&domain),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return "", err
	}
	if domain == "" {
		return defaultDomain, nil
	}
	return domain, nil
}
