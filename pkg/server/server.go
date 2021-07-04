package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/adedayo/softaudit/pkg/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	routes         = mux.NewRouter()
	apiVersion     = "0.0.0"
	allowedOrigins = []string{}
	corsOptions    = []handlers.CORSOption{
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "Accept-Language", "Origin"}),
		handlers.AllowCredentials(),
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			for _, origin := range allowedOrigins {
				if origin == r.Host {
					return true
				}
			}
			return strings.Split(r.Host, ":")[0] == "localhost" //allow localhost independent of port
		},
	}
)

func init() {
	addRoutes()
}

func addRoutes() {
	routes.HandleFunc("/api/check", check).Methods("GET")
	routes.HandleFunc("/admin/getLatest", check).Methods("GET")

}

func check(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading websocket connection %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer ws.Close()

	for {
		var query model.Query
		err = ws.ReadJSON(&query)
		if err != nil {
			log.Printf("Error %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		}
		ws.WriteJSON(checkSignature(&query))
	}

}

func checkSignature(query *model.Query) *model.Response {

	resp := model.Response{
		SHA1: query.SHA1,
		Path: query.Path,
	}

	rdsMutex.RLock()
	if pCode, exists := rds[strings.ToLower(query.SHA1)]; exists {
		resp.Found = exists
		prodMutex.RLock()
		if product, exists := rdsProduct[pCode]; exists {
			resp.Product = product
			manuMutex.RLock()
			if manf, exists := rdsManufacturer[product.Manufacturer]; exists {
				resp.Product.Manufacturer = manf.Name
			}
			manuMutex.RUnlock()
		}
		prodMutex.RUnlock()
	}
	rdsMutex.RUnlock()
	return &resp
}

func ServeAPI(config model.Config) {
	hostPort := "localhost:%d"
	if config.Local {
		//localhost electron app
		corsOptions = append(corsOptions, handlers.AllowedOrigins(allowedOrigins))
	} else {
		hostPort = ":%d"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(hostPort, config.ApiPort), handlers.CORS(corsOptions...)(routes)))
}
