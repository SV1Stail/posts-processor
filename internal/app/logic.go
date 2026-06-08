package app

import (
	"context"
	"time"

	"github.com/SV1Stail/posts-processor/internal/db"
	queue_scheduler_pb "github.com/SV1Stail/tg-project-protos/gen/go/queue_scheduler/queue_scheduler"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	publishChannel = "channel"
)

type Summary struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (pp *PostProcessor) ProcessPosts(ctx context.Context) error {
	log.Info().Ctx(ctx).Msg("start ProcessPosts ...")
	originalPosts, err := pp.DB.GetOriginalPosts(ctx, &db.GetOriginalPosts{
		Theme: "test",
	})
	if err != nil {
		return err
	}

	if len(originalPosts) <= 1 {
		log.Warn().Ctx(ctx).Msg("no posts to process")
		return nil
	}

	var userMessage string
	for _, post := range originalPosts {
		userMessage += " " + post.Data.Text
	}

	resp, err := pp.LlmClient.NewChat(ctx, summary_news, userMessage)
	if err != nil {
		return err
	}

	// raw := []byte(resp.Message)
	// raw = bytes.ReplaceAll(raw, []byte{'\n'}, []byte(``))
	// raw = bytes.ReplaceAll(raw, []byte{'\t'}, []byte(``))
	// raw = bytes.ReplaceAll(raw, []byte{'\r'}, []byte(``))

	post, err := pp.QueueSchedulerClient.CreatePost(ctx, &queue_scheduler_pb.CreatePostRequest{
		PublishChannel: publishChannel,
		Data: &queue_scheduler_pb.PublishPostData{
			Title:   "",
			Body:    resp,
			PostUrl: postUrls(originalPosts),
		},
		PublishAt: timestamppb.New(time.Now().UTC().Add(12 * time.Hour)),
	})
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("Create post failed")

		return err
	}
	log.Info().Ctx(ctx).Str("post_id", post.Post.Id).Msg("create post success")

	return nil
}

func postUrls(posts []*db.OriginalPost) []string {
	urls := make([]string, 0, 10)
	for _, post := range posts {
		urls = append(urls, *post.LinkOriginalPost)
	}

	return urls
}

func (pp *PostProcessor) pingCLient(ctx context.Context) error {
	err := pp.QueueSchedulerClient.CheckConnectionState(ctx)
	if err != nil {
		return err
	}

	return nil
}
