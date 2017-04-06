package model

import "time"

// Product represents a product.
type Product struct {
	ID          string    `json:"id,omitempty"`
	Description string    `json:"descriptions,omitempty"`
	Versions    []Version `json:"versions,omitempty"`
}

// Version represents a product version.
type Version struct {
	ID          string       `json:"id,omitempty"`
	Description string       `json:"description,omitempty"`
	Deployments []Deployment `json:"deployments,omitempty"`
}

// Status represents a Deployment request status.
type Status int8

const (
	// PENDING deployment request.
	PENDING Status = iota
	// GRANTED deployment request, waiting for actuation.
	GRANTED
	// PERFORMED deployment request.
	PERFORMED
)

// Deployment represents a potential deployment; it has an associated
// authorisation request, encoded in the Status field, which can be pending,
// or granted or performed (this last state might not be used).
// Once the authorisation is granted, this resource assumes a secondary
// representation in the form of AuthorisedDeployment, which can be used
// to trigger a CI/CD pipeline into performing the necessary steps to deploy
// the associated product/version onto the given environment.
type Deployment struct {
	Order       int       `json:"order,omitempty"`
	Environment string    `json:"environment,omitempty"`
	Status      Status    `json:"status,omitempty"`
	GrantedBy   string    `json:"grantedBy,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
}

func GetAllProducts() ([]Product, error) {
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
