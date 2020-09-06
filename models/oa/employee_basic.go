/*
* Auth : acer
* Desc : 员工基本信息
* Time : 2020/9/1 21:27
 */

package oa

import (
	"bfimpl/models"

	"github.com/jinzhu/gorm"
)

type EmployeeBasic struct {
	gorm.Model
	EmployeeID              int         `json:"employee_id"`
	IDCardFront             string      `gorm:"not null;comment:'身份证正'" json:"id_card_front"`
	IDCardBack              string      `gorm:"not null;comment:'身份证反面'" json:"id_card_back"`
	DebitCard1              string      `gorm:"size:60;not null;comment:'工资卡1'" json:"debit_card1"`
	IssuingBank1            string      `gorm:"size:60;not null;comment:'发卡行1'" json:"issuing_bank1"`
	IssuingBank1Detail      string      `gorm:"size:500;not null;comment:'发卡行1详情'" json:"issuing_bank1_detail"`
	DebitCard1Front         string      `gorm:"comment:'发卡行1正面'" json:"debit_card1_front"`
	DebitCard2              string      `gorm:"size:60;not null;comment:'工资卡2'" json:"debit_card2"`
	IssuingBank2            string      `gorm:"size:60;not null;comment:'发卡行2'" json:"issuing_bank2"`
	IssuingBank2Detail      string      `gorm:"size:500;not null;comment:'发卡行2详情'" json:"issuing_bank2_detail"`
	DebitCard2Front         string      `gorm:"comment:'发卡行2正面'" json:"debit_card2_front"`
	SocialSecurity          string      `gorm:"comment:'社保信息'" json:"social_security"`
	PublicFund              string      `gorm:"comment:'公积金信息'" json:"public_fund"`
	Degree                  string      `gorm:"size:20;not null;comment:'学历'" json:"degree"`
	DegreeProperty          string      `gorm:"size:20;not null;comment:'学历性质'" json:"degree_property"`
	Major                   string      `gorm:"size:20;not null;comment:'专业'" json:"major"`
	GraduationSchool        string      `gorm:"size:100;not null;comment:'毕业院校'" json:"graduation_school"`
	DegreeCertificationCopy string      `gorm:"comment:'毕业证书复印'" json:"degree_certification_copy"`
	DegreeVerification      string      `gorm:"size:20;not null;comment:'学历验证(未验证,已验证,无法验证)'" json:"degree_verification"`
	ENSkill                 string      `gorm:"size:20;comment:'英语技能'" json:"en_skill"`
	OtherLanguageSkill      string      `gorm:"size:100;comment:'其他语言'" json:"other_language_skill"`
	Birthday                models.Time `gorm:"type:datetime;not null;comment:'生日'" json:"birthday"`
	Birthplace              string      `gorm:"size:30;not null;comment:'籍贯'" json:"birthplace"`
	InhabitedCity           string      `gorm:"size:30;not null;comment:'居住城市'" json:"inhabited_city"`
	InhabitedDistrict       string      `gorm:"size:30;not null;comment:'区'" json:"inhabited_district"`
	InhabitedAddress        string      `gorm:"not null;comment:'地址'" json:"inhabited_address"`
	Marriage                string      `gorm:"size:10;not null;comment:'婚姻状况'" json:"marriage"`
	Children                int         `gorm:"not null;comment:'子女数'" json:"children"`
	FatherName              string      `gorm:"size:20;comment:'父名'" json:"father_name"`
	FatherTel               string      `gorm:"size:20;comment:'父联系方式'" json:"father_tel"`
	FatherCareer            string      `gorm:"size:50;comment:'父职'" json:"father_career"`
	FatherCom               string      `gorm:"size:50;comment:'企业'" json:"father_com"`
	MotherName              string      `gorm:"size:20;comment:'母名'" json:"mother_name"`
	MotherTel               string      `gorm:"size:20;comment:'联系方式'" json:"mother_tel"`
	MotherCareer            string      `gorm:"size:50;comment:'母职'" json:"mother_career"`
	MotherCom               string      `gorm:"size:50;comment:'公司'" json:"mother_com"`
	SpouseName              string      `gorm:"size:20;comment:'配偶名'" json:"spouse_name"`
	SpouseCareer            string      `gorm:"size:50;comment:'配偶职'" json:"spouse_career"`
	SpouseCom               string      `gorm:"size:50;comment:'公司'" json:"spouse_com"`
	SpouseTel               string      `gorm:"size:20;comment:'联系方式'" json:"spouse_tel"`
	EmergencyContact1       string      `gorm:"size:20;not null;comment:'紧急联系1'" json:"emergency_contact_1"`
	Contact1Relation        string      `gorm:"size:20;not null;comment:'联系1关系'" json:"contact_1_relation"`
	EmergencyMobile1        string      `gorm:"size:20;not null;comment:'紧急联系1tel'" json:"emergency_mobile_1"`
	EmergencyContact2       string      `gorm:"size:20;not null;comment:'紧急联系2'" json:"emergency_contact_2"`
	EmergencyMobile2        string      `gorm:"size:20;not null;comment:'紧急联系2tel'" json:"emergency_mobile_2"`
	Contact2Relation        string      `gorm:"size:20;not null;comment:'联系2关系'" json:"contact_2_relation"`
}
