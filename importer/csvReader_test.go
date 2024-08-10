package importer

import (
	"context"
	"testing"

	"github.com/kivson/cnpj-import/model"
	"github.com/stretchr/testify/assert"
)

var expectedEmpresa = model.Empresa{
	RazaoSocial: "OZINETE DELFINO CALDAS",
	CnpjBase:    "41273594", CodNaturezaJuridica: ptrOf(2135), CodQualificacaoResponsavel: ptrOf(50), Porte: 01, EnteResponsavel: "", CapitalSocial: 500000,
}

func ptrOf(n int) *int {
	return &n
}

var (
	expectedMotivo       = model.MotivoSituacaoCadastral{Descricao: "SEM MOTIVO", Codigo: ptrOf(0)}
	expectedCNAE         = model.CNAE{Descricao: "Cultivo de arroz", CNAE: ptrOf(111301)}
	expectedMunicipio    = model.Municipio{Nome: "GUAJARA-MIRIM", Codigo: ptrOf(1)}
	expectedNatureza     = model.NaturezaJuridica{Descricao: "Natureza Jurídica não informada", Codigo: ptrOf(0)}
	expectedPais         = model.Pais{Descricao: "COLIS POSTAUX", Codigo: ptrOf(0)}
	expectedQualificacao = model.QualificacaoSocio{Descricao: "Não informada", Codigo: ptrOf(0)}
	expectedSimples      = model.OptanteSimples{CnpjBase: "00000000", Simples: "N", SimplesInicio: "20070701", SimplesFim: "20070701", Simei: "N", SimeiInicio: "20090701", SimeiFim: "20090701"}
	expectedSocio        = model.Socio{
		CnpjBase:                       "07396865",
		TipoPessoa:                     2,
		Nome:                           "GERSON HOFFMANN",
		CPFCNPJ:                        "***240659**",
		QualificacaoCod:                ptrOf(49),
		Data:                           "20050518",
		CPFFormatadoRepresentanteLegal: "***000000**",
		NomeRepresentanteLegal:         "",
		QualificacaoRepresentanteCod:   ptrOf(0),
		FaixaEtaria:                    "5",
		PaisCod:                        nil,
		Pais:                           model.Pais{Descricao: "", Codigo: nil},
	}
	expectedEstabelecimento = model.Estabelecimento{
		CnpjBase:                   "07396865",
		CnpjOrdem:                  "0001",
		CnpjDv:                     "68",
		MatrizFilial:               1,
		Fantasia:                   "",
		SituacaoCadastral:          "",
		DataSituacaoCadastral:      "20170210",
		MotivoSituacaoCadastralCod: ptrOf(1),
		CidadeExterior:             "",
		PaisCod:                    nil,
		DataAbertura:               "20050518",
		CNAEPrincipalCod:           ptrOf(1412602),
		CNAESecundariaCod:          ptrOf(1411801),
		EnderecoTipoLogradouro:     "RUA",
		EnderecoLogradouro:         "TUCANEIRA",
		EnderecoNumero:             "30",
		EnderecoComplemento:        "",
		EnderecoBairro:             "DOS LAGOS",
		EnderecoCEP:                "89136000",
		EnderecoUF:                 "SC",
		EnderecoMunicipioCod:       ptrOf(8297),
		Telefone1DDD:               "47",
		Telefone1Numero:            "33851125",
		Telefone2DDD:               "47",
		Telefone2Numero:            "33851125",
		FaxDDD:                     "47",
		FaxNumero:                  "33851125",
		Email:                      "",
		SituacaoEspecial:           "",
		DataSituacaoEspecial:       "",
	}
)

func RunZipTest[C any](fileName string, expected C, t *testing.T) {
	out := make(chan Record[C], 4)

	go ReadZipCsv(context.Background(), fileName, out)

	record := <-out

	if record.Error != nil {
		t.Fatal(record.Error)
	}
	assert.Equal(t, &expected, record.Data)
}

func TestZipCsvParse(t *testing.T) {
	testName := "Running zip tests for Empresa"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Empresa.zip", expectedEmpresa, t)
	})
	testName = "Running zip tests for Estabelecimento"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Estabelecimentos.zip", expectedEstabelecimento, t)
	})
	testName = "Running zip tests for Simples"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Simples.zip", expectedSimples, t)
	})
	testName = "Running zip tests for Socio"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Socios.zip", expectedSocio, t)
	})
	testName = "Running zip tests for CNAE"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Cnaes.zip", expectedCNAE, t)
	})

	testName = "Running zip tests for Motivos"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Motivos.zip", expectedMotivo, t)
	})

	testName = "Running zip tests for Municipios"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Municipios.zip", expectedMunicipio, t)
	})

	testName = "Running zip tests for Paises"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Paises.zip", expectedPais, t)
	})

	testName = "Running zip tests for Naturezas"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Naturezas.zip", expectedNatureza, t)
	})

	testName = "Running zip tests for Qualificacoes"
	t.Run(testName, func(t *testing.T) {
		RunZipTest("../tests/Qualificacoes.zip", expectedQualificacao, t)
	})
}

func BenchmarkCsvParseEmpresa(b *testing.B) {
	for n := 0; n < b.N; n++ {
		out := make(chan Record[model.Empresa], 4)
		go func() {
			ReadZipCsv(context.Background(), "../tests/Empresa.zip", out)
			close(out)
		}()

		for range out {
		}
	}
}
