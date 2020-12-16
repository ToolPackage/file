package service

import (
	"bytes"
	"github.com/ToolPackage/fse/storage"
	"github.com/ToolPackage/fse/tx"
	"log"
	"net/http"
	"strings"
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

func auth(c *tx.Channel, _ *tx.Packet) {
	log.Println("client auth passed")
	c.NewPacket(Resp).StatusCode(http.StatusOK).Emit()
}

func listFiles(c *tx.Channel, p *tx.Packet) {
	prefixFilter := ""
	if v, ok := p.Headers["prefixFilter"]; ok {
		prefixFilter = v
	}

	files := storage.S.GetAllFiles(prefixFilter)
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
	filename, ok := p.Headers["filename"]
	if !ok {
		c.NewPacket(Resp).StatusCode(http.StatusBadRequest).Body("header filename missing").Emit()
		return
	}
	filename = strings.Trim(filename, " \t\r\n")
	if len(filename) == 0 {
		c.NewPacket(Resp).StatusCode(http.StatusBadRequest).Body("invalid filename").Emit()
		return
	}

	contentType, ok := p.Headers["contentType"]
	if !ok {
		c.NewPacket(Resp).StatusCode(http.StatusBadRequest).Body("header contentType missing").Emit()
		return
	}
	contentType = strings.Trim(contentType, " \t\r\n")
	if len(filename) == 0 {
		c.NewPacket(Resp).StatusCode(http.StatusBadRequest).Body("invalid contentType").Emit()
		return
	}

	file, err := storage.S.SaveFile(filename, contentType, bytes.NewReader(p.Content))
	if err != nil {
		c.NewPacket(Resp).StatusCode(http.StatusInternalServerError).Body(err.Error()).Emit()
		return
	}

	fileInfo := &FileInfo{
		FileId:      file.Id,
		FileName:    file.Name,
		ContentType: file.ContentType,
		CreatedAt:   file.CreatedAt,
		FileSize:    int64(file.Size),
		Partitions:  file.Partitions,
	}
	c.NewPacket(Resp).StatusCode(http.StatusAccepted).Body(fileInfo).Emit()
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
