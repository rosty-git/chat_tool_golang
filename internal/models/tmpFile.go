package models

type TmpFile struct {
	BaseModel

	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  uint64 `json:"size"`
	S3Key string `json:"s3_key"`
}
