package flags

import (
	"github.com/spf13/cobra"
	"log"
)

var cmd *cobra.Command

func Init(command *cobra.Command) {
	cmd = command
}

func GetString(name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		log.Fatalf("Failed to get value of %s flag: %s", name, err)
	}
	return value
}

func GetBool(name string) bool {
	value, err := cmd.Flags().GetBool(name)
	if err != nil {
		log.Fatalf("Failed to get value of %s flag: %s", name, err)
	}
	return value
}

func GetUint64(name string) uint64 {
	value, err := cmd.Flags().GetUint64(name)
	if err != nil {
		log.Fatalf("Failed to get value of %s flag: %s", name, err)
	}
	return value
}
