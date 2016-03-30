package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/memberlist"
	"github.com/spf13/viper"
)

const (
	// Name of this program
	Name = "sysnet"
	// WAN network type
	WAN = "wan"
	// LAN network type
	LAN = "lan"
	// Local network type
	Local = "local"
)

func main() {
	// Set config defaults
	viper.SetDefault("host", "localhost")
	viper.SetDefault("net", WAN)
	viper.SetDefault("configType", "yaml")

	// Tell config to read from env
	viper.SetEnvPrefix(Name)
	viper.AutomaticEnv()

	// Add config file
	viper.SetConfigName("config")
	viper.SetConfigType(viper.GetString("configType"))
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/" + Name)

	// Read config file
	err := viper.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		log.Println("Config file not found")
	} else if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s\n", err))
	}

	// Configure memberlist
	var config *memberlist.Config
	switch viper.GetString("net") {
	default:
		log.Fatal(fmt.Errorf(
			"Error configuring network type %q is not a vaild network type",
			viper.GetString("net"),
		))
	case WAN:
		config = memberlist.DefaultWANConfig()
	case LAN:
		config = memberlist.DefaultLANConfig()
	case Local:
		config = memberlist.DefaultLocalConfig()
	}

	// Create it based on the config
	list, err := memberlist.Create(config)
	if err != nil {
		log.Fatal("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member
	n, err := list.Join([]string{viper.GetString("host")})
	if err != nil {
		log.Fatal("Failed to join cluster: " + err.Error())
	}
	log.Printf("Joined cluster of %d nodes\n", n)

	// Ask for members of the cluster
	for _, member := range list.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}
}
