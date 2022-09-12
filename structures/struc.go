package structures

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorLog struct {
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
	ID           string             `json:"id"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	WorkPlace    string             `json:"workplace"`
	Expereance   string             `json:"expereance"`
	Biography    string             `json:"biography"`
	History      history
	Password     string `json:"password"`
	Login        string `json:"login"`
	RequestLogin string `json:"requestLogin"`
	Name         string `json:"name"`
	Sername      string `json:"sername"`
	Possition    string `json:"possition"`
	ProfileImage string `json:"profileimg"`
	Permissions  string `json:"permissions"`
}
type ClientLog struct {
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
	ID           string             `json:"id"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	Password     string             `json:"password"`
	Login        string             `json:"login"`
	RequestLogin string             `json:"requestLogin"`
	Name         string             `json:"name"`
	Sername      string             `json:"sername"`
	LastName     string             `json:"lastname"`
	Gender       string             `json:"gender"`
	Blood        string             `json:"blood"`
	ProfileImage string             `json:"profileimg"`
	Permissions  string             `json:"permissions"`
}
type history struct {
	Period    string `json:"period"`
	ShortInfo string `json:"shortinfo"`
	Location  string `json:"location"`
}
type FrequentlyAskedQuestion struct {
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	RequestLogin string             `json:"requestLogin"`
}
type Users struct {
	Password     string `json:"password"`
	Login        string `json:"login"`
	Permissions  string `json:"permissions"`
	Id           string `json:"id"`
	LastActive   string `json:"lastactive"` 
	UserStatus   string `json:"userStatus"`
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
}
type AdminLog struct {
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
	ID           string             `json:"id"`
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	RequestLogin string             `json:"requestLogin"`
	Login        string             `json:"login"`
	Name         string             `json:"name"`
	Sername      string             `json:"sername"`
	ProfileImage string             `json:"profileimg"`
	Permissions  string             `json:"permissions"`
}
type Online struct {
	Login       string `json:"login"`
	Permissions string `json:"permissions"`
	LastActive  string `json:"lastactive"` 
	UserStatus  string `json:"userStatus"`
}
type Results struct {
	Sickness     string `json:"sickness"`
	Diagnosis    string `json:"diagnosis"`
	RequestLogin string `json:"requestLogin"`
	Comments     string `json:"comments"`
}
type ViewReq struct {
	Sickness      string `json:"sickness"`
	Date          string `json:"date"`
	Status        string `json:"status"`
	Phone         string `json:"phone"`
	ClientId      string `json:"clientid"`
	DoctorId      string `json:"doctorid"`
	ClientName    string `json:"clientname"`
	ClientSername string `json:"clientsername"`
}
type Accept_Decline struct {
	Date     string `json:"date"`
	ClientId string `json:"clientid"`
}
type Dr_get_views struct {
	Id           string `json:"id"`
	RequestLogin string `json:"requestLogin"`
}
type Permission_get struct {
	RequestLogin string `json:"requestLogin"`
}

// structures that will be given to user for usage
type DoctorLogScreen struct {
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
	ID           string             `json:"id"`
	Email        string             `json:"email"`
	WorkPlace    string             `json:"workplace"`
	Expereance   string             `json:"expereance"`
	Biography    string             `json:"biography"`
	History      history
	Name         string             `json:"name"`
	Sername      string             `json:"sername"`
	Possition    string             `json:"possition"`
	ProfileImage string             `json:"profileimg"`
}
type ClientLogScreen struct {
	PrimitiveID  primitive.ObjectID `json:"primitiveid" bson:"_id"`
	ID           string             `json:"id"`
	Email        string             `json:"email"`
	Name         string             `json:"name"`
	Sername      string             `json:"sername"`
	LastName     string             `json:"lastname"`
	Gender       string             `json:"gender"`
	Blood        string             `json:"blood"`
	ProfileImage string             `json:"profileimg"`
}