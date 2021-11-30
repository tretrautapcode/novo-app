package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dqhieuu/novo-app/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

const limitBookAuthors = 50

type Author struct {
	Name string `json:"name" binding:"required"`
	Id int32 `json:"id" binding:"required"`
}

func BookAuthorById(id int32) (*db.BookAuthor, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	bookAuthor, err := queries.BookAuthorById(ctx, id)

	if err != nil {
		stringErr := fmt.Sprintf("Get bookAuthor by id failed: %s", err)
		return nil, errors.New(stringErr)
	}
	return &bookAuthor, err
}

func BookAuthors(page int32) ([]*db.BookAuthor, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	bookAuthors, err := queries.BookAuthors(ctx, db.BookAuthorsParams{
		Offset: (page - 1) * limitBookAuthors,
		Limit:  limitBookAuthors,
	})
	if err != nil {
		stringErr := fmt.Sprintf("Get bookAuthors list failed: %s", err)
		return nil, errors.New(stringErr)
	}
	var outData []*db.BookAuthor
	for i := 0; i < len(bookAuthors); i++ {
		outData = append(outData, &bookAuthors[i])
	}
	return outData, err
}

func UpdateBookAuthor(id int32, name, description string, imageID int32) error {

	ctx := context.Background()
	queries := db.New(db.Pool())

	var descriptionSql sql.NullString
	err := descriptionSql.Scan(description)
	if err != nil {
		stringErr := fmt.Sprintf("Update bookAuthor failed: %s", err)
		return errors.New(stringErr)
	}

	var imageIdSql = sql.NullInt32{
		Int32: imageID,
		Valid: imageID > 0,
	}

	err = queries.UpdateBookAuthor(ctx, db.UpdateBookAuthorParams{
		ID:            id,
		Name:          name,
		Description:   descriptionSql,
		AvatarImageID: imageIdSql,
	})
	if err != nil {
		stringErr := fmt.Sprintf("Update bookAuthor failed: %s", err)
		return errors.New(stringErr)
	}
	return nil
}

func CreateBookAuthor(name, description string, imageID int32) (*db.BookAuthor, error) {

	ctx := context.Background()
	queries := db.New(db.Pool())

	var descriptionSql sql.NullString
	err := descriptionSql.Scan(description)
	if err != nil {
		stringErr := fmt.Sprintf("Update bookAuthor failed: %s", err)
		return nil, errors.New(stringErr)
	}

	var imageIdSql = sql.NullInt32{
		Int32: imageID,
		Valid: imageID > 0,
	}

	bookAuthor, err := queries.InsertBookAuthor(ctx, db.InsertBookAuthorParams{
		Name:          name,
		Description:   descriptionSql,
		AvatarImageID: imageIdSql,
	})
	if err != nil {
		stringErr := fmt.Sprintf("Create bookAuthor failed: %s", err)
		return nil, errors.New(stringErr)
	}
	return &bookAuthor, nil
}

func DeleteBookAuthor(id int32) error {
	var err error
	err = DeleteBookGroupsByBookAuthor(id)
	if err != nil {
		stringErr := fmt.Sprintf("Delete bookAuthor failed: %s", err)
		return errors.New(stringErr)
	}
	ctx := context.Background()
	queries := db.New(db.Pool())
	err = queries.DeleteBookAuthor(ctx, id)
	if err != nil {
		stringErr := fmt.Sprintf("Delete bookAuthor failed: %s", err)
		return errors.New(stringErr)
	}
	return nil
}

func CheckAuthorExistByName(name string) (bool, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	result, err := queries.CheckAuthorExistByName(ctx, name)
	if err != nil {
		return false, err
	}
	return result, nil
}

func CheckAuthorExistById(id int32) (bool, error) {
	ctx := context.Background()
	queries := db.New(db.Pool())
	result, err := queries.CheckAuthorExistById(ctx, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

type CreateAuthor struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	AvatarId    int32  `json:"avatarId" form:"avatarId"`
}

func CreateAuthorHandler(c *gin.Context) {
	var a CreateAuthor
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(a.Name) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name must be less than or equal to 30 characters",
		})
		return
	}

	if len(a.Description) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "description must be less than or equal to 50 characters",
		})
		return
	}

	exist, err := CheckAuthorExistByName(a.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if exist == true {
		c.JSON(http.StatusConflict, gin.H{
			"error": "name was exist",
		})
		return
	}

	_, err = CreateBookAuthor(a.Name, a.Description, a.AvatarId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Create Author successfully",
	})
}

type UpdateAuthor struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	AvatarId    int32  `json:"avatar" form:"avatar"`
}

func UpdateAuthorHandler(c *gin.Context) {
	var authorId int32
	_, err := fmt.Sscan(c.Param("authorId"), &authorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldAuthor, err := BookAuthorById(authorId)
	if oldAuthor == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Author not exist",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var a UpdateAuthor
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(a.Name) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name must be less than or equal to 30 characters",
		})
		return
	}
	if len(a.Description) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "description must be less than or equal to 50 characters",
		})
		return
	}
	exist, err := CheckAuthorExistByName(a.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exist == true {
		c.JSON(http.StatusConflict, gin.H{"error": "name was exist"})
		return
	}

	if len(a.Name) == 0 {
		a.Name = oldAuthor.Name
	}
	if len(a.Description) == 0 {
		a.Description = oldAuthor.Description.String
	}
	if a.AvatarId == 0 {
		a.AvatarId = oldAuthor.AvatarImageID.Int32
	}
	err = UpdateBookAuthor(authorId, a.Name, a.Description, a.AvatarId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Update Author successfully",
	})
}

func DeleteAuthorHandler(c *gin.Context) {
	var authorId int32
	_, err := fmt.Sscan(c.Param("authorId"), &authorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldAuthor, err := BookAuthorById(authorId)
	if oldAuthor == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Author not exist",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = DeleteBookAuthor(authorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Author successfully",
	})
}
