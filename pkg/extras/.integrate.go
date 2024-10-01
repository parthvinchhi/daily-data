package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/parthvinchhi/jitapi"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// dbConfig give the configuration details
	dbConfig := jitapi.DbConfig{
		DBType:     os.Getenv("DB_TYPE"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSslMode:  os.Getenv("DB_SSLMODE"),
	}

	connStr := jitapi.Postgres{
		Config: dbConfig,
	}

	if err := connStr.Connect(); err != nil {
		log.Fatal(err)
	}

	// query := "select * from HARDWARE_ISSUES_INFOS WHERE created_at BETWEEN '2024-09-17 03:10:18' AND '2024-09-18 03:10:18' ORDER BY created_at DESC;"
	query, date := jitapi.GetQuery()

	data, err := connStr.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var v jitapi.Variables
	v.CountDataByString(data)

	var h jitapi.Helper
	h.WriteCustomError(data)
	h.VideoSavedFilter(data)
	h.ZeroFramesFilter(data)
	h.GetMissedVIDs(data)

	datatobewritten := jitapi.SideBySideData(v.WriteCountToCsv(), h.Result)

	err = jitapi.WriteDataToCsv(date+".csv", datatobewritten)
	if err != nil {
		log.Fatal(err)
	}

	// err = jitapi.WriteDataToCsv("test.csv", datatobewritten)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
