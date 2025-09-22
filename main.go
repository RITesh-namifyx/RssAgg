package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RITesh-namifyx/RssAgg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	// feed,err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

	fmt.Println("Hello World");

	godotenv.Load()
	port := os.Getenv("PORT")

	if port == ""{
		log.Fatal("PORT not found in env")
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == ""{
		log.Fatal("Database URL not found in env")
	}

	conn, err := sql.Open("postgres",dbUrl)
	if err != nil{
		log.Fatal("Can't connect to DATABASE:",err)
	}

	
	db := database.New(conn)
	apicfg := apiConfig{
		DB: db,
	}

	go startScrapping(db,10,time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handlerErr)
	v1Router.Post("/users",apicfg.handlerCreateUser)
	v1Router.Get("/users",apicfg.middlewareAuth(apicfg.handlerGetUser))
	v1Router.Post("/feed", apicfg.middlewareAuth(apicfg.handlerCreateFeed))
	v1Router.Get("/feeds", apicfg.handlerGetFeeds)
	v1Router.Post("/feedFollow",apicfg.middlewareAuth(apicfg.handlerCreateFeedFollow))
	v1Router.Get("/feedFollow",apicfg.middlewareAuth(apicfg.GetAllFeedFollows))
	v1Router.Delete("/feedFollow/{feedFollowID}", apicfg.middlewareAuth(apicfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts",apicfg.middlewareAuth(apicfg.handlerGetPostsForUser))
	
	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

	log.Printf("server strated at : %v ",port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("err:",err)

	}
	fmt.Println("PORT :", port);
}