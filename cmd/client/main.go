package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	v1 "news-grpc/api/news/v1"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial server: %v", err)
	}
	defer conn.Close()

	client := v1.NewNewsServiceClient(conn)
	ctx := context.Background()

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("validator init: %v", err)
	}

	// ================= VALIDATION TEST =================
	fmt.Println("=== VALIDATION TEST WITH EMPTY REQUEST ===")
	emptyReq := &v1.CreateRequest{}
	if err := validator.Validate(emptyReq); err != nil {
		fmt.Println("Validation errors for empty request:")
		fmt.Println(err)
	} else {
		fmt.Println("Empty request passed validation (unexpected).")
	}

	// ================= CREATE NEWS =================
	fmt.Println("\n=== CREATING VALID NEWS ===")
	createdNews := make([]*v1.CreateRequest, 0)

	for i := 0; i < 5; i++ {
		req := &v1.CreateRequest{
			Id:      uuid.NewString(),
			Author:  fmt.Sprintf("Author %d", i),
			Title:   fmt.Sprintf("Breaking News Title %d", i),
			Content: fmt.Sprintf("This is the full content for news item number %d. It is long enough to pass validation rules set in proto.", i),
			Summary: fmt.Sprintf("Summary of news item number %d with enough length.", i),
			Source:  "https://example.com",
			Tags:    []string{"tag1", "tag2"},
		}

		if err := validator.Validate(req); err != nil {
			log.Fatalf("validation error: %v", err)
		}

		if _, err := client.Create(ctx, req); err != nil {
			log.Fatalf("create news: %v", err)
		}

		createdNews = append(createdNews, req)
		fmt.Printf("[CREATED] ID=%s | Title=%q | Author=%q\n", req.Id, req.Title, req.Author)
	}

	// ================= GET ALL NEWS =================
	fmt.Println("\n=== FETCHING ALL NEWS ===")
	getAllStream, err := client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("get all: %v", err)
	}

	allNews := make([]*v1.GetAllResponse, 0)
	for {
		res, err := getAllStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("get all recv: %v", err)
		}
		allNews = append(allNews, res)
		printNews(res)
	}

	// ================= UPDATE NEWS =================
	fmt.Println("\n=== UPDATING NEWS ===")
	updateStream, err := client.UpdateNews(ctx)
	if err != nil {
		log.Fatalf("update stream: %v", err)
	}

	for i, n := range allNews {
		req := &v1.CreateRequest{
			Id:      n.Id,
			Author:  fmt.Sprintf("%s updated %d", n.Author, i),
			Title:   n.Title + " (Updated)",
			Content: n.Content,
			Summary: n.Summary,
			Source:  n.Source,
			Tags:    n.Tags,
		}
		if err := updateStream.Send(req); err != nil {
			log.Fatalf("update send: %v", err)
		}
		fmt.Printf("[UPDATED] ID=%s | New Author=%q | New Title=%q\n", req.Id, req.Author, req.Title)
	}

	if _, err := updateStream.CloseAndRecv(); err != nil {
		log.Fatalf("update close: %v", err)
	}

	// ================= DELETE NEWS =================
	fmt.Println("\n=== DELETING NEWS ===")
	deleteStream, err := client.DeletedNews(ctx)
	if err != nil {
		log.Fatalf("delete stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		defer close(waitc)
		for _, n := range allNews {
			if err := deleteStream.Send(&v1.NewsID{Id: n.Id}); err != nil {
				log.Fatalf("delete send: %v", err)
			}
		}
		deleteStream.CloseSend()
	}()

	for {
		_, err := deleteStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("delete recv: %v", err)
		}
	}

	<-waitc
	fmt.Println("[DELETED ALL NEWS]")

	// ================= FINAL GET =================
	fmt.Println("\n=== FETCHING REMAINING NEWS ===")
	finalStream, err := client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("final get: %v", err)
	}

	for {
		n, err := finalStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("final recv: %v", err)
		}
		printNews(n)
	}
}

// ---------- HELPER FUNCTION TO PRINT NEWS ----------
func printNews(n *v1.GetAllResponse) {
	fmt.Printf("ID: %s\nTitle: %s\nAuthor: %s\nSummary: %s\nContent: %s\nSource: %s\nTags: %v\n----------------------\n",
		n.Id, n.Title, n.Author, n.Summary, n.Content, n.Source, n.Tags)
}
