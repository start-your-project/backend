#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE USER $APP_DB_USER WITH PASSWORD '$APP_DB_PASS';
  CREATE DATABASE $APP_DB_NAME;
  GRANT ALL PRIVILEGES ON DATABASE $APP_DB_NAME TO $APP_DB_USER;
  ALTER DATABASE $APP_DB_NAME OWNER TO $APP_DB_USER;
  \connect $APP_DB_NAME $APP_DB_USER
  BEGIN;
    CREATE TABLE IF NOT EXISTS position
    (
        id serial constraint position_pk primary key,
        name varchar(250) not null,
        count_technologies int,
        requests_count int,
        fast_prediction bool,
        prediction_dttm timestamp

    );

    create unique index position_uindex
                  on position (name);
  COMMIT;

  BEGIN;
      CREATE TABLE IF NOT EXISTS technology
      (
          id serial constraint technology_pk primary key,
          name varchar(250) not null,
          hard_skill bool
      );

      create unique index technology_uindex
              on technology (name);
  COMMIT;

  BEGIN;
        CREATE TABLE IF NOT EXISTS technology_position
        (
            id_technology int REFERENCES technology,
            name_technology varchar(250) not null,
            id_position int REFERENCES position,
            name_position varchar(250) not null,
            distance float,
            professionalism float,
            PRIMARY KEY (id_technology, id_position)
        );
  COMMIT;

  BEGIN;
        CREATE TABLE IF NOT EXISTS users
        (
            id serial constraint users_pk primary key,
            username varchar(50)  not null,
            email    varchar(50)  not null,
            password varchar(250) not null,
            salt     varchar(50)  not null,
            avatar varchar(100),
            email_confirmed bool
        );

        create unique index users_email_uindex
                on users (email);
  COMMIT;

  BEGIN;
        CREATE TABLE IF NOT EXISTS user_position
        (
            id_user int REFERENCES users,
            id_position int REFERENCES position,
            PRIMARY KEY (id_user, id_position)
        );
  COMMIT;

  BEGIN;
        CREATE TABLE IF NOT EXISTS user_technology
        (
            id_user int REFERENCES users,
            id_technology int REFERENCES technology,
            PRIMARY KEY (id_user, id_technology)
        );
  COMMIT;

  BEGIN;
        create table if not exists emails
        (
            id serial constraint emails_pk primary key,
            domen varchar(100)  not null,
            name varchar(100)  not null,
            link varchar(100)  not null
        );
  COMMIT;

  INSERT INTO emails(domen, name, link) VALUES('mail.ru', 'Почта Mail.Ru', 'https://e.mail.ru/');
  INSERT INTO emails(domen, name, link) VALUES('bk.ru', 'Почта Mail.Ru (bk.ru)', 'https://e.mail.ru/');
  INSERT INTO emails(domen, name, link) VALUES('list.ru', 'Почта Mail.Ru (list.ru)', 'https://e.mail.ru/');
  INSERT INTO emails(domen, name, link) VALUES('inbox.ru', 'Почта Mail.Ru (inbox.ru)', 'https://e.mail.ru/');
  INSERT INTO emails(domen, name, link) VALUES('yandex.ru', 'Яндекс.Почта', 'https://mail.yandex.ru/');
  INSERT INTO emails(domen, name, link) VALUES('ya.ru', 'Яндекс.Почта', 'https://mail.yandex.ru/');
  INSERT INTO emails(domen, name, link) VALUES('yandex.ua', 'Яндекс.Почта', 'https://mail.yandex.ua/');
  INSERT INTO emails(domen, name, link) VALUES('yandex.by', 'Яндекс.Почта', 'https://mail.yandex.by/');
  INSERT INTO emails(domen, name, link) VALUES('yandex.kz', 'Яндекс.Почта', 'https://mail.yandex.kz/');
  INSERT INTO emails(domen, name, link) VALUES('yandex.com', 'Yandex.Mail', 'https://mail.yandex.com/');
  INSERT INTO emails(domen, name, link) VALUES('gmail.com', 'Gmail', 'https://mail.google.com/');
  INSERT INTO emails(domen, name, link) VALUES('googlemail.com', 'Gmail', 'https://mail.google.com/');
  INSERT INTO emails(domen, name, link) VALUES('outlook.com', 'Outlook.com', 'https://mail.live.com/');
  INSERT INTO emails(domen, name, link) VALUES('hotmail.com', 'Outlook.com (Hotmail)', 'https://mail.live.com/');
  INSERT INTO emails(domen, name, link) VALUES('live.ru', 'Outlook.com (live.ru)', 'https://mail.live.com/');
  INSERT INTO emails(domen, name, link) VALUES('live.com', 'Outlook.com (live.com)', 'https://mail.live.com/');
  INSERT INTO emails(domen, name, link) VALUES('me.com', 'iCloud Mail', 'https://www.icloud.com/');
  INSERT INTO emails(domen, name, link) VALUES('icloud.com', 'iCloud Mail', 'https://www.icloud.com/');
  INSERT INTO emails(domen, name, link) VALUES('rambler.ru', 'Рамблер-Почта', 'https://mail.rambler.ru/');
  INSERT INTO emails(domen, name, link) VALUES('yahoo.com', 'Yahoo! Mail', 'https://mail.yahoo.com/');
  INSERT INTO emails(domen, name, link) VALUES('ukr.net', 'Почта ukr.net', 'https://mail.ukr.net/');
  INSERT INTO emails(domen, name, link) VALUES('i.ua', 'Почта I.UA', 'http://mail.i.ua/');
  INSERT INTO emails(domen, name, link) VALUES('bigmir.net', 'Почта Bigmir.net', 'http://mail.bigmir.net/');
  INSERT INTO emails(domen, name, link) VALUES('tut.by', 'Почта tut.by', 'https://mail.tut.by/');
  INSERT INTO emails(domen, name, link) VALUES('inbox.lv', 'Inbox.lv', 'https://www.inbox.lv/');
  INSERT INTO emails(domen, name, link) VALUES('mail.kz', 'Почта mail.kz', 'http://mail.kz/');

EOSQL

