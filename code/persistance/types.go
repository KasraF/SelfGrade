package persistance

type Table interface {
	Exists() (bool, error)
	Find()   (bool, error)
	Save()   error
	Update() error
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password []byte
	Role     string
}

type App struct {
	Id          int
	Name        string
	Description string
	Url         string
	IconUrl     string
}
