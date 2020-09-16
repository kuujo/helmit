package {{ .Client.Package.Name }}

import (
    "github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type {{ .Client.Types.Interface }} interface {
    {{ .Resource.Names.Plural }}() {{ .Reader.Types.Interface }}
}

func New{{ .Client.Types.Interface }}(resources resource.Client, filter resource.Filter) {{ .Client.Types.Interface }} {
	return &{{ .Client.Types.Struct }}{
		Client: resources,
		filter: filter,
	}
}

type {{ .Client.Types.Struct }} struct {
	resource.Client
	filter resource.Filter
}

func (c *{{ .Client.Types.Struct }}) {{ .Resource.Names.Plural }}() {{ .Reader.Types.Interface }} {
    return New{{ .Reader.Types.Interface }}(c.Client, c.filter)
}
