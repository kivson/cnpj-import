package importer

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kivson/cnpj-import/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProgressUpdFun func(records int) error

type Importer struct {
	ProgressFn func(record int) error
	DbType     DBType
	DbDsn      string
}

func (imp Importer) ImportZipFolder(folder string) {
	db, err := gorm.Open(sqlite.Open(imp.DbDsn), &gorm.Config{})
	if err != nil {
		panic("unable to connect to database")
	}
	files, err := loadFiles(folder)
	if err != nil {
		panic(err)
	}

	service := NewStorageService(db)

	for _, fileName := range files {
		lowerName := strings.ToLower(fileName)
		switch {
		case strings.Contains(lowerName, "empresa"):
			readAndSave[model.Empresa](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "estabelecimento"):
			readAndSave[model.Estabelecimento](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "socio"):
			readAndSave[model.Socio](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "qualificacoes"):
			readAndSave[model.QualificacaoSocio](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "paises"):
			readAndSave[model.Pais](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "naturezas"):
			readAndSave[model.NaturezaJuridica](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "municipios"):
			readAndSave[model.Municipio](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "motivos"):
			readAndSave[model.MotivoSituacaoCadastral](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "cnae"):
			readAndSave[model.CNAE](fileName, service, imp.ProgressFn)
		case strings.Contains(lowerName, "simples"):
			readAndSave[model.OptanteSimples](fileName, service, imp.ProgressFn)
		default:
			fmt.Printf("Unable to identify type of file %s ", fileName)
		}
	}
}

func loadFiles(folder string) ([]string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))

	for i := range files {
		if strings.HasSuffix(strings.ToLower(files[i].Name()), "zip") {
			paths = append(paths, path.Join(folder, files[i].Name()))
		}
	}
	return paths, nil
}

func readAndSave[A any](filePath string, service *StorageService, updFn ProgressUpdFun) {
	output := make(chan Record[A], 10000)

	go func() {
		ReadZipCsv(context.Background(), filePath, output)
		close(output)
	}()

	var count int
	batchSize := 1000
	buffer := make([]*A, 0, batchSize)

	for record := range output {
		count++
		buffer = append(buffer, record.Data)

		if count%batchSize == 0 {
			service.Store(buffer)
			if updFn != nil {
				updFn(batchSize)
			}
			buffer = buffer[:0]
		}

	}
	if len(buffer) > 0 {
		service.Store(buffer)
	}
}
