package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/dvdmuckle/spc/cmd/helper"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

// saveDaylistCmd represents the save-daylist command
var saveDaylistCmd = &cobra.Command{
	Use:   "save-daylist",
	Short: "Saves the current Spotify Daylist playlist",
	Long: `Saves the current Spotify Daylist playlist.
Note this saves the Daylist at the current time of day and cannot retrieve past Daylists.`,
	Run: func(cmd *cobra.Command, args []string) {
		helper.SetClient(&conf)
		playlistName, _ := cmd.Flags().GetString("name")
		playlistDescription := "Spotify Daylist for " + getDaylistDateTime()
		if playlistName == "" {
			playlistName = "Daylist " + getDaylistDateTime()
		}
		isPublic, _ := cmd.Flags().GetBool("public")
		isCollaborative, _ := cmd.Flags().GetBool("collaborative")

		ctx := context.Background()
		currentUser, err := conf.Client.CurrentUser(ctx)
		if err != nil {
			helper.LogErrorAndExit(err)
		}
		if deduplicatePlaylist(playlistName, currentUser.User.ID) {
			fmt.Println("Daylist already saved")
			return
		}
		newPlaylist, err := conf.Client.CreatePlaylistForUser(ctx, currentUser.User.ID, playlistName, playlistDescription, isPublic, isCollaborative)
		if err != nil {
			helper.LogErrorAndExit(err)
		}
		searchResult, err := conf.Client.Search(ctx, "Daylist", spotify.SearchTypePlaylist)
		if err != nil {
			helper.LogErrorAndExit(err)
		}
		var daylistPlaylist spotify.ID
		for _, playlist := range searchResult.Playlists.Playlists {
			if playlist.Owner.ID == "spotify" {
				daylistPlaylist = playlist.ID
				break
			}
		}
		daylistTracks := func() spotify.PlaylistItemPage {
			playlistTracks, err := conf.Client.GetPlaylistItems(ctx, daylistPlaylist)
			if err != nil {
				helper.LogErrorAndExit(err)
			}
			return *playlistTracks
		}
		var daylistTrackIDs []spotify.ID
		for _, track := range daylistTracks().Items {
			daylistTrackIDs = append(daylistTrackIDs, track.Track.Track.ID)
		}
		_, err = conf.Client.AddTracksToPlaylist(ctx, newPlaylist.ID, daylistTrackIDs...)
		if err != nil {
			helper.LogErrorAndExit(err)
			return
		}
		fmt.Printf("Daylist saved as %s\n", playlistName)
	},
}

func getDaylistDateTime() string {
	// Gets the current date and time (hour) to differentiate Daylists
	now := time.Now()
	return fmt.Sprintf("%s %02d:%02d", now.Weekday(), now.Hour(), now.Minute())
}

func init() {
	rootCmd.AddCommand(saveDaylistCmd)
	saveDaylistCmd.Flags().StringP("name", "n", "", "Custom name for the saved playlist")
	saveDaylistCmd.Flags().BoolP("public", "p", false, "Whether to make the new playlist public")
	saveDaylistCmd.Flags().BoolP("collaborative", "c", false, "Set the playlist as collaborative")
}
