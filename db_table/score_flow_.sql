create table score_flow_
(
    id           int unsigned auto_increment
        primary key,
    oid          varchar(128)     default ''                not null comment '订单id',
    score_id     int unsigned     default 0                 not null comment '积分类型id',
    domain       varchar(64)      default ''                not null comment '域',

    o_type       tinyint unsigned default 0                 not null comment '操作类型. 1=增加, 2=扣除, 3=重置',
    old_score    bigint unsigned  default 0                 not null comment '原始积分',
    change_score bigint unsigned  default 0                 not null comment '变更积分',
    o_status     tinyint unsigned default 1                 not null comment '操作状态. 1=成功, 2=余额不足',
    result_score bigint unsigned  default 0                 not null comment '结果积分',

    uid          varchar(128)     default ''                not null comment '用户唯一标识',
    remark       varchar(1024)    default ''                not null comment '备注',

    ctime        datetime         default current_timestamp not null,
    constraint oid_index
        unique (oid)
)
    comment '积分流水';

create index uid_index on score_flow_ (uid);
