
-- users.user_type 1-管理员 2-销售 3-经理 4-资源分配 5-实施 6-hr 7-it 8-财务 9-前台 10-leader
-- 入职流程定义
-- type: Approval, Business
-- purpose: EmployeeEntry, EmployeeLeave, Overtime, Leave, Expense, DeviceRequisition
-- entity: Employee, DeviceRequisition, Overtime, Leave, Expense
insert into workflow_definitions(created_at,updated_at,workflow_type,workflow_purpose,workflow_entity)
values (now(),now(),'Business','EmployeeEntry','Employee');
insert into workflow_definitions(created_at,updated_at,workflow_type,workflow_purpose,workflow_entity)
values (now(),now(),'Business','EmployeeLeave','Employee');
insert into workflow_definitions(created_at,updated_at,workflow_type,workflow_purpose,workflow_entity)
values (now(),now(),'Approval','Overtime','Overtime');
insert into workflow_definitions(created_at,updated_at,workflow_type,workflow_purpose,workflow_entity)
values (now(),now(),'Approval','Leave','Leave');
insert into workflow_definitions(created_at,updated_at,workflow_type,workflow_purpose,workflow_entity)
values (now(),now(),'Approval','Expense','Expense');

-- type: TextField, TextArea
insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),1,1,'TextField','plan_date'),(now(),now(),1,2,'TextField','seat_number'),(now(),now(),1,3,'TextArea','device_req');

insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),3,1,'TextField','NULL'),(now(),now(),3,2,'TextArea','leader_comment'),(now(),now(),3,3,'TextArea','hr_comment');

insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),4,1,'TextField','NULL'),(now(),now(),4,2,'TextArea','leader_comment'),(now(),now(),4,3,'TextArea','hr_comment');

insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),5,1,'TextField','NULL'),
(now(),now(),5,2,'TextArea','leader_comment'),
(now(),now(),5,3,'TextArea','finance_comment');
(now(),now(),5,4,'TextField','NULL'),


-- 新增费用科目
insert into expense_account(created_at,updated_at,expense_account_name,expense_account_code) values
(now(),now(),'餐补费','10001'),
(now(),now(),'交通费(市内)','10002'),
(now(),now(),'团队激励','10003'),
(now(),now(),'活动费','10004'),
(now(),now(),'办公费','10005'),
(now(),now(),'招聘费','10006'),
(now(),now(),'通讯费','10007'),
(now(),now(),'销售费用','10008'),
(now(),now(),'充值费用','10009'),
(now(),now(),'交通费(市外)','10010'),
(now(),now(),'住宿费','10011'),
(now(),now(),'出差补贴','10012'),

                                                                                                    )
-- 新增部门
insert into users(id,created_at,updated_at, name, email,wx,phone,user_type) values
(66,now(),now(),'马俊杰','ralph.ma@broadfun.cn', '','',10),
(84,now(),now(),'孟繁秋','fanqiu.meng@broadfun.cn', '','',10),
(85,now(),now(),'李欣','lane.li@broadfun.cn', '','',10),
(86,now(),now(),'王立卿','Stanley.wang@broadfun.cn', '','',10),
(87,now(),now(),'罗超','chao.luo@broadfun.cn', '','',10),
(88,now(),now(),'姚诚诚','chengcheng.yao@broadfun.cn', '','',10),
(89,now(),now(),'何勤勤','Theresa.he@broadfun.cn', '','',10),
(90,now(),now(),'孙丹峰','danfeng.sun@broadfun.cn', '','',10),
(91,now(),now(),'陈一菲','yifei.chen@broadfun.cn', '','',10);

insert into departments(created_at,updated_at,department_name,department_leader_id) values
(now(),now(),'游戏测试一组',84),
(now(),now(),'游戏测试二组',84),
(now(),now(),'专家实施',85),
(now(),now(),'WeTest外部服务',86),
(now(),now(),'WeTest客户成功',89),
(now(),now(),'WeTest商务1',87),
(now(),now(),'WeTest商务2',88),
(now(),now(),'WeTest私有化',2),
(now(),now(),'WeTest海外先游',2),
(now(),now(),'WeTest企鹅风讯',2),
(now(),now(),'人力资源部',91),
(now(),now(),'财务部',89),
(now(),now(),'IT部',66),
(now(),now(),'QA部',66);
-- 新增财务统计engagement_code
insert into engagement_codes(created_at,updated_at,engagement_code,engagement_code_desc,department_id,code_owner_id) values
(now(),now(),'10001', 'Wetest私有化',5, 2),
(now(),now(),'10002', 'Wetest风讯',5, 3),
(now(),now(),'10003', 'Wetest先游',5, 4);

-- 新增级别
insert into levels(created_at,updated_at,department_id,level_name,cc_rate,oc_rate) values
(now(),now(),1,'T0.0',0.0,0.0),
(now(),now(),1,'T1.1',0.0,0.0),
(now(),now(),1,'T1.2',0.0,0.0),
(now(),now(),1,'T1.3',0.0,0.0),
(now(),now(),1,'T2.1',0.0,0.0),
(now(),now(),1,'T2.2',0.0,0.0),
(now(),now(),1,'T2.3',0.0,0.0),
(now(),now(),2,'T0.0',0.0,0.0),
(now(),now(),2,'T1.1',0.0,0.0),
(now(),now(),2,'T1.2',0.0,0.0),
(now(),now(),2,'T1.3',0.0,0.0),
(now(),now(),2,'T2.1',0.0,0.0),
(now(),now(),2,'T2.2',0.0,0.0),
(now(),now(),2,'T2.3',0.0,0.0),
(now(),now(),3,'T0.0',0.0,0.0),
(now(),now(),3,'T1.1',0.0,0.0),
(now(),now(),3,'T1.2',0.0,0.0),
(now(),now(),3,'T1.3',0.0,0.0),
(now(),now(),3,'T2.1',0.0,0.0),
(now(),now(),3,'T2.2',0.0,0.0),
(now(),now(),3,'T2.3',0.0,0.0),
(now(),now(),4,'T0.0',0.0,0.0),
(now(),now(),4,'T1.1',0.0,0.0),
(now(),now(),4,'T1.2',0.0,0.0),
(now(),now(),4,'T1.3',0.0,0.0),
(now(),now(),4,'T2.1',0.0,0.0),
(now(),now(),4,'T2.2',0.0,0.0),
(now(),now(),4,'T2.3',0.0,0.0),
(now(),now(),5,'T0.0',0.0,0.0),
(now(),now(),5,'T1.1',0.0,0.0),
(now(),now(),5,'T1.2',0.0,0.0),
(now(),now(),5,'T1.3',0.0,0.0),
(now(),now(),5,'T2.1',0.0,0.0),
(now(),now(),5,'T2.2',0.0,0.0),
(now(),now(),5,'T2.3',0.0,0.0),
(now(),now(),6,'T0.0',0.0,0.0),
(now(),now(),6,'T1.1',0.0,0.0),
(now(),now(),6,'T1.2',0.0,0.0),
(now(),now(),6,'T1.3',0.0,0.0),
(now(),now(),6,'T2.1',0.0,0.0),
(now(),now(),6,'T2.2',0.0,0.0),
(now(),now(),6,'T2.3',0.0,0.0),
(now(),now(),7,'T0.0',0.0,0.0),
(now(),now(),7,'T1.1',0.0,0.0),
(now(),now(),7,'T1.2',0.0,0.0),
(now(),now(),7,'T1.3',0.0,0.0),
(now(),now(),7,'T2.1',0.0,0.0),
(now(),now(),7,'T2.2',0.0,0.0),
(now(),now(),7,'T2.3',0.0,0.0),
(now(),now(),8,'T0.0',0.0,0.0),
(now(),now(),8,'T1.1',0.0,0.0),
(now(),now(),8,'T1.2',0.0,0.0),
(now(),now(),8,'T1.3',0.0,0.0),
(now(),now(),8,'T2.1',0.0,0.0),
(now(),now(),8,'T2.2',0.0,0.0),
(now(),now(),8,'T2.3',0.0,0.0),
(now(),now(),9,'T0.0',0.0,0.0),
(now(),now(),9,'T1.1',0.0,0.0),
(now(),now(),9,'T1.2',0.0,0.0),
(now(),now(),9,'T1.3',0.0,0.0),
(now(),now(),9,'T2.1',0.0,0.0),
(now(),now(),9,'T2.2',0.0,0.0),
(now(),now(),9,'T2.3',0.0,0.0);

-- 添加HR , it
insert into users(created_at,updated_at, name, email,wx,phone,user_type) values
(now(),now(),'HR','test01@broadfun.cn', '','',6);

insert into users(created_at,updated_at, name, email,wx,phone,user_type) values
(now(),now(),'楼易凯','test02@broadfun.cn', '','',7);

-- workflow  流程定义id和entityID唯一索引
alter table workflows add unique `uid_workflows_wfd_e`(`workflow_definition_id`,`entity_id`);

alter table attendances add UNIQUE uix_attendance_name_date(`name`,`attendance_date`);

alter table attendance_tmp add UNIQUE uix_attendance_name_date(`name`,`check_time`);