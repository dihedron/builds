package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v1"
)

// Deployment is the object representing a deployed
// build.
type Deployment struct {
	ID          string `db:"id" form:"id" json:"id"`
	Group       string `db:"group" form:"group" json:"group" binding:"required"`
	Artifact    string `db:"artifact" form:"artifact" json:"artifact" binding:"required"`
	Vers        string `db:"version" form:"version" json:"version" binding:"required"`
	Environment string `db:"environment" form:"environment" json:"environment" binding:"required"`
	Date        string `db:"date" form:"date" json:"date"`
	Snapshot    bool   `db:"snapshot"`
	Success     bool   `db:"success" form:"success" json:"success"`
	System      string `db:"system" form:"system" json:"date"`
	Builder     string `db:"builder" form:"builder" json:"builder"`
}

// InitDB initialises the database and the DB mapper.
func InitDB() (*gorp.DbMap, error) {

	// load the SQLITE3 DB driver
	db, err := sql.Open("sqlite3", "./builds.db")
	if err != nil {
		log.Printf("error loading database driver %q\n", err)
		return nil, err
	}

	// open and test the connection
	if err = db.Ping(); err != nil {
		log.Printf("error connecting to database %q\n", err)
		defer db.Close()
		return nil, err
	}

	// create the ORM mapper
	mapper := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// let the mapper initialise the table for us
	// add a key as an hash

	mapper.AddTableWithName(Deployment{}, "deployments").SetKeys(false, "id")

	if err = mapper.CreateTablesIfNotExists(); err != nil {
		log.Printf("error creating database table %q\n", err)
		return nil, err
	}
	return mapper, nil
}

var mapper, _ = InitDB()

// ReadDeployments returns a list of all known deployments.
func ReadDeployments(context *gin.Context) {
	var deployments []Deployment
	if _, err := mapper.Select(&deployments, "select * from deployments order by id"); err != nil {
		log.Printf("error selecting all deployments from the database %q\n", err)
		context.String(http.StatusOK, "[]")
		return
	}
	content := gin.H{}
	for k, v := range deployments {
		content[strconv.Itoa(k)] = v
	}
	context.JSON(200, content)
}

// DeleteDeployments deletes all known deployments.
func DeleteDeployments(context *gin.Context) {

}

// CreateDeployment inserts a new deployment into the database.
func CreateDeployment(context *gin.Context) {
	var deployment Deployment

	context.Bind(&deployment) // This will infer what binder to use depending on the content-type header.

	hasher := sha256.New()
	hasher.Write([]byte(deployment.Group + deployment.Artifact + deployment.Environment + deployment.Environment))
	deployment.ID = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	deployment.Date = time.Now().Format("2006-01-02@15:04:05")
	deployment.Snapshot = strings.HasSuffix(deployment.Vers, "SNAPSHOT")

	if err := mapper.Insert(&deployment); err != nil {
		log.Printf("error inserting new deployment %q\n", err)
	}
	log.Printf("inserted with id: %q\n", deployment.ID)

	/*
		content := gin.H{
			"result":  "Success",
			"title":   article.Title,
			"content": article.Content,
		}
	*/
	context.JSON(201, deployment)
	/*
		} else {
			c.JSON(500, gin.H{"result": "An error occured"})
		}
	*/
}

// ReadDeployment returns the representation of a deployment.
func ReadDeployment(context *gin.Context) {

}

// ReplaceDeployment completely replaces a given deployment with
// user-supplied information.
func ReplaceDeployment(context *gin.Context) {

}

// UpdateDeployment patches the given deployment, updating its
// contents with user-supplied data.
func UpdateDeployment(context *gin.Context) {

}

// DeleteDeployment deletes a deployment from the database.
func DeleteDeployment(context *gin.Context) {

}

func main() {
	router := gin.Default()

	router.GET("/deployments", ReadDeployments)
	router.DELETE("/deployments", DeleteDeployments)
	router.POST("/deployments", CreateDeployment)
	router.GET("/deployments/:id", ReadDeployment)
	router.PUT("/deployments/:id", ReplaceDeployment)
	router.PATCH("/deployments/:id", UpdateDeployment)
	router.DELETE("/deployments/:id", DeleteDeployment)

	router.Run(":8080")
}
