package importer

import (
	"github.com/kivson/cnpj-import/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBType string

const (
	SQLite DBType = "sqlite"
)

type StorageService struct {
	db *gorm.DB
}

func NewStorageService(db *gorm.DB) *StorageService {
	service := StorageService{db}

	service.db.AutoMigrate(
		&model.Empresa{},
		&model.CNAE{}, &model.Pais{}, &model.Socio{},
		&model.Municipio{}, &model.OptanteSimples{},
		&model.Estabelecimento{}, &model.NaturezaJuridica{},
		&model.QualificacaoSocio{}, &model.MotivoSituacaoCadastral{},
	)

	return &service
}

func (store StorageService) Store(model any) error {
	result := store.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(model)

	return result.Error
}
