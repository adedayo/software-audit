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
	"fmt"
	"log"

	find "github.com/adedayo/softaudit/pkg/find"
	hash "github.com/adedayo/softaudit/pkg/hash"
	"github.com/gorilla/websocket"

	"github.com/adedayo/softaudit/pkg/model"
	"github.com/spf13/cobra"
)

var (
	defaultServerURL = fmt.Sprintf("ws://localhost:%d", defaultServerPort)
	serverURL        string
)

// appsCmd represents the software command
var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Audit applications.",
	Long:  `Audit applications on your file system.`,
	Run:   apps,
}

func init() {
	rootCmd.AddCommand(appsCmd)

	appsCmd.Flags().StringVarP(&serverURL, "server", "s", defaultServerURL, "Websocket URL of the resolving server")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func apps(cmd *cobra.Command, args []string) {

	ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/api/check", serverURL), nil)
	if err != nil {
		log.Printf("%v\nEnsure the Server is running. Terminating ... ", err)
		return
	}
	defer ws.Close()
	go resultHandler(ws)

	for exe := range find.Executables(model.ExecPaths{
		PathRoots: args,
	}) {
		if hashSum, err := hash.FileSha(exe); err == nil {
			// println(exe, hashSum)
			ws.WriteJSON(model.Query{
				Path: exe,
				SHA1: hashSum,
			})
		}
	}
}

func resultHandler(ws *websocket.Conn) {
	for {
		var resp model.Response
		err := ws.ReadJSON(&resp)
		if err != nil {
			log.Printf("%v\n", err)
			break
		}
		if resp.Found {
			fmt.Printf("%#v\n", resp)
		}
	}
}
