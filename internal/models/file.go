package models

type File struct {
	BaseModel

	Name   string `json:"name"`
	Type   string `json:"type"`
	Size   uint64 `json:"size"`
	S3Key  string `json:"s3Key"`
	PostID string `json:"post_id" gorm:"size:191"`
}
