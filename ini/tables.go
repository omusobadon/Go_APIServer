package ini

type Shop struct {
	ID      *int    `json:"id"`
	Name    *string `json:"name"`
	Mail    *string `json:"mail"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}
