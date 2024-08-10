package importer

import (
	"testing"

	"github.com/kivson/cnpj-import/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSaveEmpresa(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{})
	assert.Nil(t, err)
	storage := NewStorageService(db)

	err = storage.Store(expectedEmpresa)
	assert.Nil(t, err)

	var dbEmpresa model.Empresa
	storage.db.First(&dbEmpresa)

	assert.Equal(t, expectedEmpresa, dbEmpresa)

	storage.Store(model.Empresa{RazaoSocial: "NEW NAME", CnpjBase: expectedEmpresa.CnpjBase})

	var count int64
	storage.db.Model(model.Empresa{}).Count(&count)

	assert.Equal(t, int64(1), count)

	storage.db.First(&dbEmpresa)
	assert.Equal(t, "NEW NAME", dbEmpresa.RazaoSocial)
}
