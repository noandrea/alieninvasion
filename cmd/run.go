/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/noandrea/alieninvasion/aliens"
	"github.com/noandrea/alieninvasion/land"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliensN, maxIterations int
var landMapFile, landOutFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVarP(&aliensN, "aliens", "n", 10, "Number of invading aliens (default: 10)")
	runCmd.Flags().IntVar(&maxIterations, "iterations", 100000, "Max number of iterations before stopping the exectution")
	runCmd.Flags().StringVarP(&landMapFile, "map", "i", "map.txt", "The file containing the land map of the world")
	runCmd.Flags().StringVarP(&landOutFile, "aftermath", "o", "map.post.txt", "The file containing the land map of the world after the invasion")
}

func run(cmd *cobra.Command, args []string) {
	fmt.Printf(`
╔═╗┬  ┬┌─┐┌┐┌╦┌┐┌┬  ┬┌─┐┌─┐┬┌─┐┌┐┌
╠═╣│  │├┤ │││║│││└┐┌┘├─┤└─┐││ ││││
╩ ╩┴─┘┴└─┘┘└┘╩┘└┘ └┘ ┴ ┴└─┘┴└─┘┘└┘ v%s`, rootCmd.Version)
	fmt.Println()
	fmt.Println("Welcome to Alien Invasion")
	fmt.Println("This invasion will see", aliensN, "aliens invading the planet")
	// initialize a new land
	world := land.NewLand()
	// load the map
	log.Debug("Input map file is ", landMapFile)
	err := land.LoadFromFile(world, landMapFile)
	if err != nil {
		fmt.Println("Something is not right, I cannot read the world map at", landMapFile)
		log.Debug(err)
		return
	}
	n, e := land.Size(world)
	fmt.Println("This planet has", n, "cities interconnected by", e, "routes")
	// initialize the invasion
	invasion := aliens.NewInvasion(world, aliensN, maxIterations)
	fmt.Println("The invasion has begun!")
	for {
		round, err := aliens.Run(world, invasion)
		fmt.Println("📅 Day", round, "has ended")
		if err != nil {
			fmt.Println("The invasion is over!")
			n, e = land.Size(world)
			fmt.Println("This planet is left with ", n, "cities interconnected by", e, "routes")
			log.Debug(err)
			land.SaveToFile(world, landOutFile)
			fmt.Println("The map was saved in ", landOutFile)
			return
		}
	}

}
