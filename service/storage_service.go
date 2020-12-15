package service

import (
	log "github.com/Luncert/slog"
	"github.com/ToolPackage/fse/storage"
	"github.com/ToolPackage/fse/tx"
	"net/http"
)

// actions
const (
	Auth     = "auth"
	List     = "list"
	Upload   = "upload"
	Download = "download"
	Delete   = "delete"
	Resp     = "resp"
)

type FileInfo struct {
	FileId      string             `json:"fileId"`
	FileName    string             `json:"fileName"`
	ContentType string             `json:"contentType"`
	CreatedAt   int64              `json:"createdAt"`
	FileSize    int64              `json:"fileSize"`
	Partitions  storage.Partitions `json:"partitions"`
}

type FileDetail struct{}

func auth(c *tx.Channel, _ *tx.Packet) {
	log.Info("client auth passed")
	c.NewPacket(Resp).StatusCode(http.StatusOK).Emit()
}

func listFiles(c *tx.Channel, _ *tx.Packet) {
	files := storage.S.GetAllFiles()
	results := make([]FileInfo, len(files))
	for idx, file := range files {
		results[idx].FileId = file.Id
		results[idx].FileName = file.Name
		results[idx].ContentType = file.ContentType
		results[idx].CreatedAt = file.CreatedAt
		results[idx].FileSize = int64(file.Size)
		results[idx].Partitions = file.Partitions
	}
	// send response packet
	c.NewPacket(Resp).StatusCode(http.StatusOK).Body(results).Emit()
}

func uploadFile(c *tx.Channel, p *tx.Packet) {
	c.NewPacket(Resp).StatusCode(http.StatusBadRequest).Emit()
}

func downloadFile(c *tx.Channel, p *tx.Packet) {
	c.NewPacket(Resp).StatusCode(http.StatusBadRequest).Emit()
}

func deleteFile(c *tx.Channel, p *tx.Packet) {
	id := string(p.Content)
	if ok := storage.S.DeleteFile(id); ok {
		c.NewPacket(Resp).StatusCode(http.StatusOK).Emit()
	} else {
		c.NewPacket(Resp).StatusCode(http.StatusInternalServerError).Emit()
	}
}

func init() {
	tx.Register(Auth, auth)
	tx.Register(List, listFiles)
	tx.Register(Upload, uploadFile)
	tx.Register(Download, downloadFile)
	tx.Register(Delete, deleteFile)
}
