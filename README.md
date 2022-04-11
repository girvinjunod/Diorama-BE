# Backend diorama.id
## Requirements
- Go 1.17
- Postgresql 14
- `.env` file
## How to Run
### Local
- Install all the requirements
- Create a postgresql database named using `diorama.sql` in the db directory
- Run `go build && ./diorama` in the src directory (Note: This is a Bash based command, change it to fit your shell)
### Docker
- Run `docker-compose up` on the root directory
