package controllers

import (
	"bfimpl/models"
	"bfimpl/services"
	"bfimpl/services/log"
	"encoding/json"
)

type UserController struct {
	BaseController
}

// @Title 新增用户
// @Description 新增用户
// @Param	name	body	string	true	"姓名"
// @Param	email	body	string	true	"邮箱"
// @Param	wx		body	string	true	"企业微信"
// @Param	phone	body	string	true	"手机"
// @Param	userType	body	int	true	"用户类型"
// @Param	leaderId	body	int	false	"组长id"
// @Success 200 {object} models.User
// @Failure 500 server err
// @router / [post]
func (u *UserController) AddUser() {
	reqUser := new(models.User)

	err := json.Unmarshal(u.Ctx.Input.RequestBody, reqUser)
	if err != nil {
		u.ErrorOK(MsgInvalidParam)
	}
	b, _ := u.valid.Valid(reqUser)
	if !b {
		log.GLogger.Error("%s:%s", u.valid.Errors[0].Field, u.valid.Errors[0].Message)
		u.ErrorOK(MsgInvalidParam)
	}

	err = services.Slave().Create(reqUser).Error
	if err != nil {
		u.ErrorOK("用户邮箱已存在")
	}
	u.Correct(reqUser)
}

// @Title 资源分配人员列表
// @Description 无
// @Success 200 {object} []models.User
// @Failure 500 server err
// @router /leaders [get]
func (u *UserController) GroupLeaders() {
	//userType = 4
	users := make([]*models.User, 0)
	err := services.Slave().Where("user_type = ?", 4).Find(&users).Error
	if err != nil {
		u.ErrorOK(err.Error())
	}
	u.Correct(users)
}

// @Title 实施人员列表
// @Description 实施人员列表, 任务指派时根据leaderId筛选
// @Param  leaderId	query	int		true	"组长id"
// @Success 200 {object} []models.User
// @Failure 500 server err
// @router /impls [get]
func (u *UserController) Implementers() {
	leadId, _ := u.GetInt("leaderId", 0)
	data := make([]*models.Impler, 0)
	err := services.Slave().Raw("SELECT u.id,u.`name`,t.app_name,t.status,s.service_name,t.real_amount,t.exp_deliver_time "+
		"FROM users u left join tasks t on u.id = t.exe_user_id left join services s on t.real_service_id = s.id WHERE "+
		"u.leader_id = ? order by t.exp_deliver_time", leadId).Scan(&data).Error
	if err != nil {
		u.ErrorOK(err.Error())
	}
	m := make(map[int][]*models.Impler)
	sortData := make(map[int]models.SortImpl)
	for i, tmp := range data {
		if d, ok := m[tmp.Id]; ok {
			if tmp.Status == models.TaskExecute || tmp.Status == models.TaskAssign {
				m[tmp.Id] = append(d, data[i])
			}
			tmpSort := sortData[tmp.Id]
			if tmp.Status == models.TaskExecute {
				tmpSort.ExeNum++
			} else if tmp.Status == models.TaskAssign {
				tmpSort.AssignNum++
			}
			sortData[tmp.Id] = tmpSort
		} else {
			m[tmp.Id] = []*models.Impler{}
			sort := models.SortImpl{Id: tmp.Id, Name: tmp.Name}
			if tmp.Status == models.TaskExecute || tmp.Status == models.TaskAssign {
				m[tmp.Id] = append(m[tmp.Id], data[i])
				if tmp.Status == models.TaskExecute {
					sort.ExeNum = 1
				} else if tmp.Status == models.TaskAssign {
					sort.AssignNum = 1
				}
			}
			sortData[tmp.Id] = sort
		}
	}
	result := make([]models.RspImpl, 0)
	sort := sortImpl(sortData)
	for _, id := range sort {
		t := models.RspImpl{
			SortImpl: sortData[id],
			List:     m[id],
		}
		result = append(result, t)
	}
	u.Correct(result)
}

func sortImpl(data map[int]models.SortImpl) []int {
	arr := make([]models.SortImpl, 0)
	result := make([]int, 0)
	for _, v := range data {
		arr = append(arr, v)
	}
	for i := 1; i < len(arr); i++ {
		for j := 0; j < len(arr)-i; j++ {
			if arr[j].ExeNum > arr[j+1].ExeNum {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			} else if arr[j].ExeNum == arr[j+1].ExeNum {
				if arr[j].AssignNum > arr[j+1].AssignNum {
					arr[j], arr[j+1] = arr[j+1], arr[j]
				}
			}
		}
	}
	for i := range arr {
		result = append(result, arr[i].Id)
	}
	return result
}

// @Title 人员列表
// @Description 按类型查询 \n1: "admin",\n2: "sale",	\n3: "manager",	\n4: "tm",	\n5: "implement"
// @Param	type		query	int		true	"类型"
// @Param	pageSize	query	int		true	"每页条数"
// @Param	pageNum		query	int		true	"页数"
// @Success 200 {object} []models.User
// @Failure 500 server err
// @router /list [get]
func (u *UserController) UserList() {
	userType, _ := u.GetInt("type", 0)
	pageSize, _ := u.GetInt("pageSize", 10)
	pageNum, _ := u.GetInt("pageNum", 1)

	users := make([]*models.User, 0)
	query := services.Slave().Model(models.User{})
	if userType != 0 {
		query = query.Where("user_type = ?", userType)
	}
	total := 0
	err := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		u.ErrorOK(err.Error())
	}
	var res = struct {
		Total int            `json:"total"`
		Users []*models.User `json:"users"`
	}{
		Total: total,
		Users: users,
	}

	u.Correct(res)
}
