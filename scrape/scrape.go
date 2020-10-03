package scrape

//Getter interface
type Getter interface {
	GetAll() []Item
}

//Adder contract for adding element
type Adder interface {
	Add(item Item)
}

//Item struct
type Item struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Interval  string `json:"interval"`
	Response  string `json:"response"`
	Duration  string `json:"duration"`
	CreatedAt string `json:"created_at"`
}

//Repo struct
type Repo struct {
	Items []Item
}

//New create new repository
func New() *Repo {
	return &Repo{
		Items: []Item{},
	}
}

//Add adds new element to existing repository
func (r *Repo) Add(item Item) {
	r.Items = append(r.Items, item)
}

//GetAll gets all elements
func (r *Repo) GetAll() []Item {
	return r.Items
}
