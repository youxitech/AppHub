
-- +migrate Up
create table app(
  id text not null primary key,
  icon text not null,
  name text not null,
  type text not null check(type = 'ios' or type = 'android'),
  bundle_id text not null,
  install_password text not null default '',
  download_count int not null default 0,
  created_at text not null, -- RFC3339
  updated_at text not null -- RFC3339
);

create table version(
  version string not null primary key, -- generated version
  app_id text not null references app(id),
  android_version_code int not null default 0,
  android_version_name text not null default '',
  ios_short_version text not null default '',
  ios_bundle_version text not null default '',
  created_at text not null,
  remark text not null default '',
  download_count int not null default 0
);

create table package(
  id string not null primary key,
  version_id string not null references version(version),
  download_count int not null default 0,
  name string not null unique,
  size int not null,
  created_at text not null,
  remark text not null default ''
);

-- +migrate Down
drop table app;
drop table version;
drop table package;
