package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/parthvinchhi/jitapi"
)

func Handler(c *gin.Context) {
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

	date := c.PostForm("date")
	parseDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid date format")
		return
	}

	formattedDate := parseDate.Format("2006-01-02")

	query := "select * from HARDWARE_ISSUES_INFOS WHERE created_at > '" + fmt.Sprintf("%v", formattedDate) + " 03:10:18' ORDER BY created_at DESC;"

	data, err := connStr.Query(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.File(date + ".csv")
}
