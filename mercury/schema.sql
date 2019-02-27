-- set search_path to <path>;

create sequence mercury_spaces_id_seq;

create table mercury_spaces
(
	space varchar
		not null,
	id integer
	    default nextval('mercury_spaces_id_seq'::regclass)
	    not null
		constraint mercury_namespace_pk
		primary key,
	notes character varying[]
		default '{}'::character varying[]
		not null,
	tags character varying[]
		default '{}'::character varying[]
		not null
);

create unique index mercury_namespace_space_uindex
	on mercury_spaces (space);

create table mercury_values
(
	id integer
		not null,
	seq integer
		not null,
	name varchar
		not null,
	values character varying[]
		default '{}'::character varying[]
		not null,
	tags character varying[]
		default '{}'::character varying[]
		not null,
	notes character varying[]
		default '{}'::character varying[]
		not null,
	
	constraint mercury_values_pk
		primary key (id, seq)
);

create index mercury_values_name_index
	on mercury_values (name);

create or replace view mercury_registry_vw as
select s.id
,      seq
,      space
,      name
,      values
,      v.notes
,      v.tags
from mercury_spaces s
join mercury_values v on (s.id = v.id);

create or replace view mercury_groups_vw as
select distinct unnest(values) user_id
,      name group_id
from mercury_registry_vw
where space = 'config.groups';

create or replace view mercury_notify_vw as
select name
,      split_part(rules,' ', 1) as "match"
,      split_part(rules,' ', 2) as "event"
,      split_part(rules,' ', 3) as "method"
,      split_part(rules,' ', 4) as "url"
from (
	select distinct name
	,      unnest(values) rules
	from mercury_registry_vw
	where space = 'config.notify'
) tt;

create or replace view mercury_rules_vw as
select user_id
,      split_part(rules,' ', 1) as "role"
,      split_part(rules,' ', 2) as "type"
,      split_part(rules,' ', 3) as "match"
from mercury_groups_vw g
join (
	select distinct name group_id
	,      unnest(values) rules
	from mercury_registry_vw
	where space = 'config.policy'
) tt on (g.group_id = tt.group_id);