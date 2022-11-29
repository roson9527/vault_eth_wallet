package factory

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
)

type IPath interface {
	Path() []*framework.Path
}

var paths = make(map[string]IPath)

func Register(name string, iPath IPath) {
	if _, ok := paths[name]; ok {
		panic("path already registered")
	}
	paths[name] = iPath
}

func Do() []*framework.Path {
	var out []*framework.Path
	for name, iPath := range paths {
		ps := iPath.Path()
		hclog.Default().Info("path:factory:do", "package", name, "paths.length", len(ps))
		for _, p := range ps {
			ops := make([]string, 0)
			if p.Operations != nil {
				for k, _ := range p.Operations {
					ops = append(ops, string(k))
				}
			}
			// 输出日志
			hclog.Default().Info("path:factory:do", "package", name, "path", p.Pattern, "operations", ops)
		}

		out = append(out, iPath.Path()...)
	}
	return out
}
