# Chirpy

A Twitter-like REST API built in Go as part of the Boot.dev backend curriculum.

## What It Does

Chirpy is a fully functional backend API that lets users create accounts, post short messages (chirps), and interact with the platform securely.

## Features

- User registration and login with hashed passwords (argon2id)
- JWT-based authentication with access and refresh tokens
- Full CRUD for chirps with ownership authorization
- Webhook endpoint for "Chirpy Red" membership upgrades via a fictional payment provider
- PostgreSQL database with goose migrations and SQLC for type-safe queries

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Migrations:** goose
- **Query generation:** SQLC
- **Auth:** argon2id, golang-jwt
