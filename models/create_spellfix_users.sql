.load ./spellfix

drop table if exists spellfix_users;
create virtual table spellfix_users using spellfix1;

insert into spellfix_users(word) select name from users;

-- drop table if exists spellfix_gyms;
-- create virtual table spellfix_gyms using spellfix1;
--
-- insert into spellfix_gyms(word) select name from gym;
