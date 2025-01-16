package repository

import (
	"fmt"

	"certimate/internal/domain"
	"certimate/internal/utils/app"

	"github.com/go-acme/lego/v4/registration"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/sync/singleflight"
)

type AcmeAccountRepository struct{}

func NewAcmeAccountRepository() *AcmeAccountRepository {
	return &AcmeAccountRepository{}
}

var g singleflight.Group

func (r *AcmeAccountRepository) GetByCAAndEmail(ca, email string) (*domain.AcmeAccount, error) {
	resp, err, _ := g.Do(fmt.Sprintf("acme_account_%s_%s", ca, email), func() (interface{}, error) {
		resp, err := app.GetApp().Dao().FindFirstRecordByFilter("acme_accounts", "ca={:ca} && email={:email}", dbx.Params{"ca": ca, "email": email})
		if err != nil {
			return nil, err
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("acme account not found")
	}

	record, ok := resp.(*models.Record)
	if !ok {
		return nil, fmt.Errorf("acme account not found")
	}

	resource := &registration.Resource{}
	if err := record.UnmarshalJSONField("resource", resource); err != nil {
		return nil, err
	}

	return &domain.AcmeAccount{
		Id:       record.GetString("id"),
		Ca:       record.GetString("ca"),
		Email:    record.GetString("email"),
		Key:      record.GetString("key"),
		Resource: resource,
		Created:  record.GetTime("created"),
		Updated:  record.GetTime("updated"),
	}, nil
}

func (r *AcmeAccountRepository) Save(ca, email, key string, resource *registration.Resource) error {
	collection, err := app.GetApp().Dao().FindCollectionByNameOrId("acme_accounts")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	record.Set("ca", ca)
	record.Set("email", email)
	record.Set("key", key)
	record.Set("resource", resource)
	return app.GetApp().Dao().Save(record)
}
