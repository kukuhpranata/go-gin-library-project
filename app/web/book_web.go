package web

type BookCreate struct {
	Title           string `validate:"required" json:"title"`
	Author          string `validate:"required" json:"author"`
	Isbn            string `validate:"required" json:"isbn"`
	PublicationYear int    `validate:"required" json:"publication_year"`
	Quantity        int    `validate:"required" json:"quantity"`
}

type BookUpdate struct {
	Id              int    `validate:"required" json:"id"`
	Title           string `validate:"required" json:"title"`
	Author          string `validate:"required" json:"author"`
	Isbn            string `validate:"required" json:"isbn"`
	PublicationYear int    `validate:"required" json:"publication_year"`
	Quantity        int    `validate:"required" json:"quantity"`
}

type BookResponse struct {
	Id              int    `json:"id"`
	Title           string `json:"name"`
	Author          string `json:"author"`
	Isbn            string `json:"isbn"`
	PublicationYear int    `json:"publication_year"`
	Quantity        int    `json:"quantity"`
}
