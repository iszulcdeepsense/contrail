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

//AnalyticsNodeRESTAPI
type AnalyticsNodeRESTAPI struct {
	DB *sql.DB
}

type AnalyticsNodeCreateRequest struct {
	Data *models.AnalyticsNode `json:"analytics-node"`
}

//Path returns api path for collections.
func (api *AnalyticsNodeRESTAPI) Path() string {
	return "/analytics-nodes"
}

//LongPath returns api path for elements.
func (api *AnalyticsNodeRESTAPI) LongPath() string {
	return "/analytics-node/:id"
}

//SetDB sets db object
func (api *AnalyticsNodeRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *AnalyticsNodeRESTAPI) Create(c echo.Context) error {
	requestData := &AnalyticsNodeCreateRequest{
		Data: models.MakeAnalyticsNode(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "analytics_node",
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
			return db.CreateAnalyticsNode(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "analytics_node",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *AnalyticsNodeRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *AnalyticsNodeRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteAnalyticsNode(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *AnalyticsNodeRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.AnalyticsNode
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowAnalyticsNode(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"analytics_node": result,
	})
}

//List handles a List REST API Request.
func (api *AnalyticsNodeRESTAPI) List(c echo.Context) error {
	var result []*models.AnalyticsNode
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListAnalyticsNode(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"analytics-nodes": result,
	})
}