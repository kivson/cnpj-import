package importer

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kivson/cnpj-import/model"
	"golang.org/x/text/encoding/charmap"
)

type Record[C any] struct {
	Data  *C
	Error error
}

func ReadZipCsv[C any](ctx context.Context, filePath string, out chan<- Record[C]) error {
	zipFile, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

loop:
	for _, file := range zipFile.File {
		csvFile, err := file.Open()
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			break loop
		default:
			ReadCsv(ctx, csvFile, out)
		}
	}

	return nil
}

func ReadCsv[C any](ctx context.Context, r io.Reader, out chan<- Record[C]) {
	decoder := charmap.ISO8859_1.NewDecoder()
	reader := csv.NewReader(decoder.Reader(r))
	reader.Comma = ';'
loop:
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		var data C

		err = Unmarshal(record, &data)
		select {
		case <-ctx.Done():
			break loop
		case out <- Record[C]{&data, err}:
		}
	}
}

func Unmarshal(csvRecord []string, data any) error {
	switch val := data.(type) {
	case *model.Estabelecimento:
		val.CnpjBase = csvRecord[0]
		val.CnpjOrdem = csvRecord[1]
		val.CnpjDv = csvRecord[2]
		val.MatrizFilial, _ = strconv.Atoi(csvRecord[3])
		val.Fantasia = csvRecord[4]
		val.SituacaoCadastral = csvRecord[5]
		val.DataSituacaoCadastral = csvRecord[6]
		val.MotivoSituacaoCadastralCod = intOrNil(csvRecord[7])
		val.CidadeExterior = csvRecord[8]
		val.PaisCod = intOrNil(csvRecord[9])
		val.DataAbertura = csvRecord[10]
		val.CNAEPrincipalCod = intOrNil(csvRecord[11])
		val.CNAESecundariaCod = intOrNil(csvRecord[12])
		val.EnderecoTipoLogradouro = csvRecord[13]
		val.EnderecoLogradouro = csvRecord[14]
		val.EnderecoNumero = csvRecord[15]
		val.EnderecoComplemento = csvRecord[16]
		val.EnderecoBairro = csvRecord[17]
		val.EnderecoCEP = csvRecord[18]
		val.EnderecoUF = csvRecord[19]
		val.EnderecoMunicipioCod = intOrNil(csvRecord[20])
		val.Telefone1DDD = csvRecord[21]
		val.Telefone1Numero = csvRecord[22]
		val.Telefone2DDD = csvRecord[23]
		val.Telefone2Numero = csvRecord[24]
		val.FaxDDD = csvRecord[25]
		val.FaxNumero = csvRecord[26]
		val.Email = csvRecord[27]
		val.SituacaoCadastral = csvRecord[28]
		val.DataSituacaoEspecial = csvRecord[29]
	case *model.Socio:
		val.CnpjBase = csvRecord[0]
		val.TipoPessoa, _ = strconv.Atoi(csvRecord[1])
		val.Nome = csvRecord[2]
		val.CPFCNPJ = csvRecord[3]
		val.QualificacaoCod = intOrNil(csvRecord[4])
		val.Data = csvRecord[5]
		val.PaisCod = intOrNil(csvRecord[6])
		val.CPFFormatadoRepresentanteLegal = csvRecord[7]
		val.NomeRepresentanteLegal = csvRecord[8]
		val.QualificacaoRepresentanteCod = intOrNil(csvRecord[9])
		val.FaixaEtaria = csvRecord[10]
	case *model.OptanteSimples:
		val.CnpjBase = csvRecord[0]
		val.Simples = csvRecord[1]
		val.SimplesInicio = csvRecord[2]
		val.SimplesFim = csvRecord[3]
		val.Simei = csvRecord[4]
		val.SimeiInicio = csvRecord[5]
		val.SimeiFim = csvRecord[6]
	case *model.Empresa:
		val.CnpjBase = csvRecord[0]
		val.RazaoSocial = csvRecord[1]
		val.CodNaturezaJuridica = intOrNil(csvRecord[2])
		val.CodQualificacaoResponsavel = intOrNil(csvRecord[3])
		val.Porte, _ = strconv.Atoi(csvRecord[5])
		val.EnteResponsavel = csvRecord[6]
		val.CapitalSocial = moneyToUint(csvRecord[4])
	case *model.CNAE:
		val.CNAE = intOrNil(csvRecord[0])
		val.Descricao = csvRecord[1]
	case *model.Pais:
		val.Codigo = intOrNil(csvRecord[0])
		val.Descricao = csvRecord[1]
	case *model.MotivoSituacaoCadastral:
		val.Codigo = intOrNil(csvRecord[0])
		val.Descricao = csvRecord[1]
	case *model.NaturezaJuridica:
		val.Codigo = intOrNil(csvRecord[0])
		val.Descricao = csvRecord[1]
	case *model.QualificacaoSocio:
		val.Codigo = intOrNil(csvRecord[0])
		val.Descricao = csvRecord[1]
	case *model.Municipio:
		val.Codigo = intOrNil(csvRecord[0])
		val.Nome = csvRecord[1]
	default:
		return fmt.Errorf("unable to Unmarshal %v", data)
	}
	return nil
}

func intOrNil(num string) *int {
	if num == "" {
		return nil
	}
	resp, err := strconv.Atoi(num)
	if err != nil {
		return nil
	}
	return &resp
}

func moneyToUint(num string) uint64 {
	var clean strings.Builder

	for idx := range num {
		if num[idx] >= '0' && num[idx] <= '9' {
			clean.WriteByte(num[idx])
		}
	}

	out, _ := strconv.Atoi(clean.String())
	return uint64(out)
}
