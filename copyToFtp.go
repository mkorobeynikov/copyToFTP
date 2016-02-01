package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dutchcoders/goftp"
	"github.com/mkorobeynikov/copyToFTP/ftp"
)

type Configuration struct {
	// FTP Host "example.com:21"
	FtpHost string `json:"ftpHost"`

	// FTP User
	FtpUser string `json:"ftpUser"`

	// FTP Password
	FtpPassword string `json:"ftpPassword"`

	// path to build direcory on FTP server
	FtpBuildPaths []string `json:"ftpBuildPaths"`

	// full path to buils dir
	BuildsPath string `json:"buildsPath"`

	// last | all. Default: last
	Mode string `json:"Mode"`
}

func main() {
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "conf.json"
	}
	log.Print(configPath)

	var config Configuration = getConfig(configPath)
	var ftpConnection *goftp.FTP
	var builds []os.FileInfo = GetAllBuilds(config.BuildsPath)

	ftpConnection = ftp.GetFtpConnection(config.FtpHost, config.FtpUser, config.FtpPassword)
	var currentPath string = ftp.GetCurrentPath(ftpConnection)

	log.Print("Connection to ftp ", config.FtpHost, " successfully established")
	log.Print("Current Path is", currentPath)

	log.Print(config.FtpBuildPaths)

	if config.Mode == "all" {
		log.Print("Mode: all builds.")
		copyAllBuilds(builds, config, ftpConnection)
	} else {
		var lastBuild os.FileInfo = builds[len(builds)-1]
		log.Print("Mode: last build.\nLast build is", lastBuild.Name())
		copyBuild(lastBuild, config, ftpConnection)
	}

	log.Print("Builds transfered successfully. Exit.")
	ftpConnection.Quit()
}

func copyAllBuilds(builds []os.FileInfo, config Configuration, ftpConnection *goftp.FTP) {
	for i := 0; i < len(builds); i++ {
		log.Print(builds[i].Name())
		copyBuild(builds[i], config, ftpConnection)
	}
}

func copyBuild(build os.FileInfo, config Configuration, ftpConnection *goftp.FTP) {
	for i := 0; i < len(config.FtpBuildPaths); i++ {
		ftp.MakeBuildDir(ftpConnection, config.FtpBuildPaths[i]+"/"+build.Name())
		var path = config.FtpBuildPaths[i]
		ModeToFTP(config, build, ftpConnection, config.BuildsPath, path, 0)
	}
}

func getConfig(path string) Configuration {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("error:", err)
	}
	return configuration
}

func GetAllBuilds(path string) []os.FileInfo {
	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(path); err != nil {
		log.Fatal(err)
	}
	return files
}

func ModeToFTP(config Configuration, build os.FileInfo, ftpConnection *goftp.FTP, buildsPath string, ftpBuildPaths string, isRec int) {
	log.Printf("%s\n%s\n\n", buildsPath, ftpBuildPaths)
	var buildFiles []os.FileInfo
	if isRec == 0 {
		buildFiles = GetBuildFiles(buildsPath + "/" + build.Name())
	} else {
		buildFiles = GetBuildFiles(buildsPath)
	}
	for i := 0; i < len(buildFiles); i++ {
		var (
			fileInfo os.FileInfo = buildFiles[i]
			file     *os.File
			err      error
			fpath    string
			ftpfpath string
		)
		if fileInfo.IsDir() {
			rBuildsPath := buildsPath + "/" + build.Name() + "/" + fileInfo.Name() + "/"
			rFTPPath := ftpBuildPaths + "/" + strings.Replace(buildsPath, config.BuildsPath, "", -1) + "/" + build.Name() + "/" + fileInfo.Name()
			ModeToFTP(config, build, ftpConnection, rBuildsPath, rFTPPath, 1)
			continue
		}
		if isRec == 0 {
			fpath = buildsPath + "/" + build.Name() + "/" + fileInfo.Name()
		} else {
			fpath = buildsPath + "/" + fileInfo.Name()
		}
		if file, err = os.Open(fpath); err != nil {
			log.Fatal(err)
		}
		var reader io.Reader = (*os.File)(file)
		if isRec == 0 {
			ftpfpath = ftpBuildPaths + "/" + build.Name() + "/" + fileInfo.Name()
		} else {
			ftpfpath = ftpBuildPaths + "/" + fileInfo.Name()
		}
		if err := ftpConnection.Stor(ftpfpath, reader); err != nil {
			log.Fatal(err)
		}
		log.Print(file.Name(), " copied to ", ftpfpath)
	}
}

func GetBuildFiles(path string) []os.FileInfo {
	var err error
	var filesAndFolders []os.FileInfo
	if filesAndFolders, err = ioutil.ReadDir(path); err != nil {
		log.Fatal(err)
	}

	var files = make([]os.FileInfo, 0)
	for i := 0; i < len(filesAndFolders); i++ {
		files = append(files, filesAndFolders[i])
	}
	return files
}
