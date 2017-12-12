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

//ApplicationPolicySetRESTAPI
type ApplicationPolicySetRESTAPI struct {
	DB *sql.DB
}

type ApplicationPolicySetCreateRequest struct {
	Data *models.ApplicationPolicySet `json:"application-policy-set"`
}

//Path returns api path for collections.
func (api *ApplicationPolicySetRESTAPI) Path() string {
	return "/application-policy-sets"
}

//LongPath returns api path for elements.
func (api *ApplicationPolicySetRESTAPI) LongPath() string {
	return "/application-policy-set/:id"
}

//SetDB sets db object
func (api *ApplicationPolicySetRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *ApplicationPolicySetRESTAPI) Create(c echo.Context) error {
	requestData := &ApplicationPolicySetCreateRequest{
		Data: models.MakeApplicationPolicySet(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "application_policy_set",
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
			return db.CreateApplicationPolicySet(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "application_policy_set",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *ApplicationPolicySetRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *ApplicationPolicySetRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeleteApplicationPolicySet(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *ApplicationPolicySetRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.ApplicationPolicySet
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowApplicationPolicySet(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"application_policy_set": result,
	})
}

//List handles a List REST API Request.
func (api *ApplicationPolicySetRESTAPI) List(c echo.Context) error {
	var result []*models.ApplicationPolicySet
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListApplicationPolicySet(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"application-policy-sets": result,
	})
}