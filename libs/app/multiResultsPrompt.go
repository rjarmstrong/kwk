package app

import (
	"github.com/kwk-links/kwk-cli/libs/openers"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
)

func NewMultiResultPrompt(o openers.IOpen, i gui.IInteraction) *MultiResultPrompt {
	return &MultiResultPrompt{Openers:o, Interaction:i}
}

type MultiResultPrompt struct {
	Openers openers.IOpen
	Interaction gui.IInteraction
}

func (m *MultiResultPrompt) CheckAndPrompt(fullKey string, list *api.AliasList, args []string){
	if list.Total == 1 {
		m.Openers.Open(&list.Items[0], args[1:])
	} else if list.Total > 1 {
		m.Interaction.Respond("chooseBetweenKeys", list.Items)
	} else {
		m.Interaction.Respond("notfound", fullKey)
	}
}
