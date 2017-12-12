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

//SecurityLoggingObjectRESTAPI
type SecurityLoggingObjectRESTAPI struct {
	DB *sql.DB
}

type SecurityLoggingObjectCreateRequest struct {
	Data *models.SecurityLoggingObject `json:"security-logging-object"`
}

//Path returns api path for collections.
func (api *SecurityLoggingObjectRESTAPI) Path() string {
	return "/security-logging-objects"
}

//LongPath returns api path for elements.
func (api *SecurityLoggingObjectRESTAPI) LongPath() string {
	return "/security-logging-object/:id"
}

//SetDB sets db object
func (api *SecurityLoggingObjectRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *SecurityLoggingObjectRESTAPI) Create(c echo.Context) error {
	requestData := &SecurityLoggingObjectCreateRequest{
		Data: models.MakeSecurityLoggingObject(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "security_logging_object",
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
			return db.CreateSecurityLoggingObject(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "security_logging_object",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *SecurityLoggingObjectRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *SecurityLoggingObjectRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteSecurityLoggingObject(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *SecurityLoggingObjectRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.SecurityLoggingObject
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowSecurityLoggingObject(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"security_logging_object": result,
	})
}

//List handles a List REST API Request.
func (api *SecurityLoggingObjectRESTAPI) List(c echo.Context) error {
	var result []*models.SecurityLoggingObject
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListSecurityLoggingObject(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"security-logging-objects": result,
	})
}