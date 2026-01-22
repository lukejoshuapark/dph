package schema

type Flag struct {
	Key         string
	Description string
	Exclude     []string
}

func FlagFromDefinition(key string, def *FlagDefinition) *Flag {
	return &Flag{
		Key:         key,
		Description: def.Description,
		Exclude:     def.Exclude,
	}
}
