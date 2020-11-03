package api

import (
	"context"
	"fmt"
	"github.com/Luncert/slog"
	"github.com/ToolPackage/fse/server/config"
	"github.com/ToolPackage/fse/server/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"path"
)

type FileInfo struct {
	FileId   string `bson:"FileId"`
	FileName string `bson:"FileName"`
	FilePath string `bson:"FilePath"`
	FileSize int64  `bson:"FileSize"`
}

func GetFilesList(ctx *gin.Context) {
	findOps := &options.FindOptions{}
	findOps.SetSort(bson.M{"_id": -1})
	cursor, err := db.MongoDb.Collection(config.FileInfoMongoCol).Find(
		context.TODO(),
		bson.M{},
		findOps,
	)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	results := make([]FileInfo, 0)
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
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
	fileSize := fileHeader.Size
	fileExt := path.Ext(fileName)
	uid, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fileId := uid.String()
	filePath := fmt.Sprintf("/upload/%s%s", fileId, fileExt)

	// read file from ctx and save to
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("failed to read uploaded file", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if file != nil {
		// TODO: read file and save to fs
	}

	// save file metadata to mongo
	_, err = db.MongoDb.Collection(config.FileInfoMongoCol).InsertOne(
		context.TODO(),
		&FileInfo{
			FileId:   fileId,
			FileName: fileName,
			FilePath: filePath,
			FileSize: fileSize,
		},
	)

	if err != nil {
		log.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//service.PersistenceServiceIns.SaveFile()

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "FileId": fileId})
}

func GetFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")
	fmt.Println(fileId)

	var info FileInfo
	err := db.MongoDb.Collection(config.FileInfoMongoCol).FindOne(
		context.TODO(),
		bson.M{
			"FileId": fileId,
		},
	).Decode(&info)
	if err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			ctx.AbortWithStatus(http.StatusNotFound)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", info.FileName))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(info.FilePath)
}

func DeleteFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")

	// query file metadata
	var info FileInfo
	err := db.MongoDb.Collection(config.FileInfoMongoCol).FindOne(
		context.TODO(),
		bson.M{
			"FileId": fileId,
		},
	).Decode(&info)
	if err != nil {
		log.Error(err)
		if err == mongo.ErrNoDocuments {
			ctx.AbortWithStatus(http.StatusNotFound)
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	// TODO: involve storage service to delete file
	err = os.Remove(info.FilePath)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete file metadata
	res, err := db.MongoDb.Collection(config.FileInfoMongoCol).DeleteOne(context.TODO(), bson.M{"FileId": fileId})
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
