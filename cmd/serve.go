/*
Copyright © 2021 Adedayo Adetoye (aka Dayo)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"github.com/adedayo/softaudit/pkg/model"
	"github.com/adedayo/softaudit/pkg/server"
	"github.com/spf13/cobra"
)

var (
	local, checkDownload, deleteDownload bool
	port                                 int
	defaultServerPort                    = 7454
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run software audit API service",
	Long:  `Run software audit API service`,
	Run:   runService,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().BoolVar(&local, "local", false, "Bind service to localhost")
	serveCmd.Flags().IntVarP(&port, "port", "p", defaultServerPort, "Service port")
	serveCmd.Flags().BoolVar(&checkDownload, "check-download", true, "Check NSRL for new versions of signatures and download as needed")
	serveCmd.Flags().BoolVar(&deleteDownload, "delete-download", false, "Delete downloaded signature archive file after extracting needed information")

}

func runService(cmd *cobra.Command, args []string) {
	go server.Download(model.DownloadConfig{
		CheckAndDownloadLatest: checkDownload,
		DeleteISO:              deleteDownload,
	})
	server.ServeAPI(model.Config{
		ApiPort: port,
		Local:   local,
	})
}
