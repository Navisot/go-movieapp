package main

import (
	"context"
	"fmt"
	"github.com/navisot/movieapp/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {

	// connect to grpc server
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("cannot dial: %v", err)
	}

	// close the connection when needed
	defer conn.Close()

	// initialize s movie service client
	movieClient := pb.NewMovieServiceClient(conn)

	// run methods
	runGetMovies(movieClient)
	runGetMovie(movieClient)
	runUpdateMovie(movieClient)
	runCreateMovie(movieClient)
	runDeleteMovie(movieClient)
}

// runGetMovies retrieves all the movies
func runGetMovies(client pb.MovieServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetMovies(ctx, req)

	if err != nil {
		log.Fatalf("runGetMovies error : %v", err)
	}

	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("runGetMovies error: %v", err)
		}
		log.Printf("Received a movie from the server: Title = %s, ISBN = %s, ID = %s \n", row.GetTitle(), row.GetIsbn(), row.GetId())
	}
}

// runGetMovie retrieves a single movie
func runGetMovie(client pb.MovieServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.Id{Value: "1"}

	movie, err := client.GetMovie(ctx, req)

	if err != nil {
		log.Printf("cannot get movie: %v", err)
	}

	log.Printf("The movie from the server is ID: %s, Title: %s and the director is "+
		"%s %s", movie.GetId(), movie.GetTitle(), movie.Director.GetFirstname(), movie.Director.GetLastname())
}

// runUpdateMovie updates a movie
func runUpdateMovie(client pb.MovieServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movieId := &pb.Id{Value: "1"}

	movie, err := client.GetMovie(ctx, movieId)

	if err != nil {
		log.Printf("cannot get movie: %v", err)
	}

	// Update movie title
	movie.Title = "New updated title"

	status, err := client.UpdateMovie(ctx, movie)

	if err != nil {
		log.Printf("cannot update movie: %v", err)
	}

	if status.Value == 1 {
		fmt.Println("movie updated")
	}

	log.Printf("New movie title is : %s", movie.GetTitle())
}

// runCreateMovie creates a new movie
func runCreateMovie(client pb.MovieServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newMovie := &pb.MovieInfo{
		Title: "Pulp Fiction",
		Isbn:  "010900",
		Director: &pb.Director{
			Firstname: "Quentin",
			Lastname:  "Tarantino",
		},
	}

	id, err := client.CreateMovie(ctx, newMovie)

	if err != nil {
		log.Println("cannot create movie")
	}

	log.Printf("movie created with movie id: %s", id)

	runGetMovies(client)
}

// runDeleteMovie deletes a movie
func runDeleteMovie(client pb.MovieServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := client.DeleteMovie(ctx, &pb.Id{Value: "1"})

	if err != nil {
		log.Println("cannot delete movie")
	}

	if status.Value == 1 {
		log.Printf("movie with id %s deleted", "1")
	} else {
		log.Println("movie not found for delete it")
	}

	runGetMovies(client)
}
