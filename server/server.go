package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CYTMWIA/secret-server/backend"
	"github.com/CYTMWIA/secret-server/crypto"
	"github.com/gin-gonic/gin"
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

type WebApp struct {
	mode   string
	server *http.Server

	backend backend.StorageBackend

	user_list []string
}

func (app *WebApp) is_vaild_user(api_key string) (bool, error) {
	hkey := crypto.Hash(api_key)
	for _, key := range app.user_list {
		if key == hkey {
			return true, nil
		}
	}

	return false, nil
}

func (app *WebApp) shutdown_when_function_mode(_ *gin.Context) {
	if app.mode == "function" {
		go func() {
			if err := app.server.Shutdown(context.Background()); err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (app *WebApp) upload(ctx *gin.Context) {
	path, file_key, api_key, ok := parse_args(ctx, true)
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	vaild, err := app.is_vaild_user(api_key)
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

	secret, err := crypto.Encrypt(content, file_key, path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = app.backend.Write(crypto.Hash(path), secret)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (app *WebApp) download(ctx *gin.Context) {
	path, file_key, _, ok := parse_args(ctx, false)
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	plain, err := app.backend.Read(crypto.Hash(path))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	content, err := crypto.Decrypt(plain, file_key, path)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Data(http.StatusOK, "text/plain", content)
}

func Serve(addr string, mode string, backend backend.StorageBackend, user_list []string) error {
	app := WebApp{mode: mode, backend: backend, user_list: user_list}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.PUT("/*path", app.upload, app.shutdown_when_function_mode)
	router.GET("/*path", app.download, app.shutdown_when_function_mode)

	app.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	err := app.server.ListenAndServe()
	if err != nil && err == http.ErrServerClosed {
		return nil
	}
	return err
}
