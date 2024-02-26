package Simulation

type Movie struct {
	Rank         int     `json:"rank"`
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Year         int     `json:"year"`
	IMDbVotes    int     `json:"imdb_votes"`
	IMDbRating   float64 `json:"imdb_rating"`
	Certificate  string  `json:"certificate"`
	Duration     int     `json:"duration"`
	Genre        string  `json:"genre"`
	CastID       string  `json:"cast_id"`
	CastName     string  `json:"cast_name"`
	DirectorID   string  `json:"director_id"`
	DirectorName string  `json:"director_name"`
	WriterName   string  `json:"writer_name"`
	WriterID     string  `json:"writer_id"`
	ImageLink    string  `json:"img_link"`
}

// dummyMovies is our sample data
var dummyMovies = []Movie{
	{
		Rank:         1,
		ID:           "tt0111161",
		Name:         "The Shawshank Redemption",
		Year:         1994,
		IMDbVotes:    2601152,
		IMDbRating:   9.3,
		Certificate:  "A",
		Duration:     142,
		Genre:        "Drama",
		CastID:       "nm0000209,nm0000151,nm0348409,nm0006669,nm0000317,nm0004743,nm0001679,nm0926235,nm0218810,nm0104594,nm0321358,nm0508742,nm0698998,nm0706554,nm0161980,nm0005204,nm0086169,nm0542957",
		CastName:     "Tim Robbins,Morgan Freeman,Bob Gunton,William Sadler,Clancy Brown,Gil Bellows,Mark Rolston,James Whitmore,Jeffrey DeMunn,Larry Brandenburg,Neil Giuntoli,Brian Libby,David Proval,Joseph Ragno,Jude Ciccolella,Paul McCrane,Renee Blaine,Scott Mann",
		DirectorID:   "nm0001104",
		DirectorName: "Frank Darabont",
		WriterName:   "Stephen King,Frank Darabont",
		WriterID:     "nm0000175,nm0001104",
		ImageLink:    "https://m.media-amazon.com/images/M/MV5BMWU4N2FjNzYtNTVkNC00NzQ0LTg0MjAtYTJlMjFhNGUxZDFmXkEyXkFqcGdeQXVyNjc1NTYyMjg@._V1_QL75_UX380_CR0",
	},
}

func GetMovies() []Movie {
	return dummyMovies
}
