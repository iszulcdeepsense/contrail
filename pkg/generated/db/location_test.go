package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/models"
)

func TestLocation(t *testing.T) {
	t.Parallel()
	db := testDB
	common.UseTable(db, "location")
	defer func() {
		common.ClearTable(db, "location")
		if p := recover(); p != nil {
			panic(p)
		}
	}()
	model := models.MakeLocation()
	model.UUID = "dummy_uuid"

	err := common.DoInTransaction(db, func(tx *sql.Tx) error {
		return CreateLocation(tx, model)
	})
	if err != nil {
		t.Fatal("create failed", err)
	}

	err = common.DoInTransaction(db, func(tx *sql.Tx) error {
		models, err := ListLocation(tx, &common.ListSpec{Limit: 1})
		if err != nil {
			return err
		}
		if len(models) != 1 {
			return fmt.Errorf("expected one element")
		}
		return nil
	})
	if err != nil {
		t.Fatal("list failed", err)
	}

	err = common.DoInTransaction(db, func(tx *sql.Tx) error {
		model, err := ShowLocation(tx, model.UUID)
		if err != nil {
			return err
		}
		if model == nil || model.UUID != "dummy_uuid" {
			return fmt.Errorf("show failed")
		}
		return nil
	})
	if err != nil {
		t.Fatal("show failed", err)
	}

	err = common.DoInTransaction(db, func(tx *sql.Tx) error {
		return DeleteLocation(tx, model.UUID)
	})
	if err != nil {
		t.Fatal("delete failed", err)
	}

	err = common.DoInTransaction(db, func(tx *sql.Tx) error {
		models, err := ListLocation(tx, &common.ListSpec{Limit: 1})
		if err != nil {
			return err
		}
		if len(models) != 0 {
			return fmt.Errorf("expected no element")
		}
		return nil
	})
	if err != nil {
		t.Fatal("list failed", err)
	}
	return
}