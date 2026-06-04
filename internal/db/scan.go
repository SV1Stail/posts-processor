package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func scanOriginalPost(row pgx.Row) (*OriginalPost, error) {
	originalPost := &OriginalPost{}
	data := make([]byte, 0)
	err := row.Scan(
		&originalPost.ID,
		&data,
		&originalPost.LinkNewPost,
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
	originalPost.Data = convertByteIntoOriginalPostData(data)

	return originalPost, nil
}

func convertByteIntoOriginalPostData(data []byte) *OriginalPostData {
	return nil
}

func scanPosts(rows pgx.Rows) ([]*OriginalPost, error) {
	var originalPosts []*OriginalPost
	data := make([]byte, 0)
	for rows.Next() {
		originalPost := &OriginalPost{}
		err := rows.Scan(
			&originalPost.ID,
			&data,
			&originalPost.LinkNewPost,
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
		originalPost.Data = convertByteIntoOriginalPostData(data)

		originalPosts = append(originalPosts, originalPost)
	}

	if err := rows.Err(); err != nil {
		log.Err(err).Msg("rows error")
		return nil, err
	}

	return originalPosts, nil
}
