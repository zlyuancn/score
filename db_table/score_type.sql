create table score_type
(
    id                      int unsigned auto_increment comment '积分类型id, 用于区分业务'
        primary key,
    score_name              varchar(32)      default ''                                            not null comment '积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么',
    start_time              datetime                                                               null comment '生效时间',
    end_time                datetime                                                               null comment '失效时间',

    order_status_expire_day tinyint unsigned default 30                                            not null comment '订单状态保留多少天, 0表示永久',

    remark                  varchar(1024)    default ''                                            not null comment '备注',
    ctime                   datetime         default current_timestamp                             not null comment '创建时间',
    utime                   datetime         default current_timestamp ON UPDATE CURRENT_TIMESTAMP not null comment '更新时间'
)
    comment '积分类型';
