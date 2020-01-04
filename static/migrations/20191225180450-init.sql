
-- +migrate Up
create table app(
  id text not null primary key,
  name text not null,
  platform text not null check(platform = 'ios' or platform = 'android'),
  bundle_id text not null unique,
  install_password text not null default '',
  download_count int not null default 0
);

create table version(
  id string not null primary key, -- generated full version string
  app_id text not null references app(id),
  android_version_code text not null default '',
  android_version_name text not null default '',
  ios_short_version text not null default '',
  ios_bundle_version text not null default '',
  sort_key int not null,
  remark text not null default '',
  download_count int not null default 0
);

create table package(
  id string not null primary key,
  version_id string not null references version(id),
  download_count int not null default 0,
  name string not null unique,
  size int not null,
  created_at datetime not null,
  remark text not null default ''
);

-- +migrate Down
drop table app;
drop table version;
drop table package;
