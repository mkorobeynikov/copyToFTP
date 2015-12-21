package main
import "fmt"
import "encoding/json"
import (
	"os"
	"io/ioutil"
	"github.com/mkorobeynikov/copyToFTP/ftp"
	"github.com/dutchcoders/goftp"
	"io"
)

type Configuration struct {
	FtpHost       string `json:"ftpHost"`
	FtpUser       string `json:"ftpUser"`
	FtpPassword   string `json:"ftpPassword"`
	FtpBuildPaths []string `json:"ftpBuildPaths"`
	BuildsPath    string `json:"buildsPath"`
}

func main() {
	configPath := os.Args[1]
	fmt.Println(configPath)

	var config Configuration = getConfig(configPath)
	var ftpConnection *goftp.FTP
	var builds []os.FileInfo = GetAllBuilds(config.BuildsPath)

	var lastBuild os.FileInfo = builds[len(builds) - 1]
	fmt.Println("Last build is", lastBuild.Name())

	ftpConnection = ftp.GetFtpConnection(config.FtpHost, config.FtpUser, config.FtpPassword)
	var currentPath string = ftp.GetCurrentPath(ftpConnection)

	fmt.Println("Connection to ftp", config.FtpHost, "successfully established")
	fmt.Println("Current Path is", currentPath)

	fmt.Println(config.FtpBuildPaths)

	for i := 0; i < len(config.FtpBuildPaths); i++ {
		ftp.MakeBuildDir(ftpConnection, config.FtpBuildPaths[i] + "/" + lastBuild.Name())
		var path = config.FtpBuildPaths[i]
		CopyBuildToFTP(config, lastBuild, ftpConnection, config.BuildsPath, path)
	}

	fmt.Println("Builds transfered successfully. Exit.")
	ftpConnection.Quit()
}

func getConfig(path string) Configuration {
	file, _ := os.Open(path)
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

func CopyBuildToFTP(config Configuration, lastBuild os.FileInfo, ftpConnection *goftp.FTP, buildsPath string, ftpBuildPaths string) {
	var buildFiles []os.FileInfo = GetBuildFiles(buildsPath + "/" + lastBuild.Name())
	for i := 0; i < len(buildFiles); i++ {
		var fileInfo os.FileInfo = buildFiles[i]
		var file *os.File
		var err error
		if file, err = os.Open(buildsPath + "/" + lastBuild.Name() + "/" + fileInfo.Name()); err != nil {
			panic(err)
		}
		var reader io.Reader = (*os.File)(file)
		if err := ftpConnection.Stor(ftpBuildPaths + "/" + lastBuild.Name() + "/" + fileInfo.Name(), reader); err != nil {
			panic(err)
		}
		fmt.Println(file.Name(), "copied to", ftpBuildPaths, "/", lastBuild.Name())
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
