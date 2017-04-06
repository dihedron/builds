package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dihedron/builds/model"
	"github.com/gin-gonic/gin"
)

type Link struct {
	Relation string `json:"rel,omitempty"`
	URI      string `json:"href,omitempty"`
}

func GetProducts(c *gin.Context) {

	products, err := model.GetAllProducts()

	type ProductInfo struct {
		ID   string `json:"id,omitempty"`
		Self Link   `json:"_link,omitempty"`
	}

	if err == nil && len(products) > 0 {
		results := make([]ProductInfo, 0, len(products))
		for _, product := range products {
			results = append(results, ProductInfo{
				ID: product.ID,
				Self: Link{
					Relation: "self",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID,
				},
			})
		}
		c.JSON(http.StatusOK, gin.H{"products": results})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": []ProductInfo{}})
}

func GetProduct(c *gin.Context) {
	productId := c.Params.ByName("productId")

	if productId == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	product, err := model.GetProductByID(productId)

	type VersionInfo struct {
		ID   string `json:"id,omitempty"`
		Link Link   `json:"_link,omitempty"`
	}

	type ProductInfo struct {
		ID          string        `json:"id,omitempty"`
		Description string        `json:"description,omitempty"`
		Links       []Link        `json:"_links,omitempty"`
		Versions    []VersionInfo `json:"versions,omitempty"`
	}

	if err != nil || product.ID == "" {
		c.JSON(http.StatusOK, gin.H{"product": ProductInfo{}})
		return
	}

	var versions []VersionInfo
	if len(product.Versions) > 0 {
		versions = make([]VersionInfo, 0, 32)
		for _, version := range product.Versions {
			versions = append(versions, VersionInfo{
				ID: version.ID,
				Link: Link{
					Relation: "self",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID,
				},
			})
		}
	}

	result := ProductInfo{
		ID:          product.ID,
		Description: product.Description,
		Links: []Link{
			{
				Relation: "self",
				URI:      "http://" + c.Request.Host + "/products/" + product.ID,
			},
			{
				Relation: "collection",
				URI:      "http://" + c.Request.Host + "/products",
			},
			{
				Relation: "versions",
				URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions",
			},
		},
		Versions: versions,
	}

	c.JSON(http.StatusOK, gin.H{"product": result})
}

func GetVersions(c *gin.Context) {
	productId := c.Params.ByName("productId")

	if productId == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	product, err := model.GetProductByID(productId)

	// if no product found, return error code
	if err != nil || product.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	type DeploymentInfo struct {
		Order int  `json:"order"`
		Link  Link `json:"_link,omitempty"`
	}

	type VersionInfo struct {
		ID          string           `json:"id,omitempty"`
		Links       []Link           `json:"_links,omitempty"`
		Deployments []DeploymentInfo `json:"deployments,omitempty"`
	}

	var versions []VersionInfo

	if len(product.Versions) > 0 {
		versions = make([]VersionInfo, 0, 32)
		for _, version := range product.Versions {

			var deployments []DeploymentInfo
			if len(version.Deployments) > 0 {
				deployments = make([]DeploymentInfo, 0, 32)
				for _, deployment := range version.Deployments {
					deployments = append(deployments, DeploymentInfo{
						Order: deployment.Order,
						Link: Link{
							Relation: "self",
							URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments/" + strconv.Itoa(deployment.Order),
						},
					})
				}
			}
			versions = append(versions, VersionInfo{
				ID: version.ID,
				Links: []Link{
					{
						Relation: "self",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID,
					},
					{
						Relation: "collection",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions",
					},
					{
						Relation: "product",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID,
					},
					{
						Relation: "deployments",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments",
					},
				},
				Deployments: deployments,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func GetVersion(c *gin.Context) {
	productId := c.Params.ByName("productId")
	versionId := c.Params.ByName("versionId")

	if productId == "" || versionId == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	product, err := model.GetProductByID(productId)

	// if no product found, return error code
	if err != nil || product.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	version, err := model.GetVersionByID(productId, versionId)
	if err != nil || version.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	type DeploymentInfo struct {
		Order int  `json:"order"`
		Link  Link `json:"_link,omitempty"`
	}

	type VersionInfo struct {
		ID          string           `json:"id,omitempty"`
		Description string           `json:"description,omitempty"`
		Links       []Link           `json:"_links,omitempty"`
		Deployments []DeploymentInfo `json:"deployments,omitempty"`
	}

	var deployments []DeploymentInfo
	if len(version.Deployments) > 0 {
		deployments = make([]DeploymentInfo, 0, 32)
		for _, deployment := range version.Deployments {
			deployments = append(deployments, DeploymentInfo{
				Order: deployment.Order,
				Link: Link{
					Relation: "self",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments/" + strconv.Itoa(deployment.Order),
				},
			})
		}
	}
	result := VersionInfo{
		ID:          version.ID,
		Description: version.Description,
		Links: []Link{
			{
				Relation: "self",
				URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID,
			},
			{
				Relation: "collection",
				URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions",
			},
			{
				Relation: "product",
				URI:      "http://" + c.Request.Host + "/products/" + product.ID,
			},
			{
				Relation: "deployments",
				URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments",
			},
		},
		Deployments: deployments,
	}

	c.JSON(http.StatusOK, gin.H{"version": result})
}

func GetDeployments(c *gin.Context) {
	productId := c.Params.ByName("productId")
	versionId := c.Params.ByName("versionId")

	if productId == "" || versionId == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	product, err := model.GetProductByID(productId)

	// if no product found, return error code
	if err != nil || product.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	version, err := model.GetVersionByID(productId, versionId)
	if err != nil || version.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	type DeploymentInfo struct {
		Order int    `json:"order"`
		Links []Link `json:"_links,omitempty"`
	}

	var deployments []DeploymentInfo
	if len(version.Deployments) > 0 {
		deployments = make([]DeploymentInfo, 0, len(version.Deployments))
	}

	for _, deployment := range version.Deployments {
		deployments = append(deployments, DeploymentInfo{
			Order: deployment.Order,
			Links: []Link{
				{
					Relation: "self",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments/" + strconv.Itoa(deployment.Order),
				},
				{
					Relation: "collection",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments",
				},
				{
					Relation: "version",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID,
				},
				{
					Relation: "product",
					URI:      "http://" + c.Request.Host + "/products/" + product.ID,
				},
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"deployments": deployments})
}

func GetDeployment(c *gin.Context) {
	productId := c.Params.ByName("productId")
	versionId := c.Params.ByName("versionId")
	deploymentId := c.Params.ByName("deploymentId")

	if productId == "" || versionId == "" || deploymentId == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	product, err := model.GetProductByID(productId)

	// if no product found, return error code
	if err != nil || product.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	version, err := model.GetVersionByID(productId, versionId)
	if err != nil || version.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	deployOrder, err := strconv.Atoi(deploymentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for _, deployment := range version.Deployments {
		if deployment.Order == deployOrder {
			// found
			type DeploymentInfo struct {
				Order       int          `json:"order"`
				Environment string       `json:"environment,omitempty"`
				Status      model.Status `json:"status,omitempty"`
				GrantedBy   string       `json:"grantedBy,omitempty"`
				Timestamp   time.Time    `json:"timestamp,omitempty"`
				Links       []Link       `json:"_links,omitempty"`
			}

			result := DeploymentInfo{
				Order:       deployment.Order,
				Environment: deployment.Environment,
				Status:      deployment.Status,
				GrantedBy:   deployment.GrantedBy,
				Timestamp:   deployment.Timestamp,
				Links: []Link{
					{
						Relation: "self",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments/" + strconv.Itoa(deployment.Order),
					},
					{
						Relation: "collection",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID + "/deployments",
					},
					{
						Relation: "version",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID + "/versions/" + version.ID,
					},
					{
						Relation: "product",
						URI:      "http://" + c.Request.Host + "/products/" + product.ID,
					},
				},
			}

			c.JSON(http.StatusOK, gin.H{"deployment": result})
			return
		}
	}
	c.JSON(http.StatusNotFound, nil)
}

func ApproveDeployment(c *gin.Context) {
	productId := c.Params.ByName("productId")
	versionId := c.Params.ByName("versionId")
	deploymentId := c.Params.ByName("deploymentId")

	if productId == "" || versionId == "" || deploymentId == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	product, err := model.GetProductByID(productId)

	// if no product found, return error code
	if err != nil || product.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	version, err := model.GetVersionByID(productId, versionId)
	if err != nil || version.ID == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	deployOrder, err := strconv.Atoi(deploymentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for _, deployment := range version.Deployments {
		if deployment.Order == deployOrder {
			deployment.Status = model.GRANTED
			deployment.GrantedBy = "d093154" // TODO: use remote user for authenticated requests
			deployment.Timestamp = time.Now()

			model.UpdateDeployment(productId, versionId, deployment)

			c.JSON(http.StatusAccepted, nil)
			return
		}
	}
	c.JSON(http.StatusNotFound, nil)

}

func main() {

	router := gin.Default()
	router.GET("/products", GetProducts)
	router.GET("/products/:productId", GetProduct)
	router.GET("/products/:productId/versions", GetVersions)
	router.GET("/products/:productId/versions/:versionId", GetVersion)
	router.GET("/products/:productId/versions/:versionId/deployments", GetDeployments)
	router.GET("/products/:productId/versions/:versionId/deployments/:deploymentId", GetDeployment)

	router.GET("/products/:productId/versions/:versionId/deployments/:deploymentId/approve")

	/*
		// list all builds
		router.GET("/builds", func(c *gin.Context) {
			c.String(200, "pong")
		})

		// Ping test
		router.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		// Get user value
		router.GET("/user/:name", func(c *gin.Context) {
			user := c.Params.ByName("name")
			value, ok := DB[user]
			if ok {
				c.JSON(200, gin.H{"user": user, "value": value})
			} else {
				c.JSON(200, gin.H{"user": user, "status": "no value"})
			}
		})

		// Authorized group (uses gin.BasicAuth() middleware)
		// Same than:
		// authorized := r.Group("/")
		// authorized.Use(gin.BasicAuth(gin.Credentials{
		//	  "foo":  "bar",
		//	  "manu": "123",
		//}))
		authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
			"foo":  "bar", // user:foo password:bar
			"manu": "123", // user:manu password:123
		}))

		authorized.POST("admin", func(c *gin.Context) {
			user := c.MustGet(gin.AuthUserKey).(string)

			// Parse JSON
			var json struct {
				Value string `json:"value" binding:"required"`
			}

			if c.Bind(&json) == nil {
				DB[user] = json.Value
				c.JSON(200, gin.H{"status": "ok"})
			}
		})
	*/

	// Listen and Server in 0.0.0.0:8080
	router.Run(":9080")
}
