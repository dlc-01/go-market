package handlers

import (
	"github.com/dlc/go-market/internal/logger"
	"github.com/dlc/go-market/internal/model/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func BindData(ginC *gin.Context, req interface{}) {

	if err := CheckContentType(ginC, "application/json"); err != nil {
		logger.Error(" content type not support")
		ginC.AbortWithStatusJSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
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

			err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

			ginC.AbortWithStatusJSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return
		}

		fallBack := apperrors.NewInternal()

		ginC.AbortWithStatusJSON(fallBack.Status(), gin.H{"error": fallBack})
		return
	}

}
