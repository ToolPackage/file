package service

import "github.com/ToolPackage/fse/common/tx"

func Init() {
	tx.Register(Auth, auth)
	tx.Register(List, listFiles)
	tx.Register(Upload, uploadFile)
	tx.Register(Download, downloadFile)
	tx.Register(Delete, deleteFile)
}
