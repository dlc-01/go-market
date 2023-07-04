package handlers

import (
	"github.com/dlc/go-market/internal/logger"
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

func BindData(ginC *gin.Context, req interface{}) (error, []invalidArgument) {

	if err := CheckContentType(ginC, "application/json"); err != nil {

		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return apperrors.NewUnsupportedMediaType("content type not support"), []invalidArgument{}
	}

	if err := ginC.ShouldBind(req); err != nil {
		logger.Errorf("error binding data: %s", err)
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

			return apperrors.NewBadRequest("Invalid request parameters. See invalidArgs"), invalidArgs
		}

		return apperrors.NewInternal(), []invalidArgument{}
	}
	return nil, []invalidArgument{}
}
