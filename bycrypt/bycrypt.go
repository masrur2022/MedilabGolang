package bycrypt

import (
	"fmt"
	"errors"
	"net"
	"regexp"
	"strings"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"io/ioutil"
)

func HashPassword(password string) (string, error){
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}
func CompareHashPasswords(HashedPasswordFromDB, PasswordToCampare string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPasswordFromDB), []byte(PasswordToCampare))
	return err == nil
}
// Compere the loginstring
func ChecktheLogin( LoginForCheck string) (string){
	UpercaseLaters := [26]string{"A","B","C","D","E","F","G","H","I","G","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
	LowercaseLaters := [26]string{"a","b","c","d","e","f","g","h","i","g","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
	// """""""""""split the login string into arrey"""""""""""
	login := strings.Split(LoginForCheck, "")
	//""""""""""""""""" sellect one letter from login string """""""""""""""""
	for i := 0; i < len(login); i++{
		//"""""""""""""""" loop through Upercase leters and define if there is any same leters""""""""""""""""
		for u := 0; u < len(UpercaseLaters); u++{
			if UpercaseLaters[u] == login[i]{
				// """""""""""""if there is one replase it with lowercase clone"""""""""""""
				login[i] = LowercaseLaters[u]
			}
		}
	}

	return strings.Join(login, "")
}
func ParseFile(c *gin.Context,directory string) (string){
	// """"""""""""""""""get the img""""""""""""""""""
	// upload of 10MB files
	c.Request.ParseMultipartForm(10 << 20)
	// formFiles haeders
	files, handler, err := c.Request.FormFile("img")
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer files.Close()
	fmt.Printf("File Name %s\n", handler.Filename)
	fmt.Printf("File Size %v\n", handler.Size)
	fmt.Printf("File Header %s\n", handler.Header)
	// create temporary files within the folder
	tempFiles, err := ioutil.TempFile(directory, "upload-*.jpg")
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer tempFiles.Close()
	// read all files to upload
	fileByte, err := ioutil.ReadAll(files)
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	// write  the byte arrey into temp files
	tempFiles.Write((fileByte))
	idString := strings.Split(tempFiles.Name(), "upload")[2]
	fmt.Println(idString)
	return idString
}
// """""""""""""""""""""""""""""""email velidation"""""""""""""""""""""""""""""""
var emailRegexp = regexp.MustCompile(`(?m)^(((((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?"((\s? +)?(([!#-[\]-~])|(\\([ -~]|\s))))*(\s? +)?"))?)?(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?<(((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?"((\s? +)?(([!#-[\]-~])|(\\([ -~]|\s))))*(\s? +)?"))@((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?\[((\s? +)?([!-Z^-~]))*(\s? +)?\]((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)))>((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?))|(((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?"((\s? +)?(([!#-[\]-~])|(\\([ -~]|\s))))*(\s? +)?"))@((((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?(([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+(\.([A-Za-z0-9!#-'*+\/=?^_\x60{|}~-])+)*)((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?)|(((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?\[((\s? +)?([!-Z^-~]))*(\s? +)?\]((((\s? +)?(\(((\s? +)?(([!-'*-[\]-~]*)|(\\([ -~]|\s))))*(\s? +)?\)))(\s? +)?)|(\s? +))?))))$`)
var ErrUnresolvableHost = errors.New("unresolvable host")
var ErrBadFormat        = errors.New("invalid format")
func ValidateMX(email string) error {
	host := strings.Split(email, "@")[1]
	if _, err := net.LookupMX(host); err != nil {
		return ErrUnresolvableHost
	}
	fmt.Println(host)
	return nil
}
func ValidateFormat(email string) bool {
	if !emailRegexp.MatchString(strings.ToLower(email)) {
		return false
	}
	return true
}