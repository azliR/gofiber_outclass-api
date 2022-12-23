package dtos

import "time"

type CreateTaskDto struct {
	ClassroomId string    `form:"classroom_id" validate:"omitempty,len=24"`
	Title       string    `form:"title" validate:"required"`
	Details     *string   `form:"detaDetails"`
	Date        time.Time `form:"date" validate:"required"`
	Repeat      string    `form:"repeat" validate:"required,oneof=none daily weekly monthly yearly"`
	Color       string    `form:"color" validate:"required,oneof=maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
}
