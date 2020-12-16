package service

import (
	"bytes"
	"fmt"
	"github.com/ToolPackage/fse/common/tx"
	"github.com/ToolPackage/fse/common/utils"
	"github.com/ToolPackage/fse/server/storage"
	"log"
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

func auth(c *tx.Channel, _ *tx.Packet) {
	// TODO
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
		c.RespBadRequest("header filename missing")
		return
	}
	filename = utils.TrimWhitespaces(filename)
	if len(filename) == 0 {
		c.RespBadRequest("invalid filename")
		return
	}

	contentType, ok := p.Headers["contentType"]
	if !ok {
		c.RespBadRequest("header contentType missing")
		return
	}
	contentType = utils.TrimWhitespaces(contentType)
	if len(filename) == 0 {
		c.RespBadRequest("invalid contentType")
		return
	}

	file, err := storage.S.SaveFile(filename, contentType, bytes.NewReader(p.Content))
	if err != nil {
		c.RespInternalServerError(err.Error())
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
	c.RespAccepted(fileInfo)
}

func downloadFile(c *tx.Channel, p *tx.Packet) {
	fileId, ok := p.Headers["fileId"]
	if !ok {
		c.RespBadRequest("header fileId missing")
		return
	}
	fileId = utils.TrimWhitespaces(fileId)
	if len(fileId) == 0 {
		c.RespBadRequest("invalid fileId")
		return
	}

	file, ok := storage.S.GetFile(fileId)
	if !ok {
		c.RespNotFound("file not found")
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file.OpenStream())
	if err != nil {
		c.RespInternalServerError(err.Error())
		return
	}
	fmt.Println(buf.Len())
	c.RespOk(buf.Bytes())
}

func deleteFile(c *tx.Channel, p *tx.Packet) {
	fileId, ok := p.Headers["fileId"]
	if !ok {
		c.RespBadRequest("header fileId missing")
		return
	}
	fileId = utils.TrimWhitespaces(fileId)
	if len(fileId) == 0 {
		c.RespBadRequest("invalid fileId")
		return
	}

	if ok := storage.S.DeleteFile(fileId); ok {
		c.NewPacket(Resp).StatusCode(http.StatusOK).Emit()
	} else {
		c.RespInternalServerError("failed to delete file")
	}
}
