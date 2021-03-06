module bfimpl

go 1.12

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190611184440-5c40567a22f8

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190616124812-15dcb6c0061f

replace golang.org/x/net => github.com/golang/net v0.0.0-20190613194153-d28f0bde5980

replace golang.org/x/text => github.com/golang/text v0.3.2

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190614205625-5aca471b1d59

replace golang.org/x/image => github.com/golang/image v0.0.0-20200927104501-e162460cd6b5

replace github.com/derekparker/delve => github.com/go-delve/delve v1.2.0

replace golang.org/x/arch => github.com/golang/arch v0.0.0-20190312162104-788fe5ffcd8c

replace github.com/derekparker/delve/terminal => github.com/go-delve/delve/terminal v1.2.0

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.1
	github.com/astaxie/beego v1.11.1
	github.com/go-ini/ini v1.61.0 // indirect
	github.com/go-ldap/ldap v0.0.0-20200627001853-45321a6717b4
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.9
	github.com/kr/pretty v0.2.0 // indirect
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/robfig/cron/v3 v3.0.0
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a
	gopkg.in/ini.v1 v1.61.0 // indirect
)
