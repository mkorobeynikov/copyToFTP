package main
import "fmt"
import "encoding/json"
import (
	"os"
	"github.com/dutchcoders/goftp"
	"github.com/mkorobeynikov/copyToFTP/ftp"
)

func main() {
	var ftpConnection *goftp.FTP
	var config Configuration = getConfig()
	ftpConnection = ftp.GetFtpConnection(config.FtpPath, config.FtpUser, config.FtpPassword)

	var folders []string = ftp.GetFolders(ftpConnection, "")
	fmt.Println(folders)
}

type Configuration struct {
	FtpPath     string `json:"ftpPath"`
	FtpUser     string `json:"ftpUser"`
	FtpPassword string `json:"ftpPassword"`
}

func getConfig() Configuration {
	file, _ := os.Open("conf/conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
