package dtos

import "time"

type CreateEventDto struct {
	ClassroomId string     `json:"classroom_id" validate:"omitempty,len=24"`
	Title       string     `json:"title" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date"`
	Repeat      string     `json:"repeat" validate:"required,oneof=none daily weekly monthly yearly"`
	Color       string     `json:"color" validate:"required,oneof=maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
	Description *string    `json:"description"`
}
