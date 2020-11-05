package api

import (
	"bytes"
	"fmt"
	"github.com/Luncert/slog"
	"github.com/ToolPackage/fse/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type FileInfo struct {
	FileId   string `bson:"FileId"`
	FileName string `bson:"FileName"`
	FileSize int64  `bson:"FileSize"`
}

func GetFilesList(ctx *gin.Context) {
	files := storage.S.GetAllFiles()
	results := make([]FileInfo, len(files))
	for idx, file := range files {
		results[idx].FileId = file.Id
		results[idx].FileName = file.Name
		results[idx].FileSize = int64(file.Size)
	}
	ctx.JSON(http.StatusOK, results)
}

func PostFile(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// build file metadata
	fileName := fileHeader.Filename
	reader, err := fileHeader.Open()
	if err != nil {
		log.Error("failed to read uploaded file", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	file, err := storage.S.SaveFile(fileName, "", reader)
	if err != nil {
		log.Error("failed to read uploaded file", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fileInfo := &FileInfo{
		FileId:   file.Id,
		FileName: file.Name,
		FileSize: int64(file.Size),
	}

	ctx.JSON(http.StatusOK, fileInfo)
}

func GetFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")
	file, ok := storage.S.GetFile(fileId)
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file.OpenStream()); err != nil {
		log.Error("failed to read file", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	ctx.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}

func DeleteFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")
	if err := storage.S.DeleteFile(fileId); err != nil {
		if err == os.ErrNotExist {
			ctx.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Error("failed to read file", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	ctx.JSON(http.StatusOK, "")
}
