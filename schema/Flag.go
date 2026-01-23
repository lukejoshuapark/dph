package schema

type Flag struct {
	Key              string
	Description      string
	DocumentationURL string
	Exclude          []string
}

func FlagFromDefinition(key string, def *FlagDefinition) *Flag {
	return &Flag{
		Key:              key,
		Description:      def.Description,
		DocumentationURL: def.DocumentationURL,
		Exclude:          def.Exclude,
	}
}
