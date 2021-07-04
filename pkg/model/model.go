package model

type ExecPaths struct {
	PathRoots []string
}

type Config struct {
	ApiPort int
	Local   bool //if set, to bind the api to localhost:port; or simply :port (web service) instead
}

type SHAQueries struct {
	Queries []Query
}

type Query struct {
	SHA1 string
	Path string
}

type Response struct {
	SHA1, Path string
	Product    NSRLProduct
	Found      bool
}

type ServiceConfig struct {
	LatestISO_SHA1 string
}

type DownloadConfig struct {
	CheckAndDownloadLatest bool
	DeleteISO              bool
	ExpectedSHA1           string
}

type NSRLOS struct {
	Name, Version, Manufacturer string
}

type NSRLManufacturer struct {
	Name string
}

type NSRLProduct struct {
	//Code string //indexed by code in memory
	Name, Version, Manufacturer, AppType string
}
