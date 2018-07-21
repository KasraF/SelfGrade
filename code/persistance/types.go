package persistance

type Table interface {
	Exists() (bool, error)
	Find() (bool, error)
	Save() error
}

type User struct {
	Id int
	Name string
	Email string
	Password []byte
	Role string
}
