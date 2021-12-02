package handlers

import (
	"business/storage"
	"business/types"
	"business/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type StaffHandler struct {
	staff storage.Staff
}

func NewStaffHandler(s storage.Staff) StaffHandler {
	return StaffHandler{
		staff: s,
	}
}

var TokenString string

func (s StaffHandler) AutoMigrate(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var d types.Product
	var err error
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = s.staff.AutoMigrate(d)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create table..")
	}
	json.NewEncoder(w).Encode("Table created successfully")

}

func (s StaffHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var d types.Staff
	var err error
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to unmarshal JSON")
	}
	// var e validator.FieldError
	// e = Translator(d, w)
	// if e != nil {
	// 	return
	// }
	err = s.staff.SignUp(d)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to signup a user")
	}
	json.NewEncoder(w).Encode("User created successfully")
}

func (s StaffHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var d types.Staff
	var re types.Staff
	var err error
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to Unmarshal JSON")
		return
	}
	re, err = s.staff.SignIn(d)
	if err != nil {
		json.NewEncoder(w).Encode("Incorrect user name or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(re.Password), []byte(d.Password))
	if err != nil {
		json.NewEncoder(w).Encode("Unable to compare password")
		return
	}

	td := &types.TokenDetails{}
	td, err = createToken(int64(d.ID), d.FullName)
	TokenString = td.AccessToken
	res := map[string]string{
		"accessT":  td.AccessToken,
		"refreshT": td.RefreshToken,
		"username": d.FullName,
	}
	json.NewEncoder(w).Encode(res)

}

func (s StaffHandler) Save(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var d types.Product
	var err error

	var data types.Product
	var clm *types.Claims

	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to decode JSON")
	}
	clm, err = verifyToken(w, r)
	data, err = s.staff.Save(d, clm.FullName)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to persist data to DB..")
	}
	json.NewEncoder(w).Encode(data)

}

func (s StaffHandler) Product(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var err error
	var d types.Product
	var result types.Product
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to unmarshal JSON")
	}
	var clm *types.Claims
	clm, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Expired session")
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", clm.FullName)))
	result, err = s.staff.Product(d)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(result)
}

func (s StaffHandler) Products(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var data []types.Product
	var err error
	var clm *types.Claims
	clm, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to verify Token")
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", clm.FullName)))
	data, err = s.staff.Products()
	if err != nil {
		json.NewEncoder(w).Encode("No products")
	}
	json.NewEncoder(w).Encode(data)

}

func createToken(id int64, fname string) (*types.TokenDetails, error) {
	var err error
	td := &types.TokenDetails{}
	td.AtExp = time.Now().Add(time.Minute * 15)
	td.RfExp = time.Now().Add(time.Hour * 24 * 7)

	aClm := types.Claims{
		FullName: fname,
		Id:       id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.AtExp.Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, aClm)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	rClm := types.Claims{
		FullName: fname,
		Id:       id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.RfExp.Unix(),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rClm)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

func verifyToken(w http.ResponseWriter, r *http.Request) (*types.Claims, error) {
	tString := TokenString
	clm := &types.Claims{}

	tkn, err := jwt.ParseWithClaims(tString, clm, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		json.NewEncoder(w).Encode("Invalid token string or claims")
	}

	if !tkn.Valid {
		json.NewEncoder(w).Encode("Token expired")
	}
	return clm, nil

}

func Translator(data types.Staff, w http.ResponseWriter) validator.FieldError {
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("Translator not fiund")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("passwd", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} is not strong enough", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})

	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	_ = v.RegisterTranslation("gender", trans, func(ut ut.Translator) error {
		return ut.Add("gender", "{0} not correct gender", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gender", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("unit", trans, func(ut ut.Translator) error {
		return ut.Add("unit", "{0} not correct unit", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unit", fe.Field())
		return t
	})

	//var err error

	err := v.Struct(data)
	var errs validator.FieldError

	for _, errs = range err.(validator.ValidationErrors) {
		json.NewEncoder(w).Encode(errs.Translate(trans))
		fmt.Println(errs.Translate(trans))
	}
	return errs

}
