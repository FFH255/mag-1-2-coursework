package proto

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error" example:"something went wrong"`
} //@name ResponseError

func WriteError(c *gin.Context, status int, err any) {
	var response Error

	switch e := err.(type) {
	case error:
		response = Error{Error: e.Error()}
	case string:
		response = Error{Error: e}
	}

	c.AbortWithStatusJSON(status, response)
}

func WriteJSON[V any](c *gin.Context, status int, v V) {
	c.JSON(status, v)
}

func MarshalTime(t time.Time) string {
	format := time.RFC3339
	timeString := t.Format(format)

	return timeString
}

func UnmarshalTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func MustUnmarshalDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}

	return d
}

func ParseMilliseconds(milliseconds *uint64) time.Duration {
	if milliseconds == nil {
		return 0
	}

	return time.Duration(*milliseconds) * time.Millisecond
}

func UnmarshalTimeFromQueryParams(c *gin.Context, name string) (*time.Time, error) {
	timeString := c.Query(name)
	if timeString == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func GetInt64FromQueryParam(c *gin.Context, key string) (*int64, error) {
	value := c.Query(key)
	if value == "" {
		return nil, nil
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &i, nil
}
