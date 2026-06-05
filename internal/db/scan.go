package db

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func scanOriginalPost(row pgx.Row) (*OriginalPost, error) {
	originalPost := &OriginalPost{}
	data := make([]byte, 0)
	err := row.Scan(
		&originalPost.ID,
		&data,
		&originalPost.LinkOriginalPost,
		&originalPost.OriginalChannel,
		&originalPost.URL,
		&originalPost.OriginalImageURL,
		&originalPost.Theme,
		&originalPost.CreatedAt,
		&originalPost.UpdatedAt,
		&originalPost.Attempts,
	)
	if err != nil {
		log.Err(err).Msg("scan failed")
		return nil, err
	}
	originalPost.Data, err = convertByteIntoOriginalPostData(data)
	if err != nil {
		log.Err(err).
			Str("original_post_id", originalPost.ID).
			Str("linked_post", *originalPost.LinkOriginalPost).
			Msg("failed unmarshal data after scan")

		return nil, err
	}

	return originalPost, nil
}

func scanPosts(rows pgx.Rows) ([]*OriginalPost, error) {
	var originalPosts []*OriginalPost
	data := make([]byte, 0)
	for rows.Next() {
		originalPost := &OriginalPost{}
		err := rows.Scan(
			&originalPost.ID,
			&data,
			&originalPost.LinkOriginalPost,
			&originalPost.OriginalChannel,
			&originalPost.URL,
			&originalPost.OriginalImageURL,
			&originalPost.Theme,
			&originalPost.CreatedAt,
			&originalPost.UpdatedAt,
			&originalPost.Attempts,
		)
		if err != nil {
			log.Err(err).Msg("scan failed")
			return nil, err
		}
		originalPost.Data, err = convertByteIntoOriginalPostData(data)
		if err != nil {
			log.Err(err).
				Str("original_post_id", originalPost.ID).
				Str("linked_post", *originalPost.LinkOriginalPost).
				Msg("failed unmarshal data after scan")

			continue
		}

		originalPosts = append(originalPosts, originalPost)
	}

	if err := rows.Err(); err != nil {
		log.Err(err).Msg("rows error")
		return nil, err
	}

	return originalPosts, nil
}

func convertByteIntoOriginalPostData(data []byte) (*OriginalPostData, error) {
	var text string
	err := json.Unmarshal(data, &text)
	if err != nil {
		return nil, err
	}

	return &OriginalPostData{
		Text: text,
	}, nil
}
