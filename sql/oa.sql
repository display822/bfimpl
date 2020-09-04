
-- 入职流程定义
-- type: Approval, Business
-- purpose: EmployeeEntry, Overtime, Leave, Expense, DeviceRequisition
-- entity: Employee, DeviceRequisition, Overtime, Leave, Expense
insert into workflow_definitions(created_at,updated_at,workflow_type,workflow_purpose,workflow_entity)
values (now(),now(),'Business','EmployeeEntry','Employee');

-- type: TextField, TextArea
insert into workflow_form_element_defs(created_at,updated_at,workflow_definition_id,element_seq,element_type, element_name) values
(now(),now(),1,1,'TextField','plan_time'),(now(),now(),1,2,'TextField','seat_number'),(now(),now(),1,3,'TextArea','device_req');

-- 新增部门
insert into employees(created_at,updated_at,name) values
(now(),now(),'孟繁秋'),
(now(),now(),'李欣'),
(now(),now(),'王立卿'),
(now(),now(),'牛茜茜');
insert into departments(created_at,updated_at,department_name,department_leader_id) values
(now(),now(),'游戏测试部',1),
(now(),now(),'专家实施',2),
(now(),now(),'WeTest深度',3),
(now(),now(),'WeTest商务',0),
(now(),now(),'合研部门',4),
(now(),now(),'先游Gamer',0),
(now(),now(),'APM',0),
(now(),now(),'B站外派',0),
(now(),now(),'职能部门',0);

-- 新增级别
insert into levels(created_at,updated_at,department_id,level_name,cc_rate,oc_rate) values
(now(),now(),1,'高级xxx',0.0,0.0),
(now(),now(),2,'高级xxx',0.0,0.0),
(now(),now(),3,'高级xxx',0.0,0.0),
(now(),now(),4,'高级xxx',0.0,0.0),
(now(),now(),5,'高级xxx',0.0,0.0),
(now(),now(),6,'高级xxx',0.0,0.0),
(now(),now(),7,'高级xxx',0.0,0.0),
(now(),now(),8,'高级xxx',0.0,0.0),
(now(),now(),9,'高级xxx',0.0,0.0);

