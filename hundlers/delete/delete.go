package delete

// type DeleteRequest struct {
// 	ID             *int  `json:"id"`
// 	DeleteRelation *bool `json:"delete_relation"`
// }

type DeleteResponseFailure struct {
	Message string `json:"message"`
}

// func deleteRequestCheck(req DeleteRequest) error {

// 	if req.ID == nil {
// 		return fmt.Errorf("id is null")
// 	}
// 	if req.DeleteRelation == nil {
// 		return fmt.Errorf("delete_relation is null")
// 	}

// 	return nil
// }
