-- +goose Up
ALTER TABLE Users
ADD is_chirpy_red BOOLEAN NOT NULL DEFAULT 'false';
