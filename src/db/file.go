package db

import (
	"database/sql"
	mydb "db/mysql"
	"fmt"
)

type TableFile struct {
	FileName sql.NullString
	FileSize sql.NullInt64
	FilePath sql.NullString
	Hash     string
}

// OnFileUploadFinished : file uploaded, store meta to db
func OnFileUploadFinished(filename string, filesize int64, filepath string, hash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tb_file(`filename`, `filesize`, `filepath`, `hash`, `status`) " +
			"values(?,?,?,?,0)")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(filename, filesize, filepath, hash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if nums, err := res.RowsAffected(); err == nil {
		if nums <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", hash)
			return false
		}
		return true
	}
	return false
}

// GetFileMeta : througn filename search fileMeta
func GetFileMeta(filename string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select filename, filesize, filepath, hash from tb_file" +
			"where filename = ? and status = 1 limit 1")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(filename).Scan(
		&tfile.FileName, &tfile.FileSize, &tfile.FilePath, &tfile.Hash)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}

// IsFileUploaded : check if hash already exists
func IsFileUploaded(hash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select 1 from tb_file where hash = ? and status = 1 limit 1")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	rows, err := stmt.Query(hash)
	if err != nil {
		return false
	} else if rows == nil || !rows.Next() {
		return false
	}
	return true
}

// GetFileMetaLists : get lists of recent file
func GetFileMetaLists(limit int) ([]TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"seletct filename, filesize, filepath, hash from tb_file" +
			"where status=1 limt ?")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileName, &tfile.FileSize, &tfile.FilePath, &tfile.Hash)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	fmt.Println(len(tfiles))
	return tfiles, nil
}

// OnFileRemoved : set the file meta col's status = 2 (unvaild)
func OnFileRemoved(filename string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tb_file set status = 2 where filename = ? and status = 1 limit 1")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(filename)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if nums, err := res.RowsAffected(); err == nil {
		if nums <= 0 {
			fmt.Println("File Remove from tb_file err:", err.Error())
			return false
		}
		return true
	}
	return false
}

// OnFileMetaUpdate : update file meta
func OnFileMetaUpdate(oldfilename, newfilename string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tb_file set filename=? where filename = ? limit 1")
	if err != nil {
		fmt.Println("Failded to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(newfilename, oldfilename)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if nums, err := res.RowsAffected(); err == nil {
		if nums <= 0 {
			fmt.Printf("File Remove from tb_file err:", err.Error())
			return false
		}
		return true
	}
	return false
}