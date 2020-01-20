
-- +migrate Up
alter table package
  add column channel text not null default 'android'
;

-- +migrate Down
drop view detail_version;

create table package2(
  id string not null primary key,
  version_id int not null references version(id) on delete cascade,
  download_count int not null default 0,
  name string not null unique,
  size int not null,
  created_at datetime not null,
  remark text not null default '',
  ios_package_type text not null default 'ad-hoc'
    check(ios_package_type = '' or ios_package_type = 'ad-hoc' or ios_package_type = 'in-house' or ios_package_type = 'app-store'),
  ios_device_list text not null default ''
);

insert into
  package2(id, version_id, download_count , name, size, created_at, remark, ios_package_type, ios_device_list)
select
  id, version_id, download_count , name, size, created_at, remark, ios_package_type, ios_device_list
from
  package;

drop table package;

alter table package2 rename to package;

create view if not exists detail_version as
  select
    v.*,
    count(*) as package_count,
    datetime(max(p.created_at)) as updated_at
  from
    version v
  left join package p on v.id = p.version_id
  group by v.id
  order by sort_key desc
;
