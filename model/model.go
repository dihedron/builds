package model

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // don't refer to SQLITE3
	"github.com/pkg/errors"
)

// Product represents a product.
type Product struct {
	ID          uint      `gorm:"primary_key;unique_index:products_pk" json:"id"`
	Code        string    `gorm:"size:63;unique_index:uix_pcode" json:"code,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Contact     string    `json:"contact,omitempty"`
	Repository  string    `json:"repository,omitempty"`
	WebSite     string    `json:"website,omitempty"`
	Versions    []Version `json:"versions,omitempty"`
	CreatedAt   time.Time `json:"created,omitempty"`
	UpdatedAt   time.Time `json:"updated,omitempty"`
}

// Version represents a product version.
type Version struct {
	ID          uint      `gorm:"primary_key;unique_index:versions_pk"  json:"id"`
	ProductID   uint      `gorm:"unique_index:uix_pv"  json:"pid"`
	Code        string    `gorm:"unique_index:uix_pv" json:"code,omitempty"`
	Description string    `gorm:"type:varchar(1024)" json:"description,omitempty"`
	Repository  string    `json:"repository,omitempty"`
	Branch      string    `json:"branch,omitempty"`
	CreatedAt   time.Time `json:"created,omitempty"`
	UpdatedAt   time.Time `json:"updated,omitempty"`
}

// String formats a Product as a JSON-encoded string.
func (p Product) String() string {
	bytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return ""
	}
	return string(bytes[:])
}

// String formats a Version as a JSON-encoded string.
func (v Version) String() string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(bytes[:])
}

var db *gorm.DB

// New loads an existing SQLITE3 database from the given path, or creates
// one if not existing, and upates the tables definitions according to
// the current object model.
func New(dbpath string) error {

	var err error
	if db, err = gorm.Open("sqlite3", dbpath); err != nil {
		return errors.Wrap(err, "failed to load database driver")
	}

	if err = db.DB().Ping(); err != nil {
		return errors.Wrap(err, "failed to connect to database manager")
	}

	// instantiate or update the schema (does not drop anything)
	db.AutoMigrate(&Product{}, &Version{})

	return nil
}

// Close closes the current SQLITE3 database.
func Close() error {
	if err := db.Close(); err != nil {
		return errors.Wrap(err, "error closing the database")
	}
	return nil
}

// GetProducts returns the full list of products.
func GetProducts() []Product {
	var products []Product
	db.Find(&products)
	return products
}

// CreateProduct creates a new Product; if it contains Version references,
// those are created too.
func CreateProduct(product *Product) {
	db.Create(product)
}

// ReadProduct reads an existing product from the database; the provided
// object should contain the search criteria.
func ReadProduct(product *Product) {
	db.Find(product)
}

// UpdateProduct updates an existing product; if it contains Versions,
// those are updated as well.
func UpdateProduct(product *Product) {
	db.Update(product)
}

// DeleteProduct deletes an existing product from the datavbase; any existing
// linked Version objects are deleted as well (cascade).
func DeleteProduct(product *Product) {
	db.Delete(product)
}

/*

func init() {

	model.New("test.sqlite")

	product := Product{
		Code:        "gaia",
		Name:        "G.A.I.A. - Servizi per il Personale",
		Description: "GAIA è il portale web dei servizi aziendali non altrimenti disponibili su piattaforma SAP.",
		Contact:     "fabio.angeli@bancaditalia.it",
		Repository:  "https://gitlab.utenze.bankit.it/gaia",
		WebSite:     "http://infogaia/",
		Versions: []Version{
			{
				Code:        "1.0.0",
				Description: "First major release, 1.0 series",
				Repository:  "https://gitlab.utenze.bankit.it/gaia",
				Branch:      "ver_1_0_0",
			},
			{
				Code:        "1.0.1",
				Description: "First bugfix release of the 1.0 series",
				Branch:      "ver_1_0_1",
			},
		},
	}
	db.Create(&product)

	product = Product{
		Code:        "siparium",
		Name:        "SIPARIUM - Sistema Integrato Processi Aziendali per le Risorse UMane",
		Description: "GAIA è il portale web dei servizi aziendali per le risorse umane su piattaforma SAP.",
		Contact:     "roberto iapichino@bancaditalia.it",
		Repository:  "https://gitlab.utenze.bankit.it/siparium",
		WebSite:     "http://portale-sap/",
		Versions: []Version{
			{
				Code:        "1.0.0",
				Description: "First major release, 1.0 series",
				Repository:  "https://gitlab.utenze.bankit.it/siparium",
				Branch:      "ver_1_0_0",
			},
			{
				Code:        "1.0.1",
				Description: "First bugfix release of the 1.0 series",
				Branch:      "ver_1_0_1",
			},
		},
	}
	db.Create(&product)
}



func GetProducts() ([]Product, error) {
	if len(model) == 0 {
		return nil, ErrorNotFound
	}
	products := make([]Product, 0, len(model))
	for _, product := range model {
		products = append(products, product)
	}
	return products, nil
}

func GetProductByID(productId string) (Product, error) {
	product, ok := model[productId]
	if !ok {
		return Product{}, ErrorNotFound
	}
	return product, nil
}

func GetVersionIDsByProduct(product Product) ([]string, error) {
	return GetVersionIDsByProductID(product.ID)
}

func GetVersionIDsByProductID(productId string) ([]string, error) {
	product, ok := model[productId]
	if !ok {
		return nil, ErrorNotFound
	}
	versions := make([]string, 0, len(product.Versions))
	for _, version := range product.Versions {
		versions = append(versions, version.ID)
	}
	return versions, nil
}

func GetVersionByID(productId string, versionID string) (Version, error) {
	product, err := GetProductByID(productId)
	if err != nil {
		return Version{}, err
	}
	for _, version := range product.Versions {
		if version.ID == versionID {
			return version, nil
		}
	}
	return Version{}, ErrorNotFound
}

func UpdateDeployment(productId, versionId string, deployment Deployment) error {
	// do nothing
	return nil
}

// model is the current (in memory) data store.
var model map[string]Product

func init() {

	model = make(map[string]Product)

	model["gaia"] = Product{
		ID:          "gaia",
		Description: "Portale dei servizi aziendali per il personale",
		Versions: []Version{
			{
				ID:          "1.0.0",
				Description: "for a changelog see http://prodinfo/gaia/1.0.0",
				Deployments: []Deployment{
					{
						Order:       0,
						Environment: "Integration",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       1,
						Environment: "Quality",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       2,
						Environment: "Certification",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       3,
						Environment: "Production",
						Status:      PENDING,
					},
				},
			},
			{
				ID:          "1.0.1",
				Description: "for a changelog see http://prodinfo/gaia/1.0.1",
				Deployments: []Deployment{
					{
						Order:       0,
						Environment: "Integration",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       1,
						Environment: "Quality",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       2,
						Environment: "Certification",
						Status:      PENDING,
					},
					{
						Order:       3,
						Environment: "Production",
						Status:      PENDING,
					},
				},
			},
			{
				ID:          "1.0.2",
				Description: "for a changelog see http://prodinfo/gaia/1.0.2",
				Deployments: []Deployment{
					{
						Order:       0,
						Environment: "Integration",
						Status:      PENDING,
					},
					{
						Order:       1,
						Environment: "Quality",
						Status:      PENDING,
					},
					{
						Order:       2,
						Environment: "Certification",
						Status:      PENDING,
					},
					{
						Order:       3,
						Environment: "Production",
						Status:      PENDING,
					},
				},
			},
		},
	}

	model["siparium"] = Product{
		ID:          "siparium",
		Description: "Enterprise Resource Planning (ERP) per le Risorse Umane",
		Versions: []Version{
			{
				ID:          "2.0.0",
				Description: "for a changelog see http://prodinfo/siparium/2.0.0",
				Deployments: []Deployment{
					{
						Order:       0,
						Environment: "Integration",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       1,
						Environment: "Quality",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       2,
						Environment: "Certification",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       3,
						Environment: "Production",
						Status:      PENDING,
					},
				},
			},
			{
				ID:          "2.0.1",
				Description: "for a changelog see http://prodinfo/siparium/2.0.1",
				Deployments: []Deployment{
					{
						Order:       0,
						Environment: "Integration",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       1,
						Environment: "Quality",
						Status:      PERFORMED,
						Timestamp:   time.Now(),
						GrantedBy:   "d093154",
					},
					{
						Order:       2,
						Environment: "Certification",
						Status:      PENDING,
					},
					{
						Order:       3,
						Environment: "Production",
						Status:      PENDING,
					},
				},
			},
			{
				ID:          "1.0.2",
				Description: "for a changelog see http://prodinfo/gaia/1.0.2",
				Deployments: []Deployment{
					{
						Order:       0,
						Environment: "Integration",
						Status:      PENDING,
					},
					{
						Order:       1,
						Environment: "Quality",
						Status:      PENDING,
					},
					{
						Order:       2,
						Environment: "Certification",
						Status:      PENDING,
					},
					{
						Order:       3,
						Environment: "Production",
						Status:      PENDING,
					},
				},
			},
		},
	}
}
*/
