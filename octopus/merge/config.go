package merge

import (
	"github.com/spf13/viper"
)

type Config struct {
	DoCommit         bool
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
		//TODO reverse DoCommit to NoCommit to simplify logic
		DoCommit:         !check,
		ChunkSize:        chunk,
		ExcludedPatterns: excludedPatterns,
		Patterns:         patterns,
	}, nil
}
