package orchestrator

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/agent"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/middlewares"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/storage"
	geteditnodes "github.com/Gandalf-Rus/distributed-calc2.0/internal/use_cases/get_edit_nodes"
	loginuser "github.com/Gandalf-Rus/distributed-calc2.0/internal/use_cases/login_user"
	postexpression "github.com/Gandalf-Rus/distributed-calc2.0/internal/use_cases/post_expression"
	registrateuser "github.com/Gandalf-Rus/distributed-calc2.0/internal/use_cases/registrate_user"
	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type Orchestrator struct {
	server     *http.Server
	grpcServer *grpc.Server
	repo       *storage.Storage
}

func New(ctx context.Context) (*Orchestrator, error) {
	const (
		defaultHTTPServerWriteTimeout = time.Second * 15
		defaultHTTPServerReadTimeout  = time.Second * 15
	)

	orch := new(Orchestrator)
	router := mux.NewRouter()

	// загрузка конфига
	l.Logger.Info("initializing config...")
	err := config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	//привязываем мидлвейр
	router.Use(middlewares.LoggingMiddleware)

	// иницилизируем структуру для работы с базой
	repo, err := storage.New(ctx)
	l.Logger.Info("DB tables initialization...")
	if err != nil {
		l.Logger.Info("DB initialization failed")
		return nil, err
	}

	// если в базе нет таблиц, создадим их
	if err := repo.CreateTablesIfNotExist(); err != nil {
		return nil, err
	}
	l.Logger.Info("DB initialization succeeds")

	// подключение хендлеров к путям
	registerHandler := http.HandlerFunc(registrateuser.MakeHandler(registrateuser.NewSvc(&repo)))
	loginHandler := http.HandlerFunc(loginuser.MakeHandler(loginuser.NewSvc(&repo)))
	postExpressionHandler := http.HandlerFunc(postexpression.MakeHandler(postexpression.NewSvc(&repo)))

	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hi")) }).Methods("GET")
	apiRouter.HandleFunc("/register", registerHandler).Methods("POST")
	apiRouter.HandleFunc("/login", loginHandler).Methods("POST")
	apiRouter.HandleFunc("/expression", postExpressionHandler).Methods("POST")

	router.PathPrefix("/api").Handler(apiRouter)

	// http сервер
	orch.server = &http.Server{
		Handler:      router,
		Addr:         ":" + strconv.Itoa(config.Cfg.ServerPort),
		WriteTimeout: defaultHTTPServerWriteTimeout,
		ReadTimeout:  defaultHTTPServerReadTimeout,
	}

	// grpc сервер
	orch.grpcServer = grpc.NewServer()
	nodeServiceServer := geteditnodes.NewServer(&repo)
	proto.RegisterNodeServiceServer(orch.grpcServer, nodeServiceServer)

	// зачистка пропавших агентов
	agent.LostAgentCollector(&repo)

	// передаем структуру для работы с БД чтоб в конце закрыть подключение
	orch.repo = &repo

	return orch, nil
}

func (o *Orchestrator) RunServer() error {
	l.Logger.Info("starting http server...")
	l.Logger.Info(fmt.Sprintf("http listener started at: %v", o.server.Addr))
	err := o.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server was stop with err: %w", err)
	}

	l.Logger.Info("server was stop")
	return nil
}

func (o *Orchestrator) RunGrpcServer() error {
	l.Logger.Info("starting grpc server...")

	err := serveGrpcServer(o.grpcServer)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("grpc server was stop with err: %w", err)
	}

	l.Logger.Info("grpc server was stop")
	return nil
}

func (o *Orchestrator) stop(ctx context.Context) error {
	l.Logger.Info("shutdowning server...")
	err := o.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}
	o.repo.ClosePoolConn()
	o.grpcServer.Stop()
	l.Logger.Info("server was shutdown")
	return nil
}

func (o *Orchestrator) GracefulStop(serverCtx context.Context, sig <-chan os.Signal, serverStopCtx context.CancelFunc) {
	<-sig
	var timeOut = 1 * time.Second
	shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, timeOut)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			l.Logger.Error("graceful shutdown timed out... forcing exit")
			os.Exit(1)
		}
	}()

	err := o.stop(shutdownCtx)
	if err != nil {
		l.Logger.Error("graceful shutdown timed out... forcing exit")
		os.Exit(1)
	}
	serverStopCtx()
	shutdownStopCtx()
}

func serveGrpcServer(server *grpc.Server) error {
	addr := fmt.Sprintf("%s:%d", config.Cfg.GrpcHost, config.Cfg.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		l.Logger.Error("error starting tcp listener: " + err.Error())
		return err
	}
	l.Logger.Info(fmt.Sprintf("tcp listener started at: %v", addr))

	if err := server.Serve(lis); err != nil {
		l.Logger.Error("error serving grpc: " + err.Error())
		return err
	}
	return nil
}
