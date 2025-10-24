package models

import "time"

type File struct {
	ID        interface{} `json:"id" bson:"_id,omitempty"`
	AlumniID  interface{} `json:"alumni_id" bson:"alumni_id"`
	FileName  string      `json:"file_name" bson:"file_name"`
	FileType  string      `json:"file_type" bson:"file_type"`
	FilePath  string      `json:"file_path" bson:"file_path"`
	FileSize  int64       `json:"file_size" bson:"file_size"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

type FileResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    *File  `json:"data,omitempty"`
}
