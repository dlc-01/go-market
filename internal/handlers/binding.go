package handlers

import (
	"github.com/dlc/go-market/internal/model/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Field string
	Value string
	Tag   string
	Param string
}

func BindData(ginC *gin.Context, req interface{}) ([]invalidArgument, error) {

	if err := CheckContentType(ginC, "application/json"); err != nil {

		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return nil, apperrors.NewUnsupportedMediaType("content type not support")
	}

	if err := ginC.ShouldBind(req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {

			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			return invalidArgs, apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")
		}

		return nil, apperrors.NewInternal()
	}
	return nil, nil
}
