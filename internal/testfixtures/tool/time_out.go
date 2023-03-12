package tool

var xx = `begin transaction;

drop table if exists naturals;
create table naturals
( n integer unique primary key asc,
  isprime bool,
  factor integer);

with recursive
  nn (n)
as (
  select 2
  union all
  select n+1 as newn from nn
  where newn < 1e14
)
insert into naturals
select n, 1, null from nn;

insert or replace into naturals
  with recursive
    product (prime,composite)
  as (
    select n, n*n as sqr
      from naturals
      where sqr <= (select max(n) from naturals)
    union all
    select prime, composite+prime as prod
    from
      product
    where
      prod <= (select max(n) from naturals)
  )
select n, 0, prime
from product join naturals
  on (product.composite = naturals.n)
;
commit;`

//if err := GetDb(cancelCtx).Exec(xx).Error; err != nil {
//t.Fatal(err)
//}
