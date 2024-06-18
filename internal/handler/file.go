package handler

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileV1Handler struct {
	fileUseCase fileUseCase
}

func NewFileV1Handler(fileUseCase fileUseCase) *FileV1Handler {
	return &FileV1Handler{
		fileUseCase: fileUseCase,
	}
}

func (fh *FileV1Handler) Create(c *gin.Context) {
	type FileForm struct {
		Name string `json:"name"`
		Size uint64 `json:"size"`
		Type string `json:"type"`
	}

	var ff FileForm
	err := c.BindJSON(&ff)
	if err != nil {
		slog.Error("BindJSON", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	slog.Info("FileV1Handler Create", "ff", ff)

	file, err := fh.fileUseCase.CreateTmp(ff.Name, ff.Type, ff.Size)
	if err != nil {
		slog.Error("BindJSON", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusCreated, file)
}

func (fh *FileV1Handler) GetPresignUrl(c *gin.Context) {
	slog.Info("GetPresignUrl", "key", c.Param("key"))

	PresignUrl, err := fh.fileUseCase.GetPresignUrl(c.Param("key"))
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"url": PresignUrl})
}

func (fh *FileV1Handler) SetS3Key(c *gin.Context) {
	file, err := fh.fileUseCase.SetS3Key(c.Param("fileID"), c.Param("s3Key"))
	if err != nil {
		slog.Error("BindJSON", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (fh *FileV1Handler) Delete(c *gin.Context) {
	err := fh.fileUseCase.DeleteTmp(c.Param("fileID"))

	if err != nil {
		slog.Error("BindJSON", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
