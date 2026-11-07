package system

import (
	"fmt"
	"strings"

	"github.com/notfoundy/dotfiles/cli/internal/shell"
	"github.com/notfoundy/dotfiles/cli/internal/ui"
)

var baseTools = []string{"ansible-galaxy", "ansible-playbook", "git"}

func CheckBaseTools(p *ui.Printer) error {
	task := p.Start("Checking base tools")

	var missing []string
	for _, tool := range baseTools {
		if !shell.Exists(tool) {
			missing = append(missing, tool)
		}
	}
	if len(missing) > 0 {
		task.Fail()
		return fmt.Errorf("%s not found — they ship with the bootc image, see image/packages/", strings.Join(missing, ", "))
	}

	task.Done()
	return nil
}
