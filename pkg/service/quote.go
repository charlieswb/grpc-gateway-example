package service

import (
	"context"
	"math/rand"
	pb "restgrpc/pkg/protobuf/message"
)

type quoteService struct {
	pb.UnimplementedQuoteServiceServer
	quotes map[string][]string
}

func NewQuoteService() pb.QuoteServiceServer {
	s := quoteService{
		quotes: make(map[string][]string),
	}
	s.quotes["Gandhi"] = []string{"Be the change you wish to see", "Action expresses priorities"}
	s.quotes["Maxwell"] = []string{"Sometimes you win, sometimes you learn", "There is no such thing as luck"}
	s.quotes["Obi-wan"] = []string{"Hello there", "Another happy landing"}
	return &s
}

func (s *quoteService) Echo(ctx context.Context, req *pb.StringMessage) (*pb.StringMessage, error) {
	res := pb.StringMessage{
		Text: "Hello",
	}
	if req != nil && req.Text != "" {
		res.Text = req.Text
	}
	return &res, nil
}

func (s *quoteService) GetQuote(ctx context.Context, req *pb.QuoteRequest) (*pb.QuoteReply, error) {
	author := "Gandhi"
	if req != nil && req.Author != "" {
		author = req.Author
	}
	var res pb.QuoteReply
	if quotes, ok := s.quotes[author]; ok {
		index := rand.Intn(2)
		res = pb.QuoteReply{
			Quote:  quotes[index],
			Author: author,
		}
	} else {
		res = pb.QuoteReply{
			Quote:  "404 Not Found",
			Author: "Internet Assigned Numbers Authority (IANA)",
		}
	}
	return &res, nil
}
