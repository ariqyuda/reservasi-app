package repositories

import "math"

type Pagination struct {
	CurrentPage     int  `json:"current_page"`
	PageSize        int  `json:"page_size"`
	TotalPage       int  `json:"total_page"`
	HasNextPage     bool `json:"has_next_page"`
	LastVisiblePage int  `json:"last_visible_page"`
	TotalData       int  `json:"total_data"`
}

func GetDataPageInfo(page int, pageSize int, totalRows int) Pagination {
	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
	HasNextPage := page < totalPages

	return Pagination{
		CurrentPage:     page,
		PageSize:        pageSize,
		TotalPage:       totalPages,
		HasNextPage:     HasNextPage,
		LastVisiblePage: int(math.Min(float64(page+2), float64(totalPages))),
		TotalData:       totalRows,
	}
}
