package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func parse_args(ctx *gin.Context, need_api_key bool) (path string, file_key string, api_key string, ok bool) {
	ok = true
	path = ctx.Params.ByName("path")
	file_key = ctx.Query("file_key")
	if file_key == "" {
		ok = false
	}
	api_key = ctx.Query("api_key")
	if need_api_key && api_key == "" {
		ok = false
	}
	return
}

func upload(ctx *gin.Context) {
	path, file_key, api_key, ok := parse_args(ctx, true)
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	vaild, err := IsVaildUser(api_key)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if !vaild {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	content, err := ctx.GetRawData()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	secret, err := Encrypt(content, file_key)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = DefaultStorageBackend.Write(Hash(path), secret)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func download(ctx *gin.Context) {
	path, file_key, _, ok := parse_args(ctx, false)
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	plain, err := DefaultStorageBackend.Read(Hash(path))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	content, err := Decrypt(plain, file_key)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Data(http.StatusOK, "text/plain", content)
}

func Serve() error {
	r := gin.Default()
	r.PUT("/*path", upload)
	r.GET("/*path", download)
	return r.Run()
}
