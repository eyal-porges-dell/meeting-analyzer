// Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
package main

import (
	"context"
	"errors"
	"meeting-analyzer/server/api/rest/generated"
	"meeting-analyzer/server/models/errorresponse"
	"meeting-analyzer/server/services/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meeting-analyzer/server/commons/constants"

	"meeting-analyzer/server/api/rest/controller"

	"meeting-analyzer/server/repositories/db"

	"eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/database"
	dbConstants "eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/database/constants"
	"eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/database/interfaces"
	"eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/database/models"
	log "eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/logger"
	"eos2git.cec.lab.emc.com/ISG-Edge/hzp-iam-lib-go/server/services/authcontextsvc"
	"github.com/gorilla/mux"
)

// main where the main application runs
func main() {
	log.Initialize(constants.ComponentName)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()
	//dbo, err := initializeDB(ctx, fetchPostgresConfig(ctx, constants.DatabaseName))
	//
	//if err != nil {
	//	log.Fatal(ctx, nil, "", err, "failed to establish connection with database")
	//}

	if err := runMain(ctx, c, nil); err != nil {
		log.Fatal(ctx, nil, "", err, "failed to start main process")
	}
}

// runMain where the server runs
func runMain(ctx context.Context, c chan os.Signal,
	dbo interfaces.Database) error {
	repo, err := db.NewRepository(dbo)
	if err != nil {
		log.Error(ctx, nil, "", err, "failed to init rules repo")
		return err
	}

	svc, err := service.NewSvc(ctx, repo)
	if err != nil {
		log.Error(ctx, nil, "", err, "failed to init service")
		return err
	}
	server := CreateServer(ctx, svc)

	go func() {
		if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			log.Fatal(ctx, nil, "", err, "unable to start server")
		}
	}()

	// Waiting for SIGINT (kill -2)
	<-c

	ctx, cancel := context.WithTimeout(ctx, constants.ContextTimeout*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Error(ctx, nil, "", err, "Unable to gracefully shutdown")
	}
	return err
}

// initializeDB where the db instance is created
func initializeDB(ctx context.Context, config *models.Config) (interfaces.Database, error) {
	dbObj, err := database.NewDatabaseInstance(ctx, config, dbConstants.Postgres, constants.MigrationFolderPath)

	if err != nil {
		return nil, err
	}

	return dbObj, nil
}

// fetchPostgresConfig fetches the postgres config values
func fetchPostgresConfig(ctx context.Context, dbName string) *models.Config {
	varNotFoundErr := errors.New("var not found")
	host, found := os.LookupEnv(constants.EnvVarDBHost)
	if !found {
		log.Error(ctx, nil, "", varNotFoundErr, constants.EnvVarDBHost)
	}
	port, found := os.LookupEnv(constants.EnvVarDBPort)
	if !found {
		log.Error(ctx, nil, "", varNotFoundErr, constants.EnvVarDBPort)
	}
	user, found := os.LookupEnv(constants.EnvVarDBUser)
	if !found {
		log.Error(ctx, nil, "", varNotFoundErr, constants.EnvVarDBUser)
	}
	pa, found := os.LookupEnv(constants.EnvVarDBPA)
	if !found {
		log.Error(ctx, nil, "", varNotFoundErr, constants.EnvVarDBPA)
	}
	return &models.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Pass:     pa,
		Database: dbName,
	}
}

// CreateServer adds the health livenesss and health dependencies
func CreateServer(ctx context.Context, svc service.Service) *HTTPServer {
	router := mux.NewRouter()
	router.Use(log.AddCorrelationIDMiddleware)
	router.Use(authcontextsvc.EnrichContextWithEstateInitiatorCtxMiddleware)
	RegisterOpenAPIHandler(ctx, router, svc)
	return NewHTTPServer(":"+GetPort(), router)
}

func RegisterOpenAPIHandler(ctx context.Context, router *mux.Router, svc service.Service) {
	swagger, err := generated.GetSwagger()
	if err != nil {
		log.Fatal(ctx, nil, "", err, "error loading swagger spec")
	}
	// This lines removes server name validation
	swagger.Servers = nil

	handler := generated.NewStrictHandlerWithOptions(controller.NewController(svc), nil, errorresponse.StrictHTTPServerOptions)
	generated.HandlerFromMux(handler, router)
}

// GetPort returns the port values
func GetPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = constants.DefaultPort
	}
	return port
}

type HTTPServer struct {
	*http.Server
}

// NewHTTPServer creates a new server
func NewHTTPServer(adr string, router *mux.Router) *HTTPServer {
	return &HTTPServer{
		&http.Server{
			Addr:              adr,
			Handler:           router,
			ReadHeaderTimeout: time.Second * constants.HTTPServerRequestReadTimeOut,
		},
	}
}
