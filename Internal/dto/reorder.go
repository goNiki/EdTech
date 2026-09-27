package dto

type ReorderItemsRequest struct {
	ItemIDs []int64 `json:"item_ids" validate:"required,min=1"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=draft published archived"`
}
