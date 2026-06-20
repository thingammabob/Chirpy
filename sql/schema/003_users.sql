-- +goose Up
ALTER TABLE Users
ADD hashed_password TEXT NOT NULL DEFAULT 'unset';
