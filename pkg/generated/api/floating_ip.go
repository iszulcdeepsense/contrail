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

//FloatingIPRESTAPI
type FloatingIPRESTAPI struct {
	DB *sql.DB
}

type FloatingIPCreateRequest struct {
	Data *models.FloatingIP `json:"floating-ip"`
}

//Path returns api path for collections.
func (api *FloatingIPRESTAPI) Path() string {
	return "/floating-ips"
}

//LongPath returns api path for elements.
func (api *FloatingIPRESTAPI) LongPath() string {
	return "/floating-ip/:id"
}

//SetDB sets db object
func (api *FloatingIPRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *FloatingIPRESTAPI) Create(c echo.Context) error {
	requestData := &FloatingIPCreateRequest{
		Data: models.MakeFloatingIP(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "floating_ip",
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
			return db.CreateFloatingIP(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "floating_ip",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *FloatingIPRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *FloatingIPRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteFloatingIP(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *FloatingIPRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.FloatingIP
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowFloatingIP(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"floating_ip": result,
	})
}

//List handles a List REST API Request.
func (api *FloatingIPRESTAPI) List(c echo.Context) error {
	var result []*models.FloatingIP
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListFloatingIP(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"floating-ips": result,
	})
}