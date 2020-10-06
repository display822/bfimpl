
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

-- type: TextField, TextArea
insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),1,1,'TextField','plan_date'),(now(),now(),1,2,'TextField','seat_number'),(now(),now(),1,3,'TextArea','device_req');

insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),3,1,'TextField','NULL'),(now(),now(),3,2,'TextArea','leader_comment'),(now(),now(),3,3,'TextArea','hr_comment');

insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),4,1,'TextField','NULL'),(now(),now(),4,2,'TextArea','leader_comment'),(now(),now(),4,3,'TextArea','hr_comment');


-- 新增部门
insert into users(created_at,updated_at, name, email,wx,phone,user_type) values
(now(),now(),'孟繁秋','fanqiu.meng@broadfun.cn', '','',1),
(now(),now(),'李欣','lane.li@broadfun.cn', '','',1),
(now(),now(),'王立卿','Stanley.wang@broadfun.cn', '','',1),
(now(),now(),'罗超','chao.luo@broadfun.cn', '','',10),
(now(),now(),'范硕瑀','barry.fan@broadfun.cn', '','',10);
insert into departments(created_at,updated_at,department_name,department_leader_id) values
(now(),now(),'游戏测试部',84),
(now(),now(),'专家实施',85),
(now(),now(),'WeTest深度',86),
(now(),now(),'WeTest商务',0),
(now(),now(),'合研部门',2),
(now(),now(),'先游Gamer',0),
(now(),now(),'APM',0),
(now(),now(),'B站外派',0),
(now(),now(),'职能部门',0);
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