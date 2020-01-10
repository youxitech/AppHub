
-- +migrate Up
create view if not exists simple_app as
  select
    id, name, alias, platform, bundle_id, download_count
  from
    app
;

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

-- +migrate Down
drop view simple_app;
drop view detail_version;
