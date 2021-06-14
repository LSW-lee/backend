package dblayer

import (
	"backend/src/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type DBLayer interface {
	GetAllProducts() ([]models.Product, error)
	GetPromos() ([]models.Product, error)
	GetCustomerByName(string, string) (models.Customer, error)
	GetCustomerByID(int) (models.Customer, error)
	GetProduct(int) (models.Product, error)
	AddUser(models.Customer) (models.Customer, error)
	SignInUser(username, password string) (models.Customer, error)
	SignOutUserById(int) error
	GetCustomerOrdersByID(int) ([]models.Order, error)
	AddOrder(models.Order) error
	GetCreditCardCID(int) (string, error)
	SaveCreditCardForCustomer(int, string) error
}

var ErrINVALIDPASSWORD = errors.New("INVALID PASSWORD")

func hashPassword(s *string) error {
	if s == nil {
		return errors.New("REFERENCE PROVIDED FOR HASHING PASSWORD IS NIL")
	}
	//bcrypt 패키지에서 사용할 수 있게 패스워드 문자열을 바이스 슬라이스로 변환
	sBytes := []byte(*s)

	//GeneratorFromPassword() 메서드는 패스워드 해시를 반환
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//패스워드 문자열을 해시 값으로 바꾼다.
	*s = string(hashedBytes[:])
	return nil //에러가 없어서 nil 반환한 듯
}

func checkPassword(existingHash, incomingPass string) bool {
	//해시와 패스워드 문자열이 일치하지 않으면 에러를 반환하기 때문에 nil과 비교값 return
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
