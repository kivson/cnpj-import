package downloader

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/assert"
)

func TestListUrls(t *testing.T) {
	html := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<html>
 <head>
  <title>Index of /CNPJ</title>
 </head>
 <body>
<h1>Index of /CNPJ</h1>
  <table>
   <tr><th valign="top"><img src="/icons/blank.gif" alt="[ICO]"></th><th><a href="?C=N;O=D">Name</a></th><th><a href="?C=M;O=A">Last modified</a></th><th><a href="?C=S;O=A">Size</a></th><th><a href="?C=D;O=A">Description</a></th></tr>
   <tr><th colspan="5"><hr></th></tr>
<tr><td valign="top"><img src="/icons/back.gif" alt="[PARENTDIR]"></td><td><a href="/">Parent Directory</a>       </td><td>&nbsp;</td><td align="right">  - </td><td>&nbsp;</td></tr>
<tr><td valign="top"><img src="/icons/compressed.gif" alt="[   ]"></td><td><a href="Cnaes.zip">Cnaes.zip</a>              </td><td align="right">2024-07-19 07:41  </td><td align="right"> 22K</td><td>&nbsp;</td></tr>
<tr><td valign="top"><img src="/icons/compressed.gif" alt="[   ]"></td><td><a href="Empresas0.zip">Empresas0.zip</a>          </td><td align="right">2024-07-19 08:37  </td><td align="right">334M</td><td>&nbsp;</td></tr>
<tr><td valign="top"><img src="/icons/compressed.gif" alt="[   ]"></td><td><a href="Empresas1.zip">Empresas1.zip</a>          </td><td align="right">2024-07-19 08:47  </td><td align="right"> 74M</td><td>&nbsp;</td></tr>`
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(html))
	}))
	defer func() { testServer.Close() }()
	downloader := Downloader{
		HttpClient: &http.Client{},
		Graber:     grab.NewClient(),
		BaseUrl:    testServer.URL,
	}

	urls, _ := downloader.retriveUrls()
	assert.Len(t, urls, 3)
}

func TestDownloadFile(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		file, _ := os.Open("../tests/Cnaes.zip")
		io.Copy(res, file)
	}))
	defer func() { testServer.Close() }()
	destination := "../tests/test_download/"
	fName := "/teste.zip"

	os.Remove(path.Join(destination, fName))

	downloader := Downloader{
		HttpClient: &http.Client{},
		Graber:     grab.NewClient(),
		BaseUrl:    testServer.URL,
		Dest:       destination,
	}
	bar := progressbar.DefaultBytes(int64(10000), "Tests downlolad")
	downloader.downloadFile(testServer.URL+fName, bar.Add)

	_, err := os.Stat(path.Join(destination, fName))
	assert.Nil(t, err)
	os.Remove(path.Join(destination, fName))
}

func TestDownloadRealFile(t *testing.T) {
	destination := "../tests/test_download/"
	fName := "/Cnaes.zip"

	os.Remove(path.Join(destination, fName))

	downloader := Downloader{
		HttpClient: &http.Client{},
		Graber:     grab.NewClient(),
		Dest:       "../tests/test_download/",
	}
	bar := progressbar.DefaultBytes(int64(10000), "Tests downlolad")
	downloader.downloadFile("https://dadosabertos.rfb.gov.br/CNPJ/Cnaes.zip", bar.Add)

	_, err := os.Stat(path.Join(destination, fName))
	assert.Nil(t, err)
	os.Remove(path.Join(destination, fName))
}
