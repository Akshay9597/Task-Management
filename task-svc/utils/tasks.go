package utils

import "time"

type(
	Task struct{
		Id	int	`json:"id" db:"id"`
		Title string `json:"title" db:"title" binding:"required,min=3"`
		CreationTime time.Time `json:"creation_time" db:"creation_time"`
		UserId int `json:"user_id" db:"user_id"`
	}
)


