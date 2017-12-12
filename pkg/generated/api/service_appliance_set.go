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

//ServiceApplianceSetRESTAPI
type ServiceApplianceSetRESTAPI struct {
	DB *sql.DB
}

type ServiceApplianceSetCreateRequest struct {
	Data *models.ServiceApplianceSet `json:"service-appliance-set"`
}

//Path returns api path for collections.
func (api *ServiceApplianceSetRESTAPI) Path() string {
	return "/service-appliance-sets"
}

//LongPath returns api path for elements.
func (api *ServiceApplianceSetRESTAPI) LongPath() string {
	return "/service-appliance-set/:id"
}

//SetDB sets db object
func (api *ServiceApplianceSetRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *ServiceApplianceSetRESTAPI) Create(c echo.Context) error {
	requestData := &ServiceApplianceSetCreateRequest{
		Data: models.MakeServiceApplianceSet(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "service_appliance_set",
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
			return db.CreateServiceApplianceSet(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "service_appliance_set",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *ServiceApplianceSetRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *ServiceApplianceSetRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteServiceApplianceSet(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *ServiceApplianceSetRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.ServiceApplianceSet
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowServiceApplianceSet(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"service_appliance_set": result,
	})
}

//List handles a List REST API Request.
func (api *ServiceApplianceSetRESTAPI) List(c echo.Context) error {
	var result []*models.ServiceApplianceSet
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListServiceApplianceSet(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"service-appliance-sets": result,
	})
}