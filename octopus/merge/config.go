package merge

import (
	"github.com/spf13/viper"
)

type Config struct {
	NoCommit         bool
	ChunkSize        int
	ExcludedPatterns []string
	Patterns         []string
}

func GetConfig(args []string) (*Config, error) {
	var check = viper.GetBool("check")
	var chunk = viper.GetInt("chunk")
	var excludedPatterns = viper.GetStringSlice("exclude")
	var patterns = args

	return &Config{
		NoCommit:         check,
		ChunkSize:        chunk,
		ExcludedPatterns: excludedPatterns,
		Patterns:         patterns,
	}, nil
}
