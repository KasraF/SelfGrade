package persistance

type Table interface {
	Find() error
	Save() error
}

type User struct {
	Id int
	Name string
	Email string
	Password string
	Role string
}
