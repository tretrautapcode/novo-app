package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dqhieuu/novo-app/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

const limitChapter = 50
const limitNameCharacter = 50

type Chapter struct {
	ChapterNumber float64 `json:"chapterNumber" binding:"required"`
	Name          string  `json:"name"`
	Id            int32   `json:"id" binding:"required"`
	TimePosted    int64   `json:"timePosted" binding:"required"`
	UserPosted    Author  `json:"userPosted" binding:"required"`
}

type HypertextChapter struct {
	ChapterNumber float64 `json:"chapterNumber" binding:"required"`
	Name          string  `json:"name"`
	TextContent   string  `json:"textContent" binding:"required"`
	BookGroupId   int32   `json:"bookGroupId" binding:"required"`
}

type ImageChapter struct {
	ChapterNumber float64 `json:"chapterNumber" binding:"required"`
	Name          string  `json:"name"`
	Images        []int32 `json:"images" binding:"required"`
	BookGroupId   int32   `json:"bookGroupId" binding:"required"`
}

func checkChapterName(name string) bool {
	if HasControlCharacters(name) {
		return false
	}
	hasNextLine, _ := regexp.MatchString(`[\r\n]`, name)
	if hasNextLine {
		return false
	}
	if len(name) > limitNameCharacter {
		return false
	}
	return true
}

func BookChapterById(id int32) (*db.BookChapter, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	bookChapter, err := queries.BookChapterById(ctx, id)

	if err != nil {
		stringErr := fmt.Sprintf("Get book Chapter by id failed: %s", err)
		return nil, errors.New(stringErr)
	}
	return &bookChapter, err
}

func BookChaptersByBookGroupId(bookGroupID, page int32) ([]*db.BookChapter, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	bookChapters, err := queries.BookChaptersByBookGroupId(ctx, db.BookChaptersByBookGroupIdParams{
		BookGroupID: bookGroupID,
		Offset:      (page - 1) * limitChapter,
		Limit:       limitChapter,
	})
	if err != nil {
		stringErr := fmt.Sprintf("Get bookChapters by book group id failed: %s", err)
		return nil, errors.New(stringErr)
	}
	var outData []*db.BookChapter
	for i := 0; i < len(bookChapters); i++ {
		outData = append(outData, &bookChapters[i])
	}
	return outData, err
}

func UpdateBookChapter(id int32, chapterNumber float64, description, textContext, chapterType string,
	bookGroupID, ownerID int32) error {

	ctx := context.Background()
	queries := db.New(db.Pool())

	descriptionSql := sql.NullString{}
	err := descriptionSql.Scan(description)
	if err != nil {
		stringErr := fmt.Sprintf("Update book chapter  failed: %s", err)
		return errors.New(stringErr)
	}

	textContextSql := sql.NullString{}
	err = textContextSql.Scan(textContext)
	if err != nil {
		stringErr := fmt.Sprintf("Update book chapter  failed: %s", err)
		return errors.New(stringErr)
	}

	err = queries.UpdateBookChapter(ctx, db.UpdateBookChapterParams{
		ID:            id,
		ChapterNumber: chapterNumber,
		Name:          descriptionSql,
		TextContext:   textContextSql,
		Type:          chapterType,
		BookGroupID:   bookGroupID,
		OwnerID:       ownerID,
	})
	if err != nil {
		stringErr := fmt.Sprintf("Update book chapter  failed: %s", err)
		return errors.New(stringErr)
	}
	return nil
}

func CreateBookChapter(chapterNumber float64, description, textContext, chapterType string,
	bookGroupID, ownerID int32) (*db.BookChapter, error) {

	ctx := context.Background()
	queries := db.New(db.Pool())

	descriptionSql := sql.NullString{}
	err := descriptionSql.Scan(description)
	if err != nil {
		stringErr := fmt.Sprintf("Create book chapter  failed: %s", err)
		return nil, errors.New(stringErr)
	}

	textContextSql := sql.NullString{}
	err = textContextSql.Scan(textContext)
	if err != nil {
		stringErr := fmt.Sprintf("Create book chapter  failed: %s", err)
		return nil, errors.New(stringErr)
	}

	bookChapter, err := queries.InsertBookChapter(ctx, db.InsertBookChapterParams{
		ChapterNumber: chapterNumber,
		Name:          descriptionSql,
		TextContext:   textContextSql,
		Type:          chapterType,
		BookGroupID:   bookGroupID,
		OwnerID:       ownerID,
	})
	if err != nil {
		stringErr := fmt.Sprintf("Create book chapter  failed: %s", err)
		return nil, errors.New(stringErr)
	}
	return &bookChapter, nil
}

func DeleteBookChapterById(id int32) error {
	ctx := context.Background()
	queries := db.New(db.Pool())
	err := queries.DeleteBookChapterById(ctx, id)
	if err != nil {
		stringErr := fmt.Sprintf("Delete book chapter by Id failed: %s", err)
		return errors.New(stringErr)
	}
	return nil
}

func DeleteBookChapterByBookGroupId(bookGroupId int32) error {
	ctx := context.Background()
	queries := db.New(db.Pool())
	err := queries.DeleteBookChapterByBookGroupId(ctx, bookGroupId)
	if err != nil {
		stringErr := fmt.Sprintf("Delete book chapter by bookGroupId failed: %s", err)
		return errors.New(stringErr)
	}
	return nil
}

func CreateHypertextChapterHandler(c *gin.Context) {
	var newHypertextChapter HypertextChapter
	if err := c.ShouldBindJSON(&newHypertextChapter); err != nil {
		ReportError(c, err, "error parsing json", http.StatusBadRequest)
		return
	}

	//check chapter name
	if !checkChapterName(newHypertextChapter.Name) {
		ReportError(c, errors.New("invalid chapter name"), "error", http.StatusBadRequest)
		return
	}

	//check content
	if HasControlCharacters(newHypertextChapter.TextContent) && CheckEmptyString(newHypertextChapter.TextContent) {
		ReportError(c, errors.New("invalid content"), "error", http.StatusBadRequest)
		return
	}

	extract := jwt.ExtractClaims(c)

	newChapter, err := CreateBookChapter(
		newHypertextChapter.ChapterNumber,
		newHypertextChapter.Name,
		newHypertextChapter.TextContent,
		"hypertext",
		newHypertextChapter.BookGroupId,
		int32(extract[UserIdClaimKey].(float64)))

	if err != nil {
		ReportError(c, err, "error creating new hypertext chapter", 500)
		return
	}
	c.JSON(200, gin.H{
		"id": newChapter.ID,
	})
}

func CreateImagesChapterHandler(c *gin.Context) {
	ctx := context.Background()
	queries := db.New(db.Pool())

	var newImageChapter ImageChapter
	if err := c.ShouldBindJSON(&newImageChapter); err != nil {
		ReportError(c, err,"error parsing json", http.StatusBadRequest)
		return
	}

	//check chapter name
	if !checkChapterName(newImageChapter.Name) {
		ReportError(c, errors.New("invalid chapter name"), "error", http.StatusBadRequest)
		return
	}

	extract := jwt.ExtractClaims(c)

	newChapter, err := CreateBookChapter(
		newImageChapter.ChapterNumber,
		newImageChapter.Name,
		"",
		"images",
		newImageChapter.BookGroupId,
		int32(extract[UserIdClaimKey].(float64)))

	if err != nil {
		ReportError(c, err, "error creating new images chapter", 500)
		return
	}

	//adding images
	for index, imageId := range newImageChapter.Images {
		peekRow, err := queries.GetImageBasedOnId(ctx, imageId)
		switch {
		case len(peekRow.Md5) == 0 || len(peekRow.Sha1) == 0:
			ReportError(c, errors.New("image does not exist"), "error", http.StatusBadRequest)
		case err == nil:
			err = queries.InsertBookChapterImage(ctx, db.InsertBookChapterImageParams{
				BookChapterID: newChapter.ID,
				ImageID:       imageId,
				Rank:          int32(index + 1),
			})
			if err != nil {
				ReportError(c, err, "error adding image chapter", 500)
				return
			}
		default:
			ReportError(c, err, "error getting image", 500)
			return
		}
	}

	c.JSON(200, gin.H{
		"id": newChapter.ID,
	})
}

func GetBookChapterContentHandler(c *gin.Context) {
	var chapterId int32
	_, err := fmt.Sscan(c.Param("chapterId"), &chapterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookChapter, err := BookChapterById(chapterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if bookChapter.Type == "images" {
		images, err := ImagesByBookChapter(chapterId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if *images == nil {
			images = &[]string{}
		}
		c.JSON(http.StatusOK, gin.H{
			"type":          bookChapter.Type,
			"bookGroupId":   bookChapter.BookGroupID,
			"chapterNumber": bookChapter.ChapterNumber,
			"name":          bookChapter.Name.String,
			"images":        *images,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"type":          bookChapter.Type,
			"bookGroupId":   bookChapter.BookGroupID,
			"chapterNumber": bookChapter.ChapterNumber,
			"name":          bookChapter.Name.String,
			"textContent":   bookChapter.TextContext.String,
		})
	}
}

func DeleteBookChapterHandler(c *gin.Context) {
	var chapterId int32
	_, err := fmt.Sscan(c.Param("chapterId"), &chapterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldChapter, err := BookChapterById(chapterId)
	if oldChapter == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Chapter not exist",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = DeleteBookChapterById(chapterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Chapter successfully",
	})
}

func LatestCreatedInBookGroup(bookGroupId int32) (*db.LastChapterInBookGroupRow, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	chapterNumber, err := queries.LastChapterInBookGroup(ctx, bookGroupId)
	if err != nil {
		return nil, err
	}
	return &chapterNumber, nil
}
