package database

import (
	"database/sql"
	"github.com/coder/flog"

	models "github.com/matttm/spoticli/spoticli-models"
)

func InsertFileMetaInfo(tx *sql.Tx, key_name, bucket_name string, file_type_cd, file_size int) {
	query := "INSERT INTO SPOTICLI_DB.FILE_META_INFO (key_name, bucket_name, file_type_cd, file_size) VALUES (?, ?, ?, ?);"
	_, err := tx.Exec(query, key_name, bucket_name, file_type_cd, file_size)
	if err != nil {
		flog.Errorf(err.Error())
	}
}
func SelectAllFileMetaInfo() []*models.FileMetaInfo {
	files := []*models.FileMetaInfo{}
	query := "SELECT id, key_name, bucket_name FROM SPOTICLI_DB.FILE_META_INFO;"
	rows, err := DB.Query(query)
	if err != nil {
		flog.Errorf(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		file := new(models.FileMetaInfo)
		if err := rows.Scan(&file.Id, &file.Key_name, &file.Bucket_name); err != nil {
			flog.Errorf(err.Error())
		}
		files = append(files, file)
	}
	return files
}
func SelectOneFileMetaInfo(id int) *models.FileMetaInfo {
	query := "SELECT key_name, bucket_name, file_size FROM SPOTICLI_DB.FILE_META_INFO WHERE ID = ?;"
	row := DB.QueryRow(query, id)
	if row == nil {
		panic("")
	}
	file := new(models.FileMetaInfo)
	if err := row.Scan(&file.Key_name, &file.Bucket_name, &file.File_size); err != nil {
		flog.Errorf(err.Error())
	}
	return file
}
