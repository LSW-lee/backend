package dblayer

import (
	"backend/src/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBORM struct {
	*gorm.DB
}

func NewORM(dbname, con string) (*DBORM, error) {
	log.Print("dbname :", dbname, "/con : ", con)
	db, err := gorm.Open(mysql.Open(con), &gorm.Config{})
	return &DBORM{
		DB: db,
	}, err
}

//select * from products 의 결과와 같다
func (db *DBORM) GetAllProducts() (products []models.Product, err error) {
	return products, db.Find(&products).Error
}

//select * from products where promotion is not null
func (db *DBORM) GetPromos() (products []models.Product, err error) {
	return products, db.Where("promotion IS NOT NULL").Find(&products).Error
}

//selct * from customers where firstname='...' and lastname='...'
func (db *DBORM) GetCustomerByName(firstname string, lastname string) (customer models.Customer, err error) {
	return customer, db.Where(&models.Customer{FirstName: firstname, LastName: lastname}).Find(&customer).Error
}

//어차피 id는 기본키이기 때문에 .Where과 .Find로 해도 동일한 결과가 나올 것이다.
func (db *DBORM) GetCustomerByID(id int) (customer models.Customer, err error) {
	return customer, db.First(&customer, id).Error
}

func (db *DBORM) GetProduct(id int) (product models.Product, err error) {
	return product, db.First(&product, id).Error
}

func (db *DBORM) AddUser(customer models.Customer) (models.Customer, error) {
	//패스워드를 해시 값으로 저장하고자 레퍼런스를 넘긴다
	hashPassword(&customer.Pass)
	customer.LoggedIn = true
	err := db.Create(&customer).Error
	//보안을 위해서 객체를 반환하기 전 패스워드 문자열 제거
	customer.Pass = ""
	return customer, err
}

func (db *DBORM) SignInUser(email, pass string) (customer models.Customer, err error) {
	//defer log.Println(&customer)
	//사용자 행을 나타내는 *gorm.DB 타입 할당
	result := db.Table("customers").Where(&models.Customer{Email: email})
	//입력된 이메일로 사용자 정보 조회
	err = result.First(&customer).Error
	if err != nil {
		return customer, err
	}

	//패스워드 문자열과 해시 값 비교
	if !checkPassword(customer.Pass, pass) {
		return customer, ErrINVALIDPASSWORD
	}

	//공유되지 않도록 패스워드 문자열은 지운다
	customer.Pass = ""
	//ㅣoggedin 필드 업데이트
	err = result.Update("loggedin", 1).Error
	if err != nil {
		return customer, err
	}

	customer.Model.ID = 0
	//사용자 행 반환
	return customer, result.Find(&customer).Error
}

func (db *DBORM) SignOutUserById(id int) error {
	//ID에 해당하는 사용자 구조체 생성
	// customer := models.Customer{
	// 	Model: gorm.Model{
	// 		ID: uint(id),
	// 	},
	// }
	//로그아웃 상태로 업데이트
	return db.Table("customers").Where("id=?", id).Update("loggedin", 0).Error
}

/*
* orders, customers, products 테이블 조인
* customers 테이블에서 전달받은 id 값에 해당하는 사용자 정보를 조회
* products 테이블에서 현재 선택된 상품 ID에 해당하는 상품 정보를 검색
 */
func (db *DBORM) GetCustomerOrdersByID(id int) (orders []models.Order, err error) {
	return orders, db.Table("orders").Select("*").
		Joins("join customers on customers.id = customer_id").
		Joins("join products on products.id = product_id").
		Where("customer_id=?", id).Scan(&orders).Error
}

//orders 테이블에 결제 내역 추가
func (db *DBORM) AddOrder(order models.Order) error {
	var simple models.OrderSimple
	simple.Model = order.Model
	simple.CustomerID = order.CustomerID
	simple.ProductID = order.ProductID
	simple.Price = order.Price
	simple.PurchaseDate = order.PurchaseDate
	return db.Table("orders").Create(&simple).Error
}

//신용카드 ID 조회
func (db *DBORM) GetCreditCardCID(id int) (string, error) {
	customerWithCCID := struct {
		models.Customer
		CCID string `gorm:"column:cc_customerid"`
	}{}
	return customerWithCCID.CCID, db.First(&customerWithCCID, id).Error
}

func (db *DBORM) SaveCreditCardForCustomer(id int, ccid string) error {
	result := db.Table("customers").Where("id=?", id)
	return result.Update("cc_customerid", ccid).Error
}
