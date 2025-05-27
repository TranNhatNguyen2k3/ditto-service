package request

type BaseFilterRequest struct {
	PageSize  int    `json:"page_size"`
	PageIndex int    `json:"page_index"`
	Names     string `json:"name"`
	IDs       []uint `json:"id_includes"`
	FromDate  int64  `json:"from_date"`
	ToDate    int64  `json:"to_date"`
	Sort      string `json:"sort"`
}

func (b BaseFilterRequest) GetOffsetAndLimit() (limit, offset int) {
	if b.PageSize == 0 {
		return 20, 0
	}
	limit = b.PageSize
	offset = b.PageIndex * b.PageSize
	return
}
