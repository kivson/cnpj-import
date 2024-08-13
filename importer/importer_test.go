package importer

import (
	"os"
	"testing"

	"github.com/kivson/cnpj-import/model"
	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestReadFilesError(t *testing.T) {
	folder := "INVALID"
	_, err := loadFiles(folder)
	assert.NotNil(t, err)
}

func TestReadFiles(t *testing.T) {
	folder := "../tests"
	files, err := loadFiles(folder)
	assert.Nil(t, err)
	assert.Contains(t, files, "../tests/Empresa.zip")
}

func TestImportFolder(t *testing.T) {
	folder := "../tests"
	dsn := "./test.db"

	importer := Importer{DbDsn: dsn, DbType: "sqlite"}
	importer.ImportZipFolder(folder)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	assert.Nil(t, err)

	var count int64
	db.Model(&model.Empresa{}).Count(&count)
	assert.Equal(t, int64(10), count)

	os.Remove("./test.db")
}

func BenchmarkInsertion(b *testing.B) {
	db, err := gorm.Open(sqlite.Open("bench.db?_journal_mode=MEMORY&_sync=OFF"), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("unable to connect to database")
	}
	service := NewStorageService(db)
	bar := progressbar.Default(
		-1,
		"Processing Records",
	)
	for i := 0; i < b.N; i++ {
		readAndSave[model.Empresa]("../zips/Empresas1.zip", service, bar.Add)
	}
}
