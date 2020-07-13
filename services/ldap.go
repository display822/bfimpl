/*
* Auth : acer
* Desc : ldap
* Time : 2020/7/10 9:58
 */

package services

import (
	"fmt"

	"bfimpl/services/log"

	"github.com/go-ldap/ldap"
)

var ldapService *LDAPService

type LDAPConfig struct {
	Addr         string
	BindUserName string
	BindPassword string
	SearchDN     string
}

type LDAPService struct {
	Conn   *ldap.Conn
	Config LDAPConfig
}

func NewLDAPService(config LDAPConfig) (*LDAPService, error) {

	conn, err := ldap.Dial("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	err = conn.Bind(config.BindUserName, config.BindPassword)
	if err != nil {
		return nil, err
	}

	return &LDAPService{Conn: conn, Config: config}, nil
}

func LdapService() *LDAPService {
	if ldapService == nil {
		config := LDAPConfig{
			Addr:         "172.16.9.230:389",
			BindUserName: "CN=测试,OU=BF-IT,OU=BF-Users,DC=broadfun,DC=cn",
			BindPassword: "123456q@",
			SearchDN:     "OU=BF-Users,DC=broadfun,DC=cn",
		}
		s, err := NewLDAPService(config)
		if err != nil {
			log.GLogger.Error(err.Error())
			return nil
		}
		ldapService = s
	}
	return ldapService
}

// Login 登录
func (l *LDAPService) Login(userName, password string) (bool, error) {
	searchRequest := ldap.NewSearchRequest(
		l.Config.SearchDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(sAMAccountName=%s))", userName),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, fmt.Errorf("user does not exist or too many entries return")
	}

	userDN := sr.Entries[0].DN
	fmt.Println(userDN)
	err = l.Conn.Bind(userDN, password)
	if err != nil {
		return false, err
	}

	//err = l.Conn.Bind(l.Config.BindUserName, l.Config.BindPassword)
	//if err != nil {
	//	return false, nil
	//}

	return true, nil
}
