package server

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/adedayo/softaudit/pkg/hash"
	"github.com/adedayo/softaudit/pkg/model"
	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v3"
)

func Download(config model.DownloadConfig) {
	shaSum, err := readLatestRDS_SHA(config)
	if err != nil {
		log.Println(err.Error())
		return
	}

	config.ExpectedSHA1 = shaSum
	rdsLocation := path.Join(serverData, shaSum)

	if _, err := os.Stat(rdsLocation); os.IsNotExist(err) {
		//decoded RDS files does not exist
		downloadISO(config)

	}
	LoadData()
}

func downloadISO(config model.DownloadConfig) {

	rds := path.Join(serverData, rdsFile)
	if _, err := os.Stat(rds); os.IsNotExist(err) {
		//download RDS file if it is not there
		log.Printf("Downloading %s ... \n", rds)
		if err = downloadFile(rds, fmt.Sprintf("%s%s", rdsURL, rdsFile)); err == nil {
			log.Println("\nDownloading Finished")
			if config.DeleteISO {
				defer func() {
					os.Remove(rds)
				}()
			}
			if sha, err := hash.FileSha(rds); err == nil {
				log.Printf("SHA1 (%s): %s \n", rds, sha)
				equal := "Equal"
				if sha != config.ExpectedSHA1 {
					equal = "NOT Equal!!!"
				}
				log.Printf("Expected SHA1 = %s, Downloaded File SHA1 = %s. %sEqual", config.ExpectedSHA1, sha, equal)
			} else {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}

	if disk, err := diskfs.Open(rds); err == nil {
		if fs, err := disk.GetFilesystem(0); err == nil {
			fInfo, err := fs.ReadDir("/")
			if err == nil {
				for _, info := range fInfo {
					if strings.Contains(strings.ToLower(info.Name()), "nsrlfile.zip") {
						extractNSRL(fs, info, config)
						continue
					}
					// if strings.Contains(strings.ToLower(info.Name()), "nsrlos") {
					// 	extractNSRLOS(fs, info, config)
					// 	continue
					// }

					if strings.Contains(strings.ToLower(info.Name()), "nsrlprod") {
						extractProd(fs, info, config)
						continue
					}

					if strings.Contains(strings.ToLower(info.Name()), "nsrlmfg") {
						extractManuFacturer(fs, info, config)
						continue
					}
				}
			}
		} else {
			log.Println(err.Error())
		}

	}
}

// func extractNSRLOS(fs filesystem.FileSystem, info fs.FileInfo, config model.DownloadConfig) {
// 	ff, err := fs.OpenFile("/"+info.Name(), os.O_RDONLY)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	defer ff.Close()
// 	scanner := bufio.NewScanner(ff)
// 	scanner.Split(bufio.ScanLines)

// 	scanner.Scan()
// 	scanner.Text() //skip first line
// 	for scanner.Scan() {

// 		d := strings.Split(scanner.Text(), ",")
// 		index := trim(&d[0])

// 		if _, exists := rdsOS[index]; !exists {

// 			data := model.NSRLOS{
// 				Name:         trim(&d[1]),
// 				Version:      trim(&d[2]),
// 				Manufacturer: trim(&d[3]),
// 			}
// 			rdsOS[index] = data
// 		}

// 	}
// 	//create directory if it doesn't exist
// 	resultPath := path.Join(serverRDS, config.ExpectedSHA1)
// 	os.MkdirAll(resultPath, 0755)

// 	out, err := os.Create(path.Join(resultPath, rdsOSFile))
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	encoder := yaml.NewEncoder(out)
// 	encoder.Encode(rdsOS)
// 	encoder.Close()
// }

func extractManuFacturer(fs filesystem.FileSystem, info fs.FileInfo, config model.DownloadConfig) {

	ff, err := fs.OpenFile("/"+info.Name(), os.O_RDONLY)
	if err != nil {
		log.Println(err.Error())
	}
	defer ff.Close()

	//create directory if it doesn't exist
	resultPath := path.Join(serverRDS, config.ExpectedSHA1)
	os.MkdirAll(resultPath, 0755)

	outTxt, err := os.Create(path.Join(resultPath, rdsManuFile))
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer outTxt.Close()

	scanner := bufio.NewScanner(ff)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	scanner.Text() //skip first line
	for scanner.Scan() {

		d := strings.Split(scanner.Text(), ",")
		index := trim(&d[0])
		manuMutex.Lock()
		if _, exists := rdsManufacturer[index]; !exists {
			name := trim(&d[1])
			data := model.NSRLManufacturer{
				Name: name,
			}
			rdsManufacturer[index] = data
			outTxt.WriteString(fmt.Sprintf("%s,%s\n", index, name))
		}
		manuMutex.Unlock()

	}

	// out, err := os.Create(path.Join(resultPath, rdsManuFile))
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }
	// encoder := yaml.NewEncoder(out)
	// manuMutex.RLock()
	// encoder.Encode(rdsManufacturer)
	// manuMutex.RUnlock()
	// encoder.Close()

}

func trim(data *string) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(*data), `"`), `"`)
}

func extractProd(fs filesystem.FileSystem, info fs.FileInfo, config model.DownloadConfig) {
	ff, err := fs.OpenFile("/"+info.Name(), os.O_RDONLY)
	if err != nil {
		log.Println(err.Error())
	}
	defer ff.Close()
	scanner := bufio.NewScanner(ff)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	scanner.Text() //skip first line
	for scanner.Scan() {

		d := strings.Split(scanner.Text(), ",")
		index := trim(&d[0])
		prodMutex.Lock()
		if _, exists := rdsProduct[index]; !exists {

			data := model.NSRLProduct{
				Name:         trim(&d[1]),
				Version:      trim(&d[2]),
				Manufacturer: trim(&d[3]),
				AppType:      trim(&d[6]),
			}
			rdsProduct[index] = data
		}
		prodMutex.Unlock()
	}
	//create directory if it doesn't exist
	resultPath := path.Join(serverRDS, config.ExpectedSHA1)
	os.MkdirAll(resultPath, 0755)

	out, err := os.Create(path.Join(resultPath, rdsProductFile))
	if err != nil {
		log.Println(err.Error())
		return
	}
	encoder := yaml.NewEncoder(out)
	prodMutex.RLock()
	encoder.Encode(rdsProduct)
	prodMutex.RUnlock()
	encoder.Close()
}

func extractNSRL(fs filesystem.FileSystem, info fs.FileInfo, config model.DownloadConfig) {

	ff, err := fs.OpenFile("/"+info.Name(), os.O_RDONLY)
	if err == nil {
		zf, err := zip.NewReader(&simpleReaderAt{r: ff}, info.Size())
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, file := range zf.File {
			if strings.Contains(strings.ToLower(file.Name), "nsrlfile.txt") {
				nsrl, err := file.Open()
				if err == nil {
					readNSRL(nsrl, config)
				} else {
					log.Println(err.Error())
				}
			}
		}
	} else {
		log.Println(err.Error())
	}

}

func readNSRL(file io.ReadCloser, config model.DownloadConfig) {
	defer file.Close()

	//create directory if it doesn't exist
	resultPath := path.Join(serverRDS, config.ExpectedSHA1)
	os.MkdirAll(resultPath, 0755)

	outTxt, err := os.Create(path.Join(resultPath, rdsNSRLFile))
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer outTxt.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	log.Println("Saving RDS File. Will take a while ... ")
	scanner.Scan()
	scanner.Text() //skip first line
	for scanner.Scan() {

		d := strings.Split(scanner.Text(), ",")
		index := strings.ToLower(trim(&d[0]))
		rdsMutex.Lock()
		if _, exists := rds[index]; !exists {
			code := trim(&d[5])
			rds[index] = code
			// log.Printf("%s, %#v\n", index, code)
			outTxt.WriteString(fmt.Sprintf("%s,%s\n", index, code))
		}
		rdsMutex.Unlock()
	}

	log.Println("Done Saving RDS File!")

}

func downloadFile(path string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, io.TeeReader(resp.Body, &downladedBytesCounter{}))
	return err
}

func readLatestRDS_SHA(config model.DownloadConfig) (string, error) {

	configFile := path.Join(serverHome, configFileName)
	if _, err := os.Stat(configFile); config.CheckAndDownloadLatest || os.IsNotExist(err) {
		//config does not exist or we are forced to check and download data
		var cFile *os.File

		if os.IsNotExist(err) {
			cFile, err = os.Create(configFile)
			if err != nil {
				return "", err
			}
		} else {
			cFile, err = os.Open(configFile)
			if err != nil {
				return "", err
			}
		}
		defer cFile.Close()

		decoder := yaml.NewDecoder(cFile)
		var serviceConfig model.ServiceConfig
		decoder.Decode(&serviceConfig)
		if err != nil {
			return "", err
		}

		resp, err := http.Get(fmt.Sprintf("%s%s.sha", rdsURL, rdsFile))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		buff, err := io.ReadAll(resp.Body)
		if err == nil {
			ss := strings.Split(string(buff), "=")
			shaSum := ""
			if len(ss) > 0 {
				shaSum = strings.TrimSpace(ss[1])
			}
			if serviceConfig.LatestISO_SHA1 != shaSum {
				serviceConfig.LatestISO_SHA1 = shaSum
				encoder := yaml.NewEncoder(cFile)
				if err := encoder.Encode(serviceConfig); err != nil {
					return "", err
				}
				encoder.Close()
				return shaSum, nil
			}
		} else {
			log.Println(err.Error())
		}
	}

	return "", nil
}

type simpleReaderAt struct {
	r io.ReadWriteSeeker
}

func (sra *simpleReaderAt) ReadAt(b []byte, offset int64) (n int, err error) {
	_, err = sra.r.Seek(offset, io.SeekStart)
	if err != nil {
		return n, err
	}
	n, err = sra.r.Read(b)
	return
}

type downladedBytesCounter struct {
	total uint64
}

func (dbc downladedBytesCounter) showProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 25))
	fmt.Printf("\r%s Downloaded ...", humanize.Bytes(dbc.total))
}

func (dbc *downladedBytesCounter) Write(data []byte) (int, error) {
	count := len(data)
	dbc.total += uint64(count)
	dbc.showProgress()
	return count, nil
}
