package model

import (
	"strconv"
	"strings"
)

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

type Empresa struct {
	RazaoSocial     string
	CnpjBase        string `gorm:"primaryKey"`
	EnteResponsavel string

	CodNaturezaJuridica *int
	NaturezaJuridica    NaturezaJuridica `gorm:"foreignKey:CodNaturezaJuridica"`

	CodQualificacaoResponsavel *int
	QualificacaoSocio          QualificacaoSocio `gorm:"foreignKey:CodQualificacaoResponsavel"`

	Porte         int
	CapitalSocial uint64
}

type Socio struct {
	CnpjBase                       string
	Empresa                        Empresa `gorm:"foreignKey:CnpjBase"`
	TipoPessoa                     int
	Nome                           string
	CPFCNPJ                        string
	QualificacaoCod                *int
	Qualificacao                   QualificacaoSocio `gorm:"foreignKey:QualificacaoCod"`
	Data                           string
	CPFFormatadoRepresentanteLegal string
	NomeRepresentanteLegal         string
	QualificacaoRepresentanteCod   *int
	QualificacaoRepresentante      QualificacaoSocio `gorm:"foreignKey:QualificacaoRepresentanteCod"`
	FaixaEtaria                    string
	PaisCod                        *int
	Pais                           Pais `gorm:"foreignKey:PaisCod"`
}

type Estabelecimento struct {
	CnpjBase                   string
	Empresa                    Empresa `gorm:"foreignKey:CnpjBase"`
	CnpjOrdem                  string
	CnpjDv                     string
	MatrizFilial               int
	Fantasia                   string
	SituacaoCadastral          string
	DataSituacaoCadastral      string
	MotivoSituacaoCadastralCod *int
	MotivoSituacaoCadastral    MotivoSituacaoCadastral `gorm:"foreignKey:MotivoSituacaoCadastralCod"`
	CidadeExterior             string
	PaisCod                    *int
	Pais                       Pais `gorm:"foreignKey:PaisCod"`
	DataAbertura               string
	CNAEPrincipalCod           *int
	CNAEPrincipal              CNAE `gorm:"foreignKey:CNAEPrincipalCod"`
	CNAESecundariaCod          *int
	CNAESecundaria             CNAE `gorm:"foreignKey:CNAESecundariaCod"`
	EnderecoTipoLogradouro     string
	EnderecoLogradouro         string
	EnderecoNumero             string
	EnderecoComplemento        string
	EnderecoBairro             string
	EnderecoCEP                string
	EnderecoUF                 string
	EnderecoMunicipioCod       *int
	EnrerecoMunicipio          Municipio `gorm:"foreignKey:EnderecoMunicipioCod"`
	Telefone1DDD               string
	Telefone1Numero            string
	Telefone2DDD               string
	Telefone2Numero            string
	FaxDDD                     string
	FaxNumero                  string
	Email                      string
	SituacaoEspecial           string
	DataSituacaoEspecial       string
}

type OptanteSimples struct {
	CnpjBase      string
	Empresa       Empresa `gorm:"foreignKey:CnpjBase"`
	Simples       string
	SimplesInicio string
	SimplesFim    string
	Simei         string
	SimeiInicio   string
	SimeiFim      string
}

type CNAE struct {
	Descricao string
	CNAE      *int `gorm:"primaryKey"`
}

type Municipio struct {
	Nome   string
	Codigo *int `gorm:"primaryKey"`
}

type NaturezaJuridica struct {
	Descricao string
	Codigo    *int `gorm:"primaryKey"`
}

type QualificacaoSocio struct {
	Descricao string
	Codigo    *int ` gorm:"primaryKey"`
}

type Pais struct {
	Descricao string
	Codigo    *int `gorm:"primaryKey"`
}

type MotivoSituacaoCadastral struct {
	Descricao string
	Codigo    *int `gorm:"primaryKey"`
}
