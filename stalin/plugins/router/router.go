package plugin_router

import (
	"stalin/message"
	. "stalin/plugins"
)

func init() {
	RegisterPlugin("Router", new(Router))
}

type Router struct {
	gConfig *GlobalConfig
}

func (r *Router) Init(g *GlobalConfig, name string) (Plugin, error) {
	plugin := &Router{gConfig: g}
	return plugin, nil
}

func (r *Router) Inject(msg *message.Message) error {
	for _, plugin := range r.gConfig.Outputs {
		plugin.Inject(msg)
	}
	return nil
}

func (r *Router) Run() error {
	return nil
}
