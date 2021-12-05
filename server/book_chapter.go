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
	"strings"
)

const limitChapter = 50
const limitNameCharacter = 50

type Chapter struct {
	ChapterNumber float64     `json:"chapterNumber"`
	Name          interface{} `json:"name"`
	Id            int32       `json:"id" binding:"required"`
	TimePosted    int64       `json:"timePosted" binding:"required"`
	UserPosted    Author      `json:"userPosted" binding:"required"`
}

type HypertextChapter struct {
	ChapterNumber interface{} `json:"chapterNumber"`
	Name          interface{} `json:"name"`
	TextContent   string      `json:"textContent" binding:"required"`
	BookGroupId   int32       `json:"bookGroupId" binding:"required"`
}

type ImageChapter struct {
	ChapterNumber interface{} `json:"chapterNumber"`
	Name          interface{} `json:"name"`
	Images        []int32     `json:"images" binding:"required"`
	BookGroupId   int32       `json:"bookGroupId" binding:"required"`
}

type UpdateHypertextChapterParams struct {
	Id            int32   `json:"id"`
	ChapterNumber float64 `json:"chapterNumber"`
	Name          string  `json:"name"`
	TextContent   string  `json:"textContent"`
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

func UpdateBookChapter(chapter UpdateHypertextChapterParams) error {

	ctx := context.Background()
	queries := db.New(db.Pool())

	nameSql := sql.NullString{}
	err := nameSql.Scan(chapter.Name)
	if err != nil {
		stringErr := fmt.Sprintf("Update book chapter  failed: %s", err)
		return errors.New(stringErr)
	}

	textContextSql := sql.NullString{}
	err = textContextSql.Scan(chapter.TextContent)
	if err != nil {
		stringErr := fmt.Sprintf("Update book chapter  failed: %s", err)
		return errors.New(stringErr)
	}

	err = queries.UpdateBookChapter(ctx, db.UpdateBookChapterParams{
		ID:            chapter.Id,
		ChapterNumber: chapter.ChapterNumber,
		Name:          nameSql,
		TextContext:   textContextSql,
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

	if newHypertextChapter.ChapterNumber == nil {
		ReportError(c, errors.New("missing chapter number"), "error", http.StatusBadRequest)
		return
	}

	//check chapter name
	var chapterName string
	if newHypertextChapter.Name != nil {
		chapterName = newHypertextChapter.Name.(string)
		chapterName = strings.TrimSpace(chapterName)
		if !checkChapterName(chapterName) || CheckEmptyString(chapterName) {
			ReportError(c, errors.New("invalid chapter name"), "error", http.StatusBadRequest)
			return
		}
	}

	//check content
	if HasControlCharacters(newHypertextChapter.TextContent) || CheckEmptyString(newHypertextChapter.TextContent) {
		ReportError(c, errors.New("invalid content"), "error", http.StatusBadRequest)
		return
	}

	extract := jwt.ExtractClaims(c)

	newChapter, err := CreateBookChapter(
		newHypertextChapter.ChapterNumber.(float64),
		chapterName,
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
		ReportError(c, err, "error parsing json", http.StatusBadRequest)
		return
	}

	if newImageChapter.ChapterNumber == nil {
		ReportError(c, errors.New("missing chapter number"), "error", http.StatusBadRequest)
		return
	}

	//check chapter name
	var chapterName string
	if newImageChapter.Name != nil {
		chapterName = newImageChapter.Name.(string)
		chapterName = strings.TrimSpace(chapterName)
		if !checkChapterName(chapterName) || CheckEmptyString(chapterName) {
			ReportError(c, errors.New("invalid chapter name"), "error", http.StatusBadRequest)
			return
		}
	}

	extract := jwt.ExtractClaims(c)

	newChapter, err := CreateBookChapter(
		newImageChapter.ChapterNumber.(float64),
		chapterName,
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
		check, err := queries.CheckImageExistById(ctx, imageId)
		if err != nil {
			ReportError(c, err, "internal error", 500)
			return
		}
		if check {
			err = queries.InsertBookChapterImage(ctx, db.InsertBookChapterImageParams{
				BookChapterID: newChapter.ID,
				ImageID:       imageId,
				Rank:          int32(index + 1),
			})
			if err != nil {
				ReportError(c, err, "error adding image chapter", 500)
				return
			}
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

	ctx := context.Background()
	queries := db.New(db.Pool())
	err = queries.DeleteBookChapterById(ctx, chapterId)
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

func UpdateHypertextChapter(c *gin.Context) {
	var newChapter UpdateHypertextChapterParams
	_, err := fmt.Sscan(c.Param("chapterId"), &newChapter.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error chapterId": err.Error()})
		return
	}
	oldChapter, err := BookChapterById(newChapter.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if oldChapter == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Chapter not exist",
		})
		return
	}
	if oldChapter.Type != "hypertext" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Type of chapter is not hypertext",
		})
		return
	}

	if err = c.ShouldBindJSON(&newChapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error json": err.Error()})
		return
	}
	if newChapter.ChapterNumber < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error json": "chapter id cannot be less than 0"})
		return
	}
	if newChapter.ChapterNumber == 0 {
		newChapter.ChapterNumber = oldChapter.ChapterNumber
	}
	if newChapter.Name == "" {
		newChapter.Name = oldChapter.Name.String
	}
	if newChapter.TextContent == "" {
		newChapter.TextContent = oldChapter.TextContext.String
	}
	newChapter.Name = strings.TrimSpace(newChapter.Name)
	if checkChapterName(newChapter.Name) == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name cannot contain control characters or newline characters",
		})
		return
	}
	if HasControlCharacters(newChapter.TextContent) == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "TextContent cannot contain control characters",
		})
		return
	}
	err = UpdateBookChapter(newChapter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update chapter successfully",
	})
}
