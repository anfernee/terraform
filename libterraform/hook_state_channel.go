package libterraform

import (
	"sync"

	"github.com/hashicorp/terraform/terraform"
)

// StateChannelHook is a hook that continuously send updated state through a channel
type StateChannelHook struct {
	terraform.NilHook
	sync.Mutex

	stateCh chan<- terraform.State // write channel to write state back
}

func (h *StateChannelHook) PostStateUpdate(
	s *terraform.State) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if h.stateCh != nil {
		h.stateCh <- *s
	}

	// Continue forth
	return terraform.HookActionContinue, nil
}
