package api

import (
	"database/sql"
	"net/http"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/db"
	"github.com/Juniper/contrail/pkg/generated/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
)

//KubernetesClusterRESTAPI
type KubernetesClusterRESTAPI struct {
	DB *sql.DB
}

type KubernetesClusterCreateRequest struct {
	Data *models.KubernetesCluster `json:"kubernetes-cluster"`
}

//Path returns api path for collections.
func (api *KubernetesClusterRESTAPI) Path() string {
	return "/kubernetes-clusters"
}

//LongPath returns api path for elements.
func (api *KubernetesClusterRESTAPI) LongPath() string {
	return "/kubernetes-cluster/:id"
}

//SetDB sets db object
func (api *KubernetesClusterRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *KubernetesClusterRESTAPI) Create(c echo.Context) error {
	requestData := &KubernetesClusterCreateRequest{
		Data: models.MakeKubernetesCluster(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "kubernetes_cluster",
		}).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	model := requestData.Data
	if model == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	if model.UUID == "" {
		model.UUID = uuid.NewV4().String()
	}
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.CreateKubernetesCluster(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "kubernetes_cluster",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *KubernetesClusterRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *KubernetesClusterRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteKubernetesCluster(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *KubernetesClusterRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.KubernetesCluster
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowKubernetesCluster(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"kubernetes_cluster": result,
	})
}

//List handles a List REST API Request.
func (api *KubernetesClusterRESTAPI) List(c echo.Context) error {
	var result []*models.KubernetesCluster
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListKubernetesCluster(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"kubernetes-clusters": result,
	})
}