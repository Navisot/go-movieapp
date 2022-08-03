package main

import (
	"context"
	"github.com/navisot/movieapp/pb"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"strconv"
)

const (
	port = ":50051"
)

// moviesCollections acts like a DB
var moviesCollection []*pb.MovieInfo

// Init movieServer
type movieServer struct{}

// NewMovieServer returns a new movie server
func NewMovieServer() pb.MovieServiceServer {
	return &movieServer{}
}

// GetMovies accepts an empty request and starts streaming the movies back to the client
func (s *movieServer) GetMovies(in *pb.Empty, stream pb.MovieService_GetMoviesServer) error {

	// Log request
	log.Printf("GetMovies Request: %v", in)

	// Loop through the slice of movies and return the movies
	for _, movie := range moviesCollection {
		// Handle error
		if err := stream.Send(movie); err != nil {
			return err
		}
	}

	// All good
	return nil
}

// GetMovie returns a single movie info back to the client
func (s *movieServer) GetMovie(ctx context.Context, movieId *pb.Id) (*pb.MovieInfo, error) {

	// Log request
	log.Printf("GetMovie Request: %v", movieId)

	// Init a movie
	res := &pb.MovieInfo{}

	// Loop through the slice of movies and return the specific movie if exists
	for _, movie := range moviesCollection {
		if movie.GetId() == movieId.GetValue() {
			res = movie
			break
		}
	}

	// Return
	return res, nil
}

// CreateMovie stores a new movie into movies collection
func (s *movieServer) CreateMovie(ctx context.Context, movie *pb.MovieInfo) (*pb.Id, error) {

	// Log request
	log.Printf("CreateMovie Request: %v", movie)

	// Init a response
	res := &pb.Id{}
	res.Value = strconv.Itoa(rand.Intn(100000))

	// Assign movie id to the random generated value
	movie.Id = res.GetValue()

	// Append the movie to the movies collection
	moviesCollection = append(moviesCollection, movie)

	// Return
	return res, nil
}

// UpdateMovie updates a given movie
func (s *movieServer) UpdateMovie(ctx context.Context, movie *pb.MovieInfo) (*pb.Status, error) {
	// Log request
	log.Printf("UpdateMovie Request: %v", movie)

	// Init response
	res := &pb.Status{}

	// Loop through the slice of movies and update the specific movie if exists
	for i, m := range moviesCollection {
		if m.GetId() == movie.GetId() {
			moviesCollection = append(moviesCollection[:i], moviesCollection[i+1:]...)
			movie.Id = m.GetId()
			moviesCollection = append(moviesCollection, movie)
			res.Value = 1
			break
		}
	}

	// Return
	return res, nil
}

// DeleteMovie removes a movie from the movie collection
func (s *movieServer) DeleteMovie(ctx context.Context, movieId *pb.Id) (*pb.Status, error) {

	// Log request
	log.Printf("DeleteMovie Request: %v", movieId)

	// Init response
	st := &pb.Status{Value: 0}

	// Loop through the slice of movies and delete the specific movie if exists
	for i, m := range moviesCollection {
		if m.GetId() == movieId.GetValue() {
			moviesCollection = append(moviesCollection[:i], moviesCollection[i+1:]...)
			st.Value = 1
			break
		}
	}

	// Return
	return st, nil
}

func main() {

	initMovies()

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterMovieServiceServer(grpcServer, &movieServer{})

	log.Printf("Server has successfully started on port %s", port)

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("cannot start the grpc server: %v", err)
	}

}

// initMovies appends some movies in the movies collection
func initMovies() {
	movie1 := &pb.MovieInfo{
		Id:    "1",
		Isbn:  "054322",
		Title: "Halftime",
		Director: &pb.Director{
			Firstname: "Amanda",
			Lastname:  "Micheli",
		}}

	movie2 := &pb.MovieInfo{
		Id:    "2",
		Isbn:  "054452",
		Title: "Stardust",
		Director: &pb.Director{
			Firstname: "Gabriel",
			Lastname:  "Range",
		}}

	movie3 := &pb.MovieInfo{
		Id:    "3",
		Isbn:  "430928",
		Title: "Red",
		Director: &pb.Director{
			Firstname: "Robert",
			Lastname:  "Schwentke",
		}}

	moviesCollection = append(moviesCollection, movie1)
	moviesCollection = append(moviesCollection, movie2)
	moviesCollection = append(moviesCollection, movie3)
}
