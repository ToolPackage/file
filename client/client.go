package internal

import (
	"encoding/json"
	"fmt"
	ini "github.com/Luncert/go-ini"
	"github.com/ToolPackage/fse/common/tx"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

const (
	Auth     = "auth"
	List     = "list"
	Upload   = "upload"
	Download = "download"
	Delete   = "delete"
	Resp     = "resp"
)

const (
	configFileName    = ".fse-cli"
	curConnEnvVarName = "currentConnection"
)

var (
	configFilePath = getConfigFilePath()
)

type FileInfo struct {
	FileId      string     `json:"fileId"`
	FileName    string     `json:"fileName"`
	ContentType string     `json:"contentType"`
	CreatedAt   int64      `json:"createdAt"`
	FileSize    int64      `json:"fileSize"`
	Partitions  Partitions `json:"partitions"`
}
type PartitionId uint32
type Partitions []PartitionId

type FseClient struct {
	cfg     *ini.Config
	channel *tx.Channel
}

func NewClient() *FseClient {
	return &FseClient{cfg: loadConfig()}
}

func (f *FseClient) login() error {
	connName := f.cfg.Section("common").Variable(curConnEnvVarName)
	if connName == nil || connName.Type() != ini.StringType || len(connName.String()) == 0 {
		return fmt.Errorf("current connection not set, use peek command to set one")
	}

	creds, err := lookupCredential(connName.String())
	if err != nil {
		return err
	}
	cred := creds[0]

	if err := f.connect(cred.UserName); err != nil {
		return err
	}

	// send auth packet
	f.channel.NewPacket(Auth).Body(string(cred.CredentialBlob)).Emit()
	res := f.channel.RecvPacket()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed: %s", string(res.Content))
	}

	return nil
}

func (f *FseClient) connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	f.channel = tx.NewChannel(conn, conn)
	return nil
}

func (f *FseClient) setConnection(connName string) {
	f.cfg.CreateIfAbsent("common").
		CreateIfAbsent(curConnEnvVarName).
		SetValue(ini.NewStringValue(connName))
	saveConfig(f.cfg)
}

func (f *FseClient) listFiles(prefixFilter string) ([]FileInfo, error) {
	f.channel.NewPacket(List).
		Header("prefixFilter", prefixFilter).
		Emit()
	res := f.channel.RecvPacket()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", string(res.Content))
	}

	data := make([]FileInfo, 0)
	err := json.Unmarshal(res.Content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f *FseClient) uploadFile(filename string, contentType string, content []byte) (*FileInfo, error) {
	f.channel.NewPacket(Upload).
		Header("filename", filename).
		Header("contentType", contentType).
		Body(content).
		Emit()
	res := f.channel.RecvPacket()
	if res.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("request failed: %s", string(res.Content))
	}

	data := &FileInfo{}
	err := json.Unmarshal(res.Content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f *FseClient) close() {

}

func loadConfig() *ini.Config {
	cfg, err := ini.ReadConfigFile(configFilePath)
	if err != nil {
		panic(err)
	}
	return cfg
}

func saveConfig(cfg *ini.Config) {
	if err := ini.WriteConfigFile(configFilePath, cfg); err != nil {
		panic(err)
	}
}

func getConfigFilePath() string {
	p := filepath.Join(getUserHomeDir(), configFileName)
	if _, err := os.Stat(p); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		err = ioutil.WriteFile(p, []byte{}, 0644)
		if err != nil {
			panic(err)
		}
	}

	return p
}
