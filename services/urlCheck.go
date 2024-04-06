package services

import (
	"github.com/Oik17/FamPay-BE/database"
	"github.com/Oik17/FamPay-BE/models"
)

func CheckUrlInDB(url string) (bool, error) {
	db := database.DB.Db
	var video []models.Video

	query := `SELECT * FROM video WHERE url=$1`
	err := db.Select(&video, query, url)
	if err != nil {
		return false, err
	}

	if len(video) == 0 {
		return true, nil
	}

	return false, nil
}
