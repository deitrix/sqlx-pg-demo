create database leasing;

create table cars
(
    id      serial primary key,
    details jsonb not null
);
