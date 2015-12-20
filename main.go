package main
import "fmt"
import "encoding/json"
import (
	"os"
/*	"github.com/dutchcoders/goftp"
	"github.com/mkorobeynikov/copyToFTP/ftp"*/
	"io/ioutil"
	"github.com/mkorobeynikov/copyToFTP/ftp"
	"github.com/dutchcoders/goftp"
	"io"
)

type Configuration struct {
	FtpHost      string `json:"ftpHost"`
	FtpUser      string `json:"ftpUser"`
	FtpPassword  string `json:"ftpPassword"`
	FtpBuildPath string `json:"ftpBuildPath"`
	BuildsPath   string `json:"buildsPath"`
}

func main() {
	var config Configuration = getConfig()
	var ftpConnection *goftp.FTP
	var builds []os.FileInfo = GetAllBuilds(config.BuildsPath)

	var lastBuild os.FileInfo = builds[len(builds) - 1]
	fmt.Println("Last build is", lastBuild.Name())

	ftpConnection = ftp.GetFtpConnection(config.FtpHost, config.FtpUser, config.FtpPassword)
	var currentPath string = ftp.GetCurrentPath(ftpConnection)

	fmt.Println("Connection to ftp", config.FtpHost, "successfully established")
	fmt.Println("Current Path is", currentPath)

	ftp.MakeBuildDir(ftpConnection, config.FtpBuildPath + "/" + lastBuild.Name())

	CopyBuildToFTP(config, lastBuild, ftpConnection)

	ftpConnection.Quit()
}

func getConfig() Configuration {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func GetAllBuilds(path string) []os.FileInfo {
	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(path); err != nil {
		panic(err)
	}
	return files
}

func CopyBuildToFTP(config Configuration, lastBuild os.FileInfo, ftpConnection *goftp.FTP) {
	var buildFiles []os.FileInfo = GetBuildFiles(config.BuildsPath + "/" + lastBuild.Name())
	for i := 0; i < len(buildFiles); i++ {
		var fileInfo os.FileInfo = buildFiles[i]
		var file *os.File
		var err error
		if file, err = os.Open(config.BuildsPath + "/" + lastBuild.Name() + "/" + fileInfo.Name()); err != nil {
			panic(err)
		}
		var reader io.Reader = (*os.File)(file)
		if err := ftpConnection.Stor(config.FtpBuildPath + "/" + lastBuild.Name() + "/" + fileInfo.Name(), reader); err != nil {
			panic(err)
		}
		fmt.Println(file.Name(), "copied to", config.FtpBuildPath, "/", lastBuild.Name())
	}
}

func GetBuildFiles(path string) []os.FileInfo {
	var err error
	var filesAndFolders []os.FileInfo
	if filesAndFolders, err = ioutil.ReadDir(path); err != nil {
		panic(err)
	}

	var files = make([]os.FileInfo, 0)
	for i := 0; i < len(filesAndFolders); i++ {
		var file os.FileInfo = filesAndFolders[i]
		if !file.IsDir() {
			files = append(files, file)
		}
	}
	return files
}
