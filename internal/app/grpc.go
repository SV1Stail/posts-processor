package app

import (
	"context"

	post_processor_pb "github.com/SV1Stail/tg-project-protos/gen/go/posts_processor"
)

func (pp *PostProcessor) CreateOriginalPost(ctx context.Context,
	req *post_processor_pb.CreateOriginalPostRequest,
) (*post_processor_pb.CreatePostCreateOriginalPostResponse, error) {
	err := pp.DB.CreateOriginalPost(ctx, req)

	return &post_processor_pb.CreatePostCreateOriginalPostResponse{}, err
}

func (pp *PostProcessor) CreateOriginalPosts(ctx context.Context,
	req *post_processor_pb.CreateOriginalPostsRequest,
) (
	*post_processor_pb.CreatePostCreateOriginalPostsResponse, error) {
	return nil, nil
}

func (pp *PostProcessor) DeleteOriginalPostsByChannel(ctx context.Context,
	req *post_processor_pb.DeleteOriginalPostsByChannelRequest,
) (*post_processor_pb.DeleteOriginalPostsByChannelResponse, error) {
	err := pp.DB.DeleteOriginalPostsByChannel(ctx, req)

	return &post_processor_pb.DeleteOriginalPostsByChannelResponse{}, err
}

func (pp *PostProcessor) GetOriginalPostsFromDB(ctx context.Context,
	req *post_processor_pb.GetOriginalPostsFromDBRequest,
) (*post_processor_pb.GetOriginalPostsFromDBResponse, error) {
	originalPosts, err := pp.DB.GetOriginalPostsFromDB(ctx, req)

	return originalPosts, err
}
