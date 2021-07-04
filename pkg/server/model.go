package server

import (
	"os"
	"path"
	"sync"

	"github.com/adedayo/softaudit/pkg/model"
	"github.com/mitchellh/go-homedir"
)

//see https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl
var (
	rdsURL         = "https://s3.amazonaws.com/rds.nsrl.nist.gov/RDS/current/"
	serverData     = ""
	serverRDS      = ""
	serverHome     = ""
	rdsFile        = "RDS_modern.iso"
	configFileName = "audit_service.yaml"
	rdsOSFile      = "rdsOS.yaml"
	rdsManuFile    = "rdsManufacturer.txt"
	rdsProductFile = "rdsProduct.yaml"
	rdsNSRLFile    = "rdsNSRL.txt"
	nothing        = struct{}{}
	rds            = make(map[string]string, 0)
	rdsMutex       = sync.RWMutex{}
	rdsProduct     = make(map[string]model.NSRLProduct)
	prodMutex      = sync.RWMutex{}
	// rdsOS           = make(map[string]model.NSRLOS)
	rdsManufacturer = make(map[string]model.NSRLManufacturer)
	manuMutex       = sync.RWMutex{}
)

func init() {

	if loc, err := homedir.Expand("~/.audit-server"); err == nil {
		serverHome = loc
		serverData = path.Join(serverHome, "data")
		serverRDS = path.Join(serverData, "rds")
	}

	//create server data/rds directory if it doesn't exist
	os.MkdirAll(serverRDS, 0755)

}
