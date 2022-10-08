package gapi

import (
	"go-grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *PostServer) GetPosts(request *pb.GetPostsRequest, stream pb.PostService_GetPostsServer) error {
	page := request.GetPage()
	limit := request.GetLimit()

	posts, err := server.postService.FindPosts(int(page), int(limit))
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	for _, post := range posts {
		stream.Send(&pb.Post{
			Id:        post.Id.Hex(),
			Title:     post.Title,
			Content:   post.Content,
			Image:     post.Image,
			CreatedAt: timestamppb.New(post.CreateAt),
			UpdatedAt: timestamppb.New(post.UpdateAt),
		})
	}
	return nil
}
