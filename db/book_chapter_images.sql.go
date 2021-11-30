// Code generated by sqlc. DO NOT EDIT.
// source: book_chapter_images.sql

package db

import (
	"context"
)

const imagesByBookChapter = `-- name: ImagesByBookChapter :many
SELECT i.path
FROM book_chapter_images AS bci
JOIN images AS i ON i.id=bci.image_id
WHERE bci.book_chapter_id = $1
ORDER BY bci.rank ASC
`

func (q *Queries) ImagesByBookChapter(ctx context.Context, bookChapterID int32) ([]string, error) {
	rows, err := q.db.Query(ctx, imagesByBookChapter, bookChapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		items = append(items, path)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertBookChapterImage = `-- name: InsertBookChapterImage :exec
INSERT INTO book_chapter_images(book_chapter_id, image_id, rank) VALUES($1, $2, $3)
`

type InsertBookChapterImageParams struct {
	BookChapterID int32 `json:"bookChapterID"`
	ImageID       int32 `json:"imageID"`
	Rank          int32 `json:"rank"`
}

func (q *Queries) InsertBookChapterImage(ctx context.Context, arg InsertBookChapterImageParams) error {
	_, err := q.db.Exec(ctx, insertBookChapterImage, arg.BookChapterID, arg.ImageID, arg.Rank)
	return err
}
