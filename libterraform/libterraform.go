package libterraform

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform/config/module"
	tf "github.com/hashicorp/terraform/terraform"
)

// Terraform is the interface that does terraform operations.
type Terraform interface {
	// Apply applies the blueprint to a deploy specified in state
	Apply(path string, state *tf.State, stateCh chan<- tf.State) Task

	// Destroy destroyes the resources specified in state.
	Destroy(path string, state *tf.State, stateCh chan<- tf.State) Task
}

// Context represent all the context needed to do a terraform plan and apply.
type Context interface {
	// Apply does a apply or a destroy
	Apply() (*tf.State, error)

	// Plan does a plan before doing an apply
	Plan() (*tf.Plan, error)

	// Stop stops an execution
	Stop()
}

type libterraform struct {
	createContext func(path string, destroy bool, state *tf.State, stateCh chan<- tf.State) (Context, error)
}

func NewTerraform() *libterraform {
	return &libterraform{
		createContext: createContext,
	}
}

func createContext(path string, destroy bool, state *tf.State, stateCh chan<- tf.State) (Context, error) {
	// Load the root module
	mod, err := module.NewTreeModule("", path)
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %s", err)
	}

	storage := &getter.FolderStorage{
		StorageDir: filepath.Join(path, "modules"),
	}
	err = mod.Load(storage, module.GetModeGet) // Hard-code to GetModeGet
	if err != nil {
		return nil, fmt.Errorf("Error downloading modules: %s", err)
	}

	// Get providers
	BuiltinConfig.Discover()
	fmt.Println(BuiltinConfig.Providers)
	providerFactories := BuiltinConfig.ProviderFactories()

	// New hooks
	stateHook := new(StateChannelHook)
	stateHook.stateCh = stateCh

	opts := &tf.ContextOpts{
		Destroy:      destroy,
		Module:       mod,
		Parallelism:  10,
		Targets:      nil,
		Variables:    nil,
		State:        state,
		Providers:    providerFactories,
		Hooks:        []tf.Hook{stateHook},
		Diff:         nil,
		Provisioners: nil,
		UIInput:      nil,
	}

	return tf.NewContext(opts), nil
}
