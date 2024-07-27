package types

type Post struct {
	ID            int
	Title         string
	Content       string
	ShortContent  string
	Owner         int
	Timestamp     string
	OwnerUsername string
	TotalVotes    int
	Upvotes       int
	Downvotes     int
}
