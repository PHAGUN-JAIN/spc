package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dvdmuckle/spc/cmd/helper"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var seekCmd = &cobra.Command{
	Use:   "seek",
	Args:  cobra.ExactArgs(1),
	Short: "Seek to a specific position in the currently playing song from Spotify",
	Long: `Seek to a specific position in the currently playing song from Spotify. This command requires
	exactly one argument, a number between 0 and the length of the currently playing song in seconds to seek to.`,
	Run: func(cmd *cobra.Command, args []string) {
		helper.SetClient(&conf)
		seconds, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Passed value for seconds must be an integer.")
			os.Exit(1)
		}

		err = conf.Client.Seek(seconds * 1000)
		if err != nil {
			glog.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(seekCmd)
}
