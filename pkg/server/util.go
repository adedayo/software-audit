package server

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"

	"github.com/adedayo/softaudit/pkg/model"
	"gopkg.in/yaml.v3"
)

func LoadData() {

	cFile, err := os.Open(path.Join(serverHome, configFileName))
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	decoder := yaml.NewDecoder(cFile)
	var config model.ServiceConfig
	decoder.Decode(&config)
	loadManufacturer(config)
	loadProduct(config)
	loadNSRL(config)
}

// func loadOS(config model.ServiceConfig) {
// 	file, err := os.Open(path.Join(serverRDS, config.LatestISO_SHA1, rdsOSFile))
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 		return
// 	}
// 	defer file.Close()
// 	decoder := yaml.NewDecoder(file)
// 	err = decoder.Decode(rdsOS)
// 	if err != nil {
// 		log.Printf("%v\n", err)
// 	}
// }

func loadManufacturer(config model.ServiceConfig) {
	file, err := os.Open(path.Join(serverRDS, config.LatestISO_SHA1, rdsManuFile))
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	manuMutex.Lock()
	for scanner.Scan() {
		d := strings.Split(scanner.Text(), ",")
		rdsManufacturer[d[0]] = model.NSRLManufacturer{Name: d[1]}
	}
	manuMutex.Unlock()

}

func loadProduct(config model.ServiceConfig) {
	file, err := os.Open(path.Join(serverRDS, config.LatestISO_SHA1, rdsProductFile))
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	prodMutex.RLock()
	err = decoder.Decode(rdsProduct)
	prodMutex.RUnlock()
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func loadNSRL(config model.ServiceConfig) {
	file, err := os.Open(path.Join(serverRDS, config.LatestISO_SHA1, rdsNSRLFile))
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	rdsMutex.Lock()

	for scanner.Scan() {
		d := strings.Split(scanner.Text(), ",")
		rds[d[0]] = d[1]
	}

	rdsMutex.Unlock()

}
