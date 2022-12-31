package main

import ( 
	"context"
	"fmt"
	"log"
	"net"
	// "net/http"
	"os"
	// "os/signal"
	// "syscall"
	// "time"

	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/config"
	// "golang.org/x/sync/errgroup"

)

func main() {
	// err := http.ListenAndServe(
	// 	":18080",
	// 	http.HandlerFunc( func( w http.ResponseWriter, r *http.Request ) {
	// 		fmt.Fprintf( w, "Hello, %s!", r.URL.Path[1:] )
	// 	}),
	// )
	// if err != nil {
	// 	fmt.Printf( "failed to terminate server: %v", err)
	// 	os.Exit(1)
	// }
	// if len(os.Args) != 2{
	// 	log.Printf("need port number\n")
	// 	os.Exit(1)
	// }

	// p := os.Args[1]
	// l, err := net.Listen("tcp", ":"+p)
	// if err != nil{
	// 	log.Fatalf("failed to listen port %s: %v", p,err)
	// }

	// if err := run(context.Background(), l); err != nil {
	// 	log.Printf("failed to terminate server: %v", err)
	// 	os.Exit(1)
	// }	
	if err := run(context.Background() ); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}	


}

func run(ctx context.Context) error {
	// ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	// mux := NewMux()
	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()
	s := NewServer(l, mux)
	return s.Run(ctx)

	// server.go に移設

	// s := &http.Server{
	// 	// Addr: ":18080",
	// 	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		// コマンドラインで実験するため
	// 		time.Sleep(5 * time.Second)
	// 		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	// 	}),
	// }
	// eg, ctx := errgroup.WithContext(ctx)
	// // 別ゴルーチンでHTTPサーバーを起動する
	// eg.Go(func() error {
	// 	// http.ErrServerClosed は
	// 	// http.Server.Shutdown() が正常に終了したことを示すので異常ではない。
	// 	// if err := s.ListenAndServe(); err != nil &&
	// 	if err := s.Serve(l); err != nil &&
	// 		err != http.ErrServerClosed {
	// 		log.Printf("failed to close: %+v", err)
	// 		return err
	// 	}
	// 	return nil
	// })
	// // チャネルからの通知（終了通知）を待機する
	// <-ctx.Done()
	// if err := s.Shutdown(context.Background()); err != nil {
	// 	log.Printf("failed to shutdown: %+v", err)
	// }
	// // Goメソッドで起動した別ゴルーチンの終了を待つ。
	// return eg.Wait()
}




